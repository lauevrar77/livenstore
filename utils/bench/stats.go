package bench

import (
	"slices"
)

type Stats struct {
	FirstQuartile int64
	Median        int64
	Mean          int64
	ThirdQuartile int64
}

func ComputeStats(times []int64) Stats {
	var sum int64
	slices.Sort(times)
	for _, t := range times {
		sum += t
	}
	return Stats{
		FirstQuartile: times[len(times)/4],
		Median:        times[len(times)/2],
		Mean:          sum / int64(len(times)),
		ThirdQuartile: times[3*len(times)/4],
	}
}
