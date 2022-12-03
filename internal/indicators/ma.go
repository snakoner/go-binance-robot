package indicators

// Exponential moving average
func Ema(data []float64, n int) []float64 {
	emaRet := make([]float64, len(data))
	emaRet[0] = data[0]
	alpha := 2./(1. + float64(n))

	for i := 1; i < len(data); i++ {
		emaRet[i] = alpha * data[i] + (1. - alpha) *emaRet[i - 1]
	}

	return emaRet
}
