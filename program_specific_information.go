package tsgo

type ProgramSpecificInformation struct {
	PointerField int
	Sections     []*PSISection
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
	TableID int

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

	// TableData -- Could be any one of these options
	PAT *ProgramAssociationTable
}

func ParseProgramSpecificInformation(b []byte) *ProgramSpecificInformation {
	psi := &ProgramSpecificInformation{
		PointerField: int(b[0]),
		Sections:     make([]*PSISection, 0),
	}

	for i := psi.PointerField; i < len(b); {
		header := ParsePSISectionHeader(b[i : i+3])
		i += 3

		if header.TableID == 0xff {
			return psi
		}

		syntax := ParsePSISectionSyntax(b[i : i+5])
		i += 5

		switch header.TableID {
		case 0:
			syntax.PAT = ParseProgramAssociationTable(b[i:])
		}

		psi.Sections = append(psi.Sections, &PSISection{
			Header: header,
			Syntax: syntax,
		})

		i += int(header.SectionLength)
	}

	return psi
}

func ParsePSISectionHeader(b []byte) *PSISectionHeader {
	return &PSISectionHeader{
		TableID:                int(b[0]),
		SectionSyntaxIndicator: b[1]>>7&0x1 > 0,
		PrivateBit:             b[1]>>6&0x1 > 0,
		SectionLength:          uint16(b[1]&0xf)<<8 | uint16(b[2]),
	}
}

func ParsePSISectionSyntax(b []byte) *PSISectionSyntax {
	return &PSISectionSyntax{
		TableIDExtension:     uint16(b[0])<<8 | uint16(b[1]),
		VersionNumber:        uint8(b[2]&0x3f) >> 1,
		CurrentNextIndicator: b[2]&0x1 > 0,
		SectionNumber:        uint8(b[3]),
		LastSectionNumber:    uint8(b[4]),
	}
}
