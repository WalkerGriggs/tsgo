package tsgo

type PacketizedElementaryStream struct {
	// Audio streams (0xC0-0xDF), Video streams (0xE0-0xEF)
	StreamID uint8

	// Specifies the number of bytes remaining in the packet after this
	// field. Can be zero. If the PES packet length is set to zero, the PES
	// packet can be of any length. A value of zero for the PES packet
	// length can be used only when the PES packet payload is a video
	// elementary stream.
	PacketLength uint16

	// not present in case of Padding stream & Private stream 2 (navigation
	// data)
	OptionalPESHeader *OptionalPacketizedElementaryHeader

	// See elementary stream. In the case of private streams the first byte
	// of the payload is the sub-stream number.
	Data []byte
}

type OptionalPacketizedElementaryHeader struct {
	// 00 implies not scrambled
	ScramblingControl uint8

	// 1 indicates that the PES packet header is immediately followed by the
	// video start code or audio syncword
	DataAlignmentIndicator bool

	// 1 implied copyrighted
	Copyright bool

	// 1 implied original
	Original bool

	// 11 = both present, 01 is forbidden, 10 = only PTS, 00 = no PTS or DTS
	PTSDTSIndicator uint8

	Priority               bool
	ESCRFlag               bool
	ESRateFlag             bool
	DSMTrickModeFlag       bool
	AdditionalCopyInfoFlag bool
	CRCFlag                bool
	ExtensionFlag          bool

	// Gives the length of the remainder of the PES header in bytes
	HeaderLength uint8
}

func (s *PacketizedElementaryStream) IsVideoStream() bool {
	return s.StreamID == 0xe0 || s.StreamID == 0xfd
}

func (p *Parser) ParsePacketizedElementaryStream() *PacketizedElementaryStream {
	bs := p.ReadBytes(3)
	pes := &PacketizedElementaryStream{
		StreamID:     uint8(bs[0]),
		PacketLength: uint16(bs[1])<<8 | uint16(bs[2]),
	}

	bs = p.ReadBytes(3)
	pes.OptionalPESHeader = &OptionalPacketizedElementaryHeader{
		ScramblingControl:      uint8(bs[0]) >> 4 & 0x3,
		Priority:               uint8(bs[0])&0x8 > 0,
		DataAlignmentIndicator: uint8(bs[0])&0x4 > 0,
		Copyright:              uint8(bs[0])&0x2 > 0,
		Original:               uint8(bs[0])&0x1 > 0,
		PTSDTSIndicator:        uint8(bs[1]) >> 6 & 0x3,
		ESCRFlag:               uint8(bs[1])&0x20 > 0,
		ESRateFlag:             uint8(bs[1])&0x10 > 0,
		DSMTrickModeFlag:       uint8(bs[1])&0x8 > 0,
		AdditionalCopyInfoFlag: uint8(bs[1])&0x4 > 0,
		CRCFlag:                uint8(bs[1])&0x2 > 0,
		ExtensionFlag:          uint8(bs[1])&0x1 > 0,
		HeaderLength:           uint8(bs[2]),
	}

	return pes
}
