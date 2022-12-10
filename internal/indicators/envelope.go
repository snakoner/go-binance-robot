package indicators

import (
	"math"
)

const (
	envelopeLen  = 500
	envelopeH    = 8
	envelopeMult = 3
)

func Envelope(data []float64) (bool, bool) {
	if len(data) >= envelopeLen {
		data = data[len(data)-envelopeLen:]
	} else {
		return false, false
	}

	up := make([]float64, len(data))
	down := make([]float64, len(data))
	y := []float64{}

	sum_e := 0.0
	for i := 0; i < envelopeLen; i++ {
		summ := 0.0
		summw := 0.0
		for j := 0; j < envelopeLen; j++ {
			w := math.Exp((-(math.Pow(float64(i-j), 2))) / float64(envelopeH*envelopeH*2))
			summ += data[j] * w
			summw += w
		}

		y2 := summ / summw
		sum_e += math.Abs(data[i] - y2)
		y = append(y, y2)
	}

	mae := float64(envelopeMult) * sum_e / float64(envelopeLen)
	for i := 0; i < len(y); i++ {
		up[i] = y[i] + mae
		down[i] = y[i] - mae
	}

	long := data[len(data)-1] <= down[len(data)-1]
	short := data[len(data)-1] >= up[len(data)-1]

	return long, short
}
