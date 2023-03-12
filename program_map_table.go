package tsgo

type ProgramMapTable struct {
	// The packet identifier that contains the program clock reference used to
	// improve the random access accuracy of the stream's timing that is derived
	// from the program timestamp. If this is unused. then it is set to 0x1FFF
	// (all bits on).
	ProgramClockReferencePID uint16

	// The number of bytes that follow for the program descriptors.
	ProgramInfoLength uint16

	// When the program info length is non-zero, this is the program info length
	// number of program descriptor bytes.
	ProgramDescriptors []byte

	// The streams used in this program map.
	ElementaryStreamInfoData []*ElementaryStreamInfo
}

type ElementaryStreamInfo struct {
	// This defines the structure of the data contained within the elementary
	// packet identifier.
	StreamType uint8

	// The packet identifier that contains the stream type data.
	ElementaryPID uint16

	// The number of bytes that follow for the elementary stream descriptors.
	ElementaryStreamLength uint16

	// TODO ElementaryStreamDescriptors
}

func (p *Parser) ParseProgramMapTable(l int) *ProgramMapTable {
	bs := p.ReadBytes(4)

	pmt := &ProgramMapTable{
		ProgramClockReferencePID: uint16(bs[0]&0x1f)<<8 | uint16(bs[1]),
		ProgramInfoLength:        uint16(bs[2]&0x3)<<8 | uint16(bs[3]),
		ElementaryStreamInfoData: make([]*ElementaryStreamInfo, 0),
	}

	// Iterate over N 8-byte chunks, where N is the program info length, starting
	// at the 4th byte.
	for i := 0; i < int(pmt.ProgramInfoLength); i++ {
		_ = p.ReadBytes(8)
		l -= 8
		// todo
	}

	// Iterate over 8 byte chunks until the end of the section, accounting for the
	// 4 byte CRC32.
	for i := 0; i < l/8; i += 8 {
		bs = p.ReadBytes(8)

		info := &ElementaryStreamInfo{
			StreamType:             uint8(bs[0]),
			ElementaryPID:          uint16(bs[1]&0x1f)<<8 | uint16(bs[2]),
			ElementaryStreamLength: uint16(bs[3]&0x3)<<8 | uint16(bs[4]),
		}

		pmt.ElementaryStreamInfoData = append(pmt.ElementaryStreamInfoData, info)
	}

	return pmt
}
