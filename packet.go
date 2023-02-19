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

	AdaptationField *AdaptationField `json:",omitempty"`

	ProgramAssociationTable *ProgramAssociationTable `json:",omitempty"`
}

func ParsePacket(b []byte) (*Packet, error) {
	if b[0]&0xFF != 0x47 {
		return nil, fmt.Errorf("Not a sync byte")
	}

	i := 4 // read head
	h := &Packet{
		PacketID:                   uint16(b[1]&0x1F)<<8 | uint16(b[2]),
		TransportErrorIndicator:    b[1]>>7&0x1 > 0,
		PayloadUnitStartIndicator:  b[1]>>6&0x1 > 0,
		TransportPriority:          b[1]>>5&0x1 > 0,
		TransportScramblingControl: uint8(b[3] & 0xC0),
		AdaptationFieldControl:     uint8(b[3] >> 4 & 0x3),
		ContinuityCounter:          uint8(b[3] >> 6 & 0x3),
	}

	if h.AdaptationFieldControl >= 3 {
		af, err := ParseAdaptationField(b[i:])
		if err != nil {
			return nil, err
		}

		i += int(af.AdaptationFieldLength)
		if ext := af.AdaptationExtension; ext != nil {
			i += int(ext.AdaptationExtensionLength)
		}

		h.AdaptationField = af
	}

	if h.AdaptationFieldControl % 2 == 1 {
		i += 1 // TODO (parse Payload pointer)
		switch h.PacketID {
		case 0:
			h.ProgramAssociationTable = ParseProgramAssociationTable(b[i:])
		}
	}

	return h, nil
}
