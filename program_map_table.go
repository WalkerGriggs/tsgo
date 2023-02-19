package tsgo

type ProgramMapTable struct {
	ProgramClockReferencePID uint16
	ProgramInfoLength        uint16
	ProgramDescriptors       []byte
	ElementaryStreamInfoData []*ElementaryStreamInfo
}

type ElementaryStreamInfo struct {
	StreamType             uint8
	ElementaryPID          uint16
	ElementaryStreamLength uint16
	// TODO ElementaryStreamDescriptors
}
