package timeHelper

import (
	"time"

	"../log"
)

// GetCurrentTime is a convience method to get milliseconds since epoch.
func GetCurrentTime() int64 {
	return time.Now().Unix()
}

// FewDaysLater return time that days after
func FewDaysLater(day int) time.Time {
	return FewDurationLater(time.Duration(day) * 24 * time.Hour)
}

// TwentyFourHoursLater return time that 24 hrs after
func TwentyFourHoursLater() time.Time {
	return FewDurationLater(time.Duration(24) * time.Hour)
}

// SixHoursLater return time that 6 hrs after
func SixHoursLater() time.Time {
	return FewDurationLater(time.Duration(6) * time.Hour)
}

// FewDurationLater return time that hrs after
func FewDurationLater(duration time.Duration) time.Time {
	// When Save time should considering UTC
	baseTime := time.Now()
	log.Debugf("basetime : %s", baseTime)
	fewDurationLater := baseTime.Add(duration)
	log.Debugf("time : %s", fewDurationLater)
	return fewDurationLater
}

// FewDurationLaterFrom return time that hrs after
func FewDurationLaterFrom(baseTime time.Time, duration time.Duration) time.Time {
	// When Save time should considering UTC
	log.Debugf("basetime : %s", baseTime)
	fewDurationLater := baseTime.Add(duration)
	log.Debugf("time : %s", fewDurationLater)
	return fewDurationLater
}

// ShortDate returns only date
func ShortDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Now().Location())
}

// ShortHour returns till date, hr
func ShortHour(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 0, 0, 0, time.Now().Location())
}

// ScheduleDate returns date and hour with schedule info
func ScheduleDate(date time.Time, tim time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), tim.Hour(), tim.Minute(), tim.Second(), 0, time.Now().Location())
}

// FewDaysLaterDate return time that days after
func FewDaysLaterDate(date time.Time, day int) time.Time {
	fewDurationLater := date.Add(time.Duration(day) * 24 * time.Hour)
	return fewDurationLater
}

// DateFromTime returns only date from time
func DateFromTime(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// IsExpired return expire result
func IsExpired(expirationTime time.Time) bool {
	baseTime := time.Now()
	log.Debugf("basetime : %s", baseTime)
	log.Debugf("expirationTime : %s", expirationTime)
	elapsed := time.Since(expirationTime)
	log.Debugf("elapsed : %s", elapsed)
	after := time.Now().After(expirationTime)
	log.Debugf("after : %t", after)
	return after
}
