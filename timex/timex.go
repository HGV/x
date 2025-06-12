package timex

import (
	"time"
)

const (
	Day time.Duration = 24 * time.Hour
)

// Deprecated: timex.GenerateSeries is deprecated. Use timex.MakeTimeSeries instead.
func GenerateSeries(start, stop time.Time, interval time.Duration) []time.Time {
	return MakeTimeSeries(start, stop, interval)
}

func MakeTimeSeries(start, stop time.Time, interval time.Duration) []time.Time {
	if start.After(stop) || interval <= 0 {
		return nil
	}

	l := int64(stop.Sub(start) / interval)
	series := make([]time.Time, 0, l)
	for t := start; !t.After(stop); t = t.Add(interval) {
		series = append(series, t)
	}
	return series
}

func MakeDateSeries(start, stop Date, intervalDays int) []Date {
	if start.After(stop) || intervalDays <= 0 {
		return nil
	}

	series := make([]Date, 0, stop.DaysSince(start))
	for d := start; !d.After(stop); d = d.AddDays(intervalDays) {
		series = append(series, d)
	}
	return series
}
