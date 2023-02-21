package tsgo

import (
	"fmt"
)

type Packet struct {
	// PacketID is the 13-bit field, indicating the type of the data stored in the
	// packet payload. PacketIDs are not unique and many are reserved.
	PacketID uint16

	// TransportErrorIndicator is used to indicate when at least 1 uncorrected
	// error bit exists in the packaged. When set to 1, this bit should never be
	// reset to '0' unless the error has been corrected.
	TransportErrorIndicator bool

	// PayloadUnitStartIndicator is used to indicate when a packet contains the
	// first byte of a new payload unit. This field is useful for a receiver who
	// started reader mid-transmission, and needs to know when to start extracting
	// payload data.
	//
	// When the payload of the  packet contains PES packet data and the flag is
	// true, the payload of the packet will start with the first  byte of a PES
	// packet.
	//
	// When the payload of the packet contains PSI data and the flag is set to
	// true, the first byte of the payload carried the pointer field.
	PayloadUnitStartIndicator bool

	// TransportPriority is used to indicate that the current packet has a higher
	// priority than other packets with the same pid.
	TransportPriority bool

	// TransportScramblingControl is used to select the DVB-CSA or ATSC DES packet
	// encryption. Any value greater than 1 indicates that the packet payload has
	// been scrambled.
	TransportScramblingControl uint8

	// AdaptationFieldControl is used to specify both if and where the adaptation
	// field can be found relative to the payload.
	AdaptationFieldControl uint8

	// ContinuityCounter is used to enumerate packets with the same PID. Each
	// packet that contains a payload (as indicated by the adaptation field
	// control) will increment the continuity counter.
	//
	// A packet is continuous if counter is one greater than the counter of the
	// previous packet with the same PID (ignoring packets without a payload).
	//
	// A packet is discontinuous if the above conditions are not met, or the
	// adaptation field's discontinuity indicator is true.
	ContinuityCounter uint8

	// An optional variable-length extension field of the fixed-length TS
	// Packet header, intended to convey clock references and timing and
	// synchronization information as well as stuffing over an MPEG-2 Multiplex
	AdaptationField *AdaptationField `json:",omitempty"`

	// Program specific metadata like program association, conditional access,
	// program mapping, and network information. This information will never be
	// scrambled.
	ProgramSpecificInformation *ProgramSpecificInformation `json:",omitempty"`
}

type Parser struct {
	ByteReader
	ProgramMap map[uint16]uint16
}

func (p *Parser) ParsePacket(in []byte) (*Packet, error) {
	p.ByteReader = ByteReader{b: in}

	if p.ReadByte()&0xFF != 0x47 {
		return nil, fmt.Errorf("Not a sync byte")
	}

	bs := p.ReadBytes(3)
	h := &Packet{
		TransportErrorIndicator:    bs[0]>>7&0x1 > 0,
		PayloadUnitStartIndicator:  bs[0]>>6&0x1 > 0,
		TransportPriority:          bs[0]>>5&0x1 > 0,
		PacketID:                   uint16(bs[0])&0x1F<<8 | uint16(bs[1]),
		TransportScramblingControl: uint8(bs[2] & 0xC0),
		AdaptationFieldControl:     uint8(bs[2] >> 4 & 0x3),
		ContinuityCounter:          uint8(bs[2] >> 6 & 0x3),
	}

	if h.AdaptationFieldControl >= 3 {
		var err error
		h.AdaptationField, err = p.ParseAdaptationField()
		if err != nil {
			return nil, err
		}
	}

	if h.AdaptationFieldControl%2 == 1 && isPSI(h.PacketID, p.ProgramMap) {
		h.ProgramSpecificInformation = p.ParseProgramSpecificInformation()
	}

	return h, nil
}

func isPSI(pid uint16, pm map[uint16]uint16) bool {
	if _, ok := pm[pid]; ok {
		return true
	}
	return pid == 0
}

func isPES(b []byte) bool {
	return uint32(b[0])<<16|uint32(b[1])<<8|uint32(b[2]) == 1
}
