package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/go-binance-robot/internal/indicators"
	"github.com/go-binance-robot/internal/robot"
)

func GetHistoricalKlines(r *robot.Robot, symbol string, numberOfKlines int) error {
	client := binance.NewClient(r.BinanceApiPublic, r.BinanceApiPrivate)
	klines, err := client.NewKlinesService().Symbol(symbol + r.StableCurrency).
		Interval(r.Timeframe).Limit(numberOfKlines).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	klinesFloat := make([]float64, numberOfKlines)
	for i := 0; i < numberOfKlines; i++ {
		klinesFloat[i], _ = strconv.ParseFloat(klines[i].Close, 64)
	}

	r.ActiveTrade.Close = klinesFloat
	r.ActiveTrade.LastTime = time.Unix(klines[len(klines)-1].OpenTime/1000, 0)

	return nil
}

func WebSocketTracking(r *robot.Robot, symbol string, close chan float64, t chan time.Time) {

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		currClose, _ := strconv.ParseFloat(event.Kline.Close, 64)
		close <- currClose
		t <- time.Unix(event.Kline.StartTime/1000, 0)
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneC, _, err := binance.WsKlineServe(symbol+r.StableCurrency, r.Timeframe, wsKlineHandler, errHandler)

	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}

func WebSocketRun(r *robot.Robot, symbol string, numberOfKlines int) {
	// if no close data - get close data
	if len(r.ActiveTrade.Close) == 0 {
		err := GetHistoricalKlines(r, symbol, numberOfKlines)
		if err != nil {
			return
		}
	}
	// web socket tracking
	close := make(chan float64)
	t := make(chan time.Time)
	go WebSocketTracking(r, symbol, close, t)

	var prev, curr float64
	var tm time.Time
	for {
		curr = <-close
		tm = <-t
		if curr != prev {
			r.ActiveTrade.Close[len(r.ActiveTrade.Close)-1] = curr
			r.ActiveTrade.CurrentPrice = curr

			if tm != r.ActiveTrade.LastTime {
				r.ActiveTrade.Close = append(r.ActiveTrade.Close, curr)
				r.ActiveTrade.Close = r.ActiveTrade.Close[1:]
				r.ActiveTrade.LastTime = tm
			}
			fmt.Printf("Last time: %v, Time: %v\n", r.ActiveTrade.LastTime, tm)
			fmt.Println("Cur price: ", curr)
			fmt.Printf("Len: %d\n", len(r.ActiveTrade.Close))
			// estimator calculute
			result, err := indicators.Envelope(r.ActiveTrade.Close, true)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Result of long: ", result)
		}
		prev = curr
	}
}
