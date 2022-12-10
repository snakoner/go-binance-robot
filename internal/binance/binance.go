package binance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/go-binance-robot/internal/logger"
	"github.com/go-binance-robot/internal/robot"
	"github.com/go-binance-robot/internal/trailing"
)

// Get N = numberOfKlines historical klines
func GetHistoricalKlines(r *robot.Robot, symbol string, numberOfKlines int) error {
	client := binance.NewClient(r.BinanceApiPublic, r.BinanceApiPrivate)
	klines, err := client.NewKlinesService().Symbol(symbol + r.StableCurrency).
		Interval(r.Timeframe).Limit(numberOfKlines).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil
	}

	klinesFloat := make([]float64, numberOfKlines)
	for i := 0; i < numberOfKlines; i++ {
		klinesFloat[i], _ = strconv.ParseFloat(klines[i].Close, 64)
	}

	r.TradingSession.Close = klinesFloat
	r.TradingSession.LastTime = time.Unix(klines[len(klines)-1].OpenTime/1000, 0)

	return nil
}

// Get penultimate close price
func GetPenultimatePrice(r *robot.Robot, symbol string) (float64, error) {
	client := binance.NewClient(r.BinanceApiPublic, r.BinanceApiPrivate)
	klines, err := client.NewKlinesService().Symbol(symbol + r.StableCurrency).
		Interval(r.Timeframe).Limit(2).Do(context.Background())
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
func WebSocketTracking(r *robot.Robot, symbol string, close chan float64, t chan time.Time) {
	fmt.Println("Web socket is running")
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		currClose, _ := strconv.ParseFloat(event.Kline.Close, 64)
		close <- currClose
		t <- time.Unix(event.Kline.StartTime/1000, 0)
	}

	errHandler := func(err error) {
		log.Println(err)
		return
	}

	doneC, _, err := binance.WsKlineServe(symbol+r.StableCurrency, r.Timeframe, wsKlineHandler, errHandler)

	if err != nil {
		log.Fatal(err)
		return
	}
	<-doneC
}

// Function to run from main
func WebSocketRun(r *robot.Robot, symbol string, numberOfKlines int) {
	// for logger
	var s string
	logger := logger.New(fmt.Sprintf("../../envelope_%s.log", symbol))
	logger.Open()
	// get data
	if err := GetHistoricalKlines(r, symbol, numberOfKlines); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Hello")
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
			r.TradingSession.Close[len(r.TradingSession.Close)-1] = curr

			if tm != r.TradingSession.LastTime {
				price, _ := GetPenultimatePrice(r, symbol)
				r.TradingSession.Close[len(r.TradingSession.Close)-1] = price
				r.TradingSession.Close = append(r.TradingSession.Close, curr)
				r.TradingSession.Close = r.TradingSession.Close[1:]
			}

			fmt.Println(prev, curr)
			r.TradingSession.LastTime = tm
			r.Strategy.Apply(r.TradingSession.Close)
			result := r.Strategy.IsLong()
			fmt.Println(result)
			if result {
				s = "Envelope long OK"
				logger.Write(s)
			}

			// result = true

			// log.Println(r.TradingSession.Close[len(r.TradingSession.Close)-2:], "long: ", result)

			if result && r.TradingSession.Active == false {
				// [todo] buy, set sl and tp and set fields
				r.TradingSession.Active = true
				r.TradingSession.Token = symbol
				r.TradingSession.BuyValue = r.StartBalance // [todo]
				r.TradingSession.OpenPrice, _ = GetLastPrice(r, symbol)
				r.TradingSession.StopLossValue = (1 - r.StopLoss/100.0) * r.TradingSession.OpenPrice
				r.TradingSession.TakeProfitValue = (1 + r.TakeProfit/100.0) * r.TradingSession.OpenPrice
				r.TradingSession.Result.StartTime = time.Now()
				r.TradingSession.LastPriceForSLChange = r.TradingSession.OpenPrice
				r.TradingSession.Quantity = r.TradingSession.BuyValue / r.TradingSession.OpenPrice
				break
			}
		}
		prev = curr
	}

	// main trading session
	if r.TradingSession.Active {
		for r.TradingSession.Active {
			prev = curr
			// fmt.Println("Waiting for change 0")
			curr = <-close
			<-t
			// fmt.Println("Waiting for change 1")
			priceChange := r.TradingSession.Quantity*curr - r.TradingSession.BuyValue
			// s = fmt.Sprintf("Price change from beginning: %.4f", priceChange)
			// logger.Write(s)
			if curr <= r.TradingSession.StopLossValue {
				// [todo] sell with loss
				r.TradingSession.Active = false
				r.TradingSession.Result.EndTime = time.Now()
				r.TradingSession.Result.Profit = priceChange
				s = fmt.Sprint("##########\nClose deal\n###########")
				logger.Write(s)

				s = fmt.Sprintf("SL: %f, TP:%f", r.TradingSession.StopLossValue, r.TradingSession.TakeProfitValue)
				logger.Write(s)

				s = fmt.Sprintf("Deal is closed due SL: %f", priceChange)
				logger.Write(s)

				break
			} else if curr >= r.TradingSession.TakeProfitValue { // [todo] sell with profit
				r.TradingSession.Active = false
				r.TradingSession.Result.EndTime = time.Now()
				r.TradingSession.Result.Profit = priceChange
				s = fmt.Sprint("##########\nClose deal\n###########")
				logger.Write(s)

				s = fmt.Sprintf("SL: %f, TP:%f", r.TradingSession.StopLossValue, r.TradingSession.TakeProfitValue)
				logger.Write(s)

				s = fmt.Sprintf("Deal is closed due TP: %f", priceChange)
				logger.Write(s)

				break
			}

			diff := curr - r.TradingSession.LastPriceForSLChange
			prev = curr
			s = fmt.Sprintf("Price: %f", curr)
			logger.Write(s)

			// s = fmt.Sprintf("Diff: %.5f", diff)
			// logger.Write(s)

			if trailing.RecalcTrailingStop(r, diff) {
				// change OCO order
				s = fmt.Sprintf("SL: %f, TP:%f", r.TradingSession.StopLossValue, r.TradingSession.TakeProfitValue)
				logger.Write(s)
				r.TradingSession.LastPriceForSLChange = curr
			}
			fmt.Println()
		}

	}
}
