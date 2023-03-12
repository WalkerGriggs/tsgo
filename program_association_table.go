package tsgo

type ProgramAssociationTable struct {
	// ProgramMapPID is used to specify the packet ID which contains the program
	// map table.
	//
	//  - The key is the program_num, which is used to specify the table ID extension
	//    for the associated PMT. The default value of 0 is reserved for a NIT
	//    packet.
	//  - The value is the program map PID which contains the PMT.
	//
	// For example, 1:4096 means the packet ID 4096 with program specific table ID
	// 1 contains the program map table.
	ProgramMap map[uint16]uint16
}

func (p *Parser) ParseProgramAssociationTable(l int) *ProgramAssociationTable {
	t := &ProgramAssociationTable{
		ProgramMap: make(map[uint16]uint16),
	}

	for i := 0; i < l/4; i += 4 {
		bs := p.ReadBytes(4)
		num := uint16(bs[0])<<8 | uint16(bs[1])
		pid := uint16(bs[2]&0x1f)<<8 | uint16(bs[3])
		t.ProgramMap[pid] = num // local, PAT map
		p.ProgramMap[pid] = num // global, parser map
	}

	return t
}
