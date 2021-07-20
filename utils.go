package main

import (
	"github.com/CHTJonas/timecheckd/ntp"
	"github.com/cloudflare/roughtime/mjd"
)

// Number of Gregorian calendar days between 1858-11-17 and 1970-01-01.
const daysPreEpoch = 40587
const secsPerDay = 24 * 60 * 60

func getMjd() mjd.Mjd {
	target := "ntp0a.cl.cam.ac.uk"
	t, err := ntp.GetTime(target)
	if err != nil {
		panic(err)
	}
	daysPostEpoch := t.Unix() / secsPerDay
	days := uint64(daysPreEpoch + daysPostEpoch)
	micros := float64((t.Unix()-daysPostEpoch*secsPerDay)*1e6) + float64(t.Nanosecond())/1000
	return mjd.New(days, micros)
}
