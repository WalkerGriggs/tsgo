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

func (p *Parser) ParseAdaptationField(b []byte) (*AdaptationField, error) {
	i := 2 // read head
	a := &AdaptationField{
		AdaptationFieldLength:             uint8(b[0]),
		DiscontinuityIndicator:            b[1]>>7&0x1 > 0,
		RandomAccessIndicator:             b[1]>>6&0x1 > 0,
		ElementaryStreamPriorityIndicator: b[1]>>5&0x1 > 0,
		ProgramClockReferenceFlag:         b[1]>>4&0x1 > 0,
		OriginalProgramClockReferenceFlag: b[1]>>3&0x1 > 0,
		SplicingPointFlag:                 b[1]>>2&0x1 > 0,
		TransportPrivateDataFlag:          b[1]>>1&0x1 > 0,
		AdaptationFieldExtensionFlag:      b[1]&0x1 > 0,
	}

	if a.ProgramClockReferenceFlag {
		a.ProgramClockReference = ParseClockReference(b[i : i+7])
		i += 6
	}

	if a.OriginalProgramClockReferenceFlag {
		a.OriginalProgramClockReference = ParseClockReference(b[i : i+7])
		i += 6
	}

	if a.SplicingPointFlag {
		a.SpliceCountdown = int8(b[i])
		i += 1
	}

	if a.TransportPrivateDataFlag {
		a.TransportPrivateDataLength = uint8(b[i])
		i += 1

		a.TransportPrivateData = b[i : i+int(a.TransportPrivateDataLength+1)]
		i += int(a.TransportPrivateDataLength)
	}

	if a.AdaptationFieldExtensionFlag {
		l := int(b[i])

		var err error
		if a.AdaptationExtension, err = ParseAdaptationExtension(b[i : i+l+1]); err != nil {
			return nil, err
		}
	}

	return a, nil
}
