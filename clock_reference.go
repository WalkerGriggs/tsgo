package tsgo

import (
	"time"
)

type ClockReference struct {
	Base      int64
	Extension int64
	Duration  time.Duration
}

func ParseClockReference(b []byte) *ClockReference {
	pcr := uint64(b[0])<<40 |
		uint64(b[1])<<32 |
		uint64(b[2])<<24 |
		uint64(b[3])<<16 |
		uint64(b[4])<<8 |
		uint64(b[5])

	r := &ClockReference{
		Base:      int64(pcr >> 15),
		Extension: int64(pcr & 0x1ff),
	}

	x := time.Duration((r.Base * 1e9) / 90_000)
	y := time.Duration((r.Extension * 1e9) / 27_000_000)
	r.Duration = x + y

	return r
}
