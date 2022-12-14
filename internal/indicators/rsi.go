package indicators

const (
	rsiLen = 14
	rsiLongLimit = 30.0
	rsiShortLimit = 70.0
)

func Rsi(data []float64) (bool, bool) {
	u := make([]float64, len(data) - 1)
	d := make([]float64, len(data) - 1)
	r := make([]float64, len(data) - 1)

	for i := 0; i < len(data) - 1; i++ {
		if data[i + 1] >= data[i] {
			u[i] = data[i + 1] - data[i]
		} else {
			d[i] = data[i] - data[i + 1]
		}
	}

	emaU := Ema(u, rsiLen)
	emaD := Ema(d, rsiLen)

	for i := 0; i < len(data) - 1; i++ {
		r[i] = 100.0
		if emaD[i] != 0.0 {
			r[i] = 100.0 - 100.0 / (1 + emaU[i] / emaD[i])
		}
	}

	lastValue := r[len(r) - 1]

	return lastValue < rsiLongLimit, lastValue > rsiShortLimit
}
