package indicators

import (
	// "fmt"
	"log"

	"github.com/go-binance-robot/internal/errno"
)

const (
	macdS = 12
	macdL = 26
	macdA = 9
)

// Difference between two slices of type float64
func slicesDiff(a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return []float64{}, errno.ErrLengthNotEqual
	}

	diff := make([]float64, len(a))

	for i := 0; i < len(a); i++ {
		diff[i] = a[i] - b[i]
	}

	return diff, nil
}

// Macd histogram
func Macd(data []float64) ([]float64, error) {
	macdRet, err := slicesDiff(Ema(data, macdS), Ema(data, macdL))
	if err != nil {
		log.Fatal(err)
		return macdRet, err
	}

	signal := Ema(macdRet, macdA)
	macdHist, err := slicesDiff(macdRet, signal)
	if err != nil {
		log.Fatal(err)
		return macdHist, err
	}

	return macdHist, nil
}
