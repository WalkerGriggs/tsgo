package tsgo

type AdaptationField struct {
	AdaptationFieldLength             uint8
	DiscontinuityIndicator            bool
	RandomAccessIndicator             bool
	ElementaryStreamPriorityIndicator bool
	ProgramClockReferenceFlag         bool
	OriginalProgramClockReferenceFlag bool
	SplicingPointFlag                 bool
	TransportPrivateDataFlag          bool
	AdaptationFieldExtensionFlag      bool
	ProgramClockReference             *ClockReference      `json:",omitempty"`
	OriginalProgramClockReference     *ClockReference      `json:",omitempty"`
	SpliceCountdown                   int8                 `json:",omitempty"`
	TransportPrivateDataLength        uint8                `json:",omitempty"`
	TransportPrivateData              []byte               `json:",omitempty"`
	AdaptationExtension               *AdaptationExtension `json:",omitempty"`
}

func (p *Parser) ParseAdaptationField() (*AdaptationField, error) {
	bs := p.ReadBytes(2)

	a := &AdaptationField{
		AdaptationFieldLength:             uint8(bs[0]),
		DiscontinuityIndicator:            bs[1]>>7&0x1 > 0,
		RandomAccessIndicator:             bs[1]>>6&0x1 > 0,
		ElementaryStreamPriorityIndicator: bs[1]>>5&0x1 > 0,
		ProgramClockReferenceFlag:         bs[1]>>4&0x1 > 0,
		OriginalProgramClockReferenceFlag: bs[1]>>3&0x1 > 0,
		SplicingPointFlag:                 bs[1]>>2&0x1 > 0,
		TransportPrivateDataFlag:          bs[1]>>1&0x1 > 0,
		AdaptationFieldExtensionFlag:      bs[1]&0x1 > 0,
	}

	if a.ProgramClockReferenceFlag {
		a.ProgramClockReference = ParseClockReference(p.ReadBytes(6))
	}

	if a.OriginalProgramClockReferenceFlag {
		a.OriginalProgramClockReference = ParseClockReference(p.ReadBytes(6))
	}

	if a.SplicingPointFlag {
		a.SpliceCountdown = int8(p.ReadByte())
	}

	if a.TransportPrivateDataFlag {
		a.TransportPrivateDataLength = uint8(p.ReadByte())
		a.TransportPrivateData = p.ReadBytes(uint(a.TransportPrivateDataLength))
	}

	if a.AdaptationFieldExtensionFlag {
		l := uint(p.ReadByte())
		p.Dec(1)

		var err error
		if a.AdaptationExtension, err = ParseAdaptationExtension(p.ReadBytes(l)); err != nil {
			return nil, err
		}
	}

	return a, nil
}
