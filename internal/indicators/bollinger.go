package indicators

const (
  bollingerN = 20
  bollingerStd = 2.0
)


func bollingerBands(data []float64) ([]float64, []float64) {
  sma := Sma(data, bollingerN)
  std := make([]float64, len(data) - bollingerN + 1)
  for i := 0; i < len(data) - bollingerN + 1; i++ {
    sum := 0.0
    for j := i; j < i + bollingerN; j++ {
      sum += (data[j] - sma[i]) * (data[j] - sma[i])
    }
    sum /= bollingerN
    std[i] = bollingerStd * math.Sqrt(sum)
  }
  
  low := make([]float64, len(data) - bollingerN + 1)
  high := make([]float64, len(data) - bollingerN + 1)
  
  for i := 0; i < len(data) - bollingerN + 1 {
    low[i] = sma[i] - std[i]
    high[i] = sma[i] + std[i]
  }
  
  return low, high
}

func Bollinger(data []float64) (bool, bool) {
  low, high := bollingerBands(data)
  lastLow := low[len(low) - 1]
  lastHigh := high[len(high) - 1]
  
  return data[len(data) - 1] < lastLow, data[len(data) - 1] > lastHigh
}
