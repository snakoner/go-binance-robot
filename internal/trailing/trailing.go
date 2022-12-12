package trailing

import "github.com/go-binance-robot/internal/robot"

func RecalcTrailingStop(ts *robot.Trade, diff float64) bool {
	if diff > 0.0 {
		ts.StopLossValue += diff
		ts.TakeProfitValue += diff
		return true
	}

	return false
}
