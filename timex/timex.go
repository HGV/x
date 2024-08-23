package timex

import (
	"time"
)

const (
	Day time.Duration = 24 * time.Hour
)

func GenerateSeries(start, stop time.Time, interval time.Duration) []time.Time {
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
