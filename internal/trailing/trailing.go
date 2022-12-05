package trailing

import "github.com/go-binance-robot/internal/robot"

func RecalcTrailingStop(r *robot.Robot, diff float64) bool {
	if diff > 0 {
		r.TradingSession.StopLossValue += diff
		r.TradingSession.TakeProfitValue += diff
		return true
	}

	return false
}
