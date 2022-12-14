package indicators

// Simple moving average
func Sma(data []float64, n int) []float64 {
	result := make([]float64, len(data) - n + 1)
	for i := 0; i < len(result); i++ {
		for _, val := range data[i:i+n] {
			result[i] += val
		}
		result[i] /= float64(n)
	}
	
	return result
}

// Exponential moving average
func Ema(data []float64, n int) []float64 {
	result := make([]float64, len(data))
	result[0] = data[0]
	alpha := 2./(1. + float64(n))

	for i := 1; i < len(data); i++ {
		result[i] = alpha * data[i] + (1. - alpha) *result[i - 1]
	}

	return result
}
