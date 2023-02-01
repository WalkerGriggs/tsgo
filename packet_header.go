package tsgo

import (
	"fmt"
)

type PacketHeader struct {
	PacketID                   uint16
	TransportErrorIndicator    bool
	PayloadUnitStartIndicator  bool
	TransportPriority          bool
	TransportScramblingControl uint8
	AdaptationFieldControl     uint8
	ContinuityCounter          uint8
	AdaptationField            *AdaptationField `json:",omitempty"`
}

func ParseHeader(b []byte) (*PacketHeader, error) {
	if b[0]&0xFF != 0x47 {
		return nil, fmt.Errorf("Not a sync byte")
	}

	h := &PacketHeader{
		PacketID:                   uint16(b[1]&0x1F)<<8 | uint16(b[2]),
		TransportErrorIndicator:    b[1]>>7&0x1 > 0,
		PayloadUnitStartIndicator:  b[1]>>6&0x1 > 0,
		TransportPriority:          b[1]>>5&0x1 > 0,
		TransportScramblingControl: uint8(b[3] & 0xC0),
		AdaptationFieldControl:     uint8(b[3] >> 4 & 0x3),
		ContinuityCounter:          uint8(b[3] >> 6 & 0x3),
	}

	if h.AdaptationFieldControl >= 3 {
		var err error
		if h.AdaptationField, err = ParseAdaptationField(b[4:]); err != nil {
			return nil, err
		}
	}

	return h, nil
}
