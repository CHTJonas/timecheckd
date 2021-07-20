package ntp

import "time"

var (
	nanoPerSec = uint64(time.Second)
	ntpEpoch   = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

func parseDuration(t uint64) time.Duration {
	sec := (t >> 32) * nanoPerSec
	frac := (t & 0xffffffff) * nanoPerSec
	nsec := frac >> 32
	if uint32(frac) >= 0x80000000 {
		nsec++
	}
	return time.Duration(sec + nsec)
}

func parseTime(t uint64) time.Time {
	return ntpEpoch.Add(parseDuration(t))
}
