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
	Data   *PSISectionData
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
}

type PSISectionData struct {
	TableID TableID                  `json:"-"`
	PAT     *ProgramAssociationTable `json:",omitempty"`
	PMT     *ProgramMapTable         `json:",omitempty"`
}

func (p *Parser) ParseProgramSpecificInformation(b []byte) *ProgramSpecificInformation {
	pointerField := uint8(b[0])
	i := int(pointerField) + 1

	sections := p.ParseProgramSpecificInformationSections(b[i:])

	return &ProgramSpecificInformation{
		PointerField: pointerField,
		Sections:     sections,
	}
}

func (p *Parser) ParseProgramSpecificInformationSections(b []byte) []*PSISection {
	sections := make([]*PSISection, 0)

	for i := 0; i < len(b); {
		header := p.ParsePSISectionHeader(b[i : i+3])
		i += 3

		tableID := TableID(header.TableID)
		if tableID == NIL || !isKnownTable(tableID) {
			return sections
		}

		syntax := p.ParsePSISectionSyntax(b[i : i+5])
		i += 5

		dataStart := i
		dataEnd := dataStart + int(header.SectionLength) - 4 // CRC
		data := p.ParsePSISectionData(tableID, b[dataStart:dataEnd])
		i += int(header.SectionLength) - 4

		syntax.CRC32 = uint32(b[i])<<24 | uint32(b[i+1])<<16 |
			uint32(b[i+2])<<8 | uint32(b[i+3])
		i += 4

		sections = append(sections, &PSISection{
			Header: header,
			Syntax: syntax,
			Data:   data,
		})
	}

	return sections
}


func (p *Parser) ParsePSISectionHeader(b []byte) *PSISectionHeader {
	return &PSISectionHeader{
		TableID:                uint8(b[0]),
		SectionSyntaxIndicator: b[1]>>7&0x1 > 0,
		PrivateBit:             b[1]>>6&0x1 > 0,
		SectionLength:          uint16(b[1]&0xf)<<8 | uint16(b[2]),
	}
}

func (p *Parser) ParsePSISectionSyntax(b []byte) *PSISectionSyntax {
	return &PSISectionSyntax{
		TableIDExtension:     uint16(b[0])<<8 | uint16(b[1]),
		VersionNumber:        uint8(b[2]&0x3f) >> 1,
		CurrentNextIndicator: b[2]&0x1 > 0,
		SectionNumber:        uint8(b[3]),
		LastSectionNumber:    uint8(b[4]),
	}
}

func (p *Parser) ParsePSISectionData(tableID TableID, b []byte) *PSISectionData {
	data := &PSISectionData{
		TableID: tableID,
	}

	switch tableID {
	case PAT:
		data.PAT = p.ParseProgramAssociationTable(b)

	case PMT:
		data.PMT = p.ParseProgramMapTable(b)
	}

	return data
}
