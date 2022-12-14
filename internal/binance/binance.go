package binance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/go-binance-robot/internal/robot"
	"github.com/go-binance-robot/pkg/logger"
)

// Get N = numberOfKlines historical klines
func GetHistoricalKlines(ts *robot.Trade, numberOfKlines int) error {
	client := binance.NewClient(ts.Root.BinanceApiPublic, ts.Root.BinanceApiPrivate)
	klines, err := client.NewKlinesService().Symbol(ts.Token + ts.Root.StableCurrency).
		Interval(ts.Root.Timeframe).Limit(numberOfKlines).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil
	}

	klinesFloat := make([]float64, numberOfKlines)
	for i := 0; i < numberOfKlines; i++ {
		klinesFloat[i], _ = strconv.ParseFloat(klines[i].Close, 64)
	}

	ts.Close = klinesFloat
	ts.LastTime = time.Unix(klines[len(klines)-1].OpenTime/1000, 0)

	return nil
}

// Get penultimate close price
func GetPenultimatePrice(ts *robot.Trade) (float64, error) {
	client := binance.NewClient(ts.Root.BinanceApiPublic, ts.Root.BinanceApiPrivate)
	klines, err := client.NewKlinesService().Symbol(ts.Token + ts.Root.StableCurrency).
		Interval(ts.Root.Timeframe).Limit(2).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return 0.0, err
	}

	price, _ := strconv.ParseFloat(klines[0].Close, 64)
	return price, nil
}

// Get last close price
func GetLastPrice(r *robot.Robot, symbol string) (float64, error) {
	client := binance.NewClient(r.BinanceApiPublic, r.BinanceApiPrivate)
	klines, err := client.NewKlinesService().Symbol(symbol + r.StableCurrency).
		Interval(r.Timeframe).Limit(1).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return 0.0, err
	}

	price, _ := strconv.ParseFloat(klines[0].Close, 64)
	return price, nil
}

// Goroutine function to return time and close price for token
func WebSocketTracking(ts *robot.Trade, close chan float64, t chan time.Time) {
	fmt.Printf("Web socket is running: %s\n", ts.Token)
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		currClose, _ := strconv.ParseFloat(event.Kline.Close, 64)
		close <- currClose
		t <- time.Unix(event.Kline.StartTime/1000, 0)
	}

	errHandler := func(err error) {
		log.Println(err)
		return
	}

	doneC, _, err := binance.WsKlineServe(ts.Token+ts.Root.StableCurrency, ts.Root.Timeframe, wsKlineHandler, errHandler)

	if err != nil {
		log.Fatal(err)
		return
	}
	<-doneC
}

// Function to run from main
func WebSocketRun(ts *robot.Trade, numberOfKlines int) {
	// logger
	var s string
	logger := logger.New(fmt.Sprintf("../../log/envelope_%s.log", ts.Token))
	logger.Open()
	defer logger.Close()
	profits := []float64{}
	// get data
	if err := GetHistoricalKlines(ts, numberOfKlines); err != nil {
		log.Fatal(err)
		return
	}

	// web socket tracking
	close := make(chan float64)
	t := make(chan time.Time)
	go WebSocketTracking(ts, close, t)

	// get price, timestamp from WebSocketTracking goroutine
	var prev, curr float64
	var tm time.Time

	for {
		curr = <-close
		tm = <-t
		if curr != prev {
			ts.Close[len(ts.Close)-1] = curr

			if tm != ts.LastTime {
				price, _ := GetPenultimatePrice(ts)
				ts.Close[len(ts.Close)-1] = price
				ts.Close = append(ts.Close, curr)
				ts.Close = ts.Close[1:]
			}

			ts.LastTime = tm
			ts.Strategy.Apply(ts.Close)
			result := ts.Strategy.IsLong()
			if result {
				s = "Envelope long OK"
				logger.Write(s)
			}

			if result && ts.Active == false {
				ts.Active = true
				ts.BuyValue = ts.Root.StartBalance // [todo]
				ts.OpenPrice, _ = GetLastPrice(ts.Root, ts.Token)
				ts.StopLossValue = (1 - ts.Root.StopLoss/100.0) * ts.OpenPrice
				ts.TakeProfitValue = (1 + ts.Root.TakeProfit/100.0) * ts.OpenPrice
				ts.LastPriceForSLChange = ts.OpenPrice
				ts.Quantity = ts.BuyValue / ts.OpenPrice
				ts.Result.StartTime = time.Now()
				s = fmt.Sprintf("Price: %f", ts.OpenPrice)
				logger.Write(s)
				fmt.Printf("Strategy started for %s\n", ts.Token)

				log.Printf("Strategy started for %s\n", ts.Token)

				break
			}
		}
		prev = curr
	}

	// main trading session
	if ts.Active {
		for ts.Active {
			prev = curr
			curr = <-close
			<-t
			priceChange := ts.Quantity*curr - ts.BuyValue
			if curr != prev {
				profit := (curr - ts.OpenPrice) / ts.OpenPrice * 100.0
				s = fmt.Sprintf("Profit: %.3f o/o", profit)
				logger.Write(s)
				profits = append(profits, profit)
			}

			if curr <= ts.StopLossValue || curr >= ts.TakeProfitValue {
				// [todo] sell
				ts.Active = false
				ts.Result.EndTime = time.Now()
				ts.Result.ProfitPrice = priceChange
				ts.Result.ProfitPerc = priceChange / ts.OpenPrice

				// logging
				s = fmt.Sprint("\n##########\nClose deal\n###########\n")
				logger.Write(s)
				s = fmt.Sprintf("SL: %f, TP:%f", ts.StopLossValue, ts.TakeProfitValue)
				logger.Write(s)
				s = fmt.Sprintf("Price sell: %f", curr)
				logger.Write(s)
				if curr <= ts.StopLossValue {
					s = fmt.Sprintf("Deal is closed due SL: %f", priceChange)
				} else {
					s = fmt.Sprintf("Deal is closed due TP: %f", priceChange)
				}
				logger.Write(s)

				break
			}

			diff := curr - ts.LastPriceForSLChange
			prev = curr

			if recalcTrailingStop(ts, diff) {
				// change OCO order
				s = fmt.Sprintf("SL: %f, TP:%f", ts.StopLossValue, ts.TakeProfitValue)
				logger.Write(s)
				ts.LastPriceForSLChange = curr
			}
		}

	}
	if len(profits) != 0 {
		s = fmt.Sprintf("Max profit: %f", maxFloat(profits))
		logger.Write(s)
	}
}

func maxFloat(s []float64) float64 {
	max := s[0]

	for _, val := range s {
		if max < val {
			max = val
		}
	}
	return max
}

func recalcTrailingStop(ts *robot.Trade, diff float64) bool {
	if diff > 0 {
		ts.StopLossValue += diff
		ts.TakeProfitValue += diff
		return true
	}

	return false
}
