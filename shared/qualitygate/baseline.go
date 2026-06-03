package qualitygate

import "math"

// BaselineStats holds rolling statistics for a metric from historical successful runs.
type BaselineStats struct {
	Samples int
	Mean    float64
	StdDev  float64
	Values  []float64
}

func ComputeBaselineStats(values []float64) BaselineStats {
	if len(values) == 0 {
		return BaselineStats{}
	}
	mean := mean(values)
	if len(values) == 1 {
		return BaselineStats{
			Samples: 1,
			Mean:    mean,
			StdDev:  0,
			Values:  append([]float64(nil), values...),
		}
	}
	return BaselineStats{
		Samples: len(values),
		Mean:    mean,
		StdDev:  stddev(values, mean),
		Values:  append([]float64(nil), values...),
	}
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func stddev(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0
	}
	var sumSq float64
	for _, v := range values {
		d := v - mean
		sumSq += d * d
	}
	return math.Sqrt(sumSq / float64(len(values)-1))
}
