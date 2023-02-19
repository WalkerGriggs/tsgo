package tsgo

type ProgramAssociationTable struct {
	// ProgramNum is used to specify the table ID extension for the associated
	// PMT. The default value of 0 is reserved for a NIT packet.
	//ProgramNum uint16

	// ProgramMapPID is used to specify the packet ID which contains the program
	// map table.
	// ProgramMapPID uint16
	ProgramMap map[uint16]uint16
}

func ParseProgramAssociationTable(b []byte) *ProgramAssociationTable {
	t := &ProgramAssociationTable {
		ProgramMap: make(map[uint16]uint16),
	}

	for i := 0; i < len(b) / 4; i += 4 {
		num := uint16(b[0])<<8 | uint16(b[1])
		pid := uint16(b[2]&0x1f)<<8 | uint16(b[3])
		t.ProgramMap[num] = pid
	}

	return t
}
