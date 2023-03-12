package tsgo

type TableID uint16

const (
	PAT TableID = 0x00
	PMT         = 0x02
	NIL         = 0xff
)

func isKnownTable(id TableID) bool {
	switch id {
	case PAT, PMT:
		return true

	default:
		return false
	}
}

type ProgramSpecificInformation struct {
	// Present at the start of the TS packet payload signaled by
	// the payload_unit_start_indicator bit in the TS header. Used to set packet
	// alignment bytes or content before the start of tabled payload data.
	PointerField uint8

	// Program specific information tables repeated until the end of the payload.
	Sections []*PSISection
}

type PSISection struct {
	Header *PSISectionHeader
	Syntax *PSISectionSyntax
}

type PSISectionHeader struct {
	// Table ID, that defines the structure of the syntax section and other
	// contained data. As an exception, if this is the byte that immediately
	// follow previous table section and is set to 0xFF, then it indicates that
	// the repeat of table section end here and the rest of TS packet payload
	// shall be stuffed with 0xFF. Consequently, the value 0xFF shall not be used
	// for the Table Identifier.
	TableID uint8

	// A flag that indicates if the syntax section follows the section length. The
	// PAT, PMT, and CAT all set this to 1.
	SectionSyntaxIndicator bool

	// The PAT, PMT, and CAT all set this to 0. Other tables set this to 1. The
	// PAT, PMT, and CAT all set this to 0. Other tables set this to 1.
	PrivateBit bool

	// The number of bytes that follow for the syntax section (with CRC value)
	// and/or table data. These bytes must not exceed a value of 1021.
	SectionLength uint16
}

type PSISectionSyntax struct {
	// Informational only identifier. The PAT uses this for the transport stream
	// identifier and the PMT uses this for the Program number.
	TableIDExtension uint16

	// Syntax version number. Incremented when data is changed and wrapped around
	// on overflow for values greater than 32.
	VersionNumber uint8

	// Indicates if data is current in effect or is for future use. If the bit is
	// flagged on, then the data is to be used at the present moment.
	CurrentNextIndicator bool

	// This is an index indicating which table this is in a related sequence of
	// tables. The first table starts from 0.
	SectionNumber uint8

	// This indicates which table is the last table in the sequence of tables.
	LastSectionNumber uint8

	// A checksum of the entire table excluding the pointer field, pointer filler
	// bytes and the trailing CRC32.
	CRC32 uint32

	PAT *ProgramAssociationTable `json:",omitempty"`
	PMT *ProgramMapTable         `json:",omitempty"`
}

func (p *Parser) ParseProgramSpecificInformation() *ProgramSpecificInformation {
	return &ProgramSpecificInformation{
		Sections: p.ParseProgramSpecificInformationSections(),
	}
}

func (p *Parser) ParseProgramSpecificInformationSections() []*PSISection {
	sections := make([]*PSISection, 0)

	for i := 0; i < len(p.b); {
		header := p.ParsePSISectionHeader()

		if !isKnownTable(TableID(header.TableID)) {
			return sections
		}

		syntax := p.ParsePSISectionSyntax()

		l := int(header.SectionLength) - 4
		switch TableID(header.TableID) {
		case PAT:
			syntax.PAT = p.ParseProgramAssociationTable(l)

		case PMT:
			syntax.PMT = p.ParseProgramMapTable(l)
		}

		bs := p.ReadBytes(4)
		syntax.CRC32 = uint32(bs[0])<<24 | uint32(bs[1])<<16 |
			uint32(bs[2])<<8 | uint32(bs[3])

		sections = append(sections, &PSISection{
			Header: header,
			Syntax: syntax,
		})
	}

	return sections
}

func (p *Parser) ParsePSISectionHeader() *PSISectionHeader {
	bs := p.ReadBytes(3)

	return &PSISectionHeader{
		TableID:                uint8(bs[0]),
		SectionSyntaxIndicator: bs[1]>>7&0x1 > 0,
		PrivateBit:             bs[1]>>6&0x1 > 0,
		SectionLength:          uint16(bs[1]&0xf)<<8 | uint16(bs[2]),
	}
}

func (p *Parser) ParsePSISectionSyntax() *PSISectionSyntax {
	bs := p.ReadBytes(5)

	return &PSISectionSyntax{
		TableIDExtension:     uint16(bs[0])<<8 | uint16(bs[1]),
		VersionNumber:        uint8(bs[2]&0x3f) >> 1,
		CurrentNextIndicator: bs[2]&0x1 > 0,
		SectionNumber:        uint8(bs[3]),
		LastSectionNumber:    uint8(bs[4]),
	}
}
