package tsgo

type AdaptationExtension struct {
	AdaptationExtensionLength uint8
	LegalTimeWindowFlag       bool
	PiecewiseRateFlag         bool
	SeamlessSpliceFlag        bool
	LegalTimeWindowValidFlag  bool   `json:",omitempty"`
	LegalTimeWindowOffset     uint16 `json:",omitempty"`
	PiecewiseRate             uint32 `json:",omitempty"`
	SpliceType                uint8  `json:",omitempty"`
	DTSNextAccessUnit         uint64 `json:",omitempty"`
}

func ParseAdaptationExtension(b []byte) (*AdaptationExtension, error) {
	i := 2 // read head
	e := &AdaptationExtension{
		AdaptationExtensionLength: uint8(b[0]),
		LegalTimeWindowFlag:       b[1]>>7&0x1 > 0,
		PiecewiseRateFlag:         b[1]>>6&0x1 > 0,
		SeamlessSpliceFlag:        b[1]>>7&0x1 > 0,
	}

	if e.LegalTimeWindowFlag {
		e.LegalTimeWindowValidFlag = b[i]>>7*0x1 > 0
		e.LegalTimeWindowOffset = uint16(b[i]&0x7f)<<8 | uint16(b[i+1])
		i += 2
	}

	if e.PiecewiseRateFlag {
		e.PiecewiseRate = uint32(b[i]&0x3f)<<16 | uint32(b[i+1])<<8 | uint32(b[i+2])
		i += 3
	}

	if e.SeamlessSpliceFlag {
		e.SpliceType = uint8(b[i]&0xf0) >> 4
		// TODO(walker) DTS Next Access Unit
	}

	return e, nil
}
