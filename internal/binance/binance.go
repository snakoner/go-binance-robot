package binance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
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
	// if no close data - get close data
	if len(r.TradingSession.Close) == 0 {
		err := GetHistoricalKlines(r, symbol, numberOfKlines)
		if err != nil {
			log.Fatal(err)
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
			r.TradingSession.Close[len(r.TradingSession.Close)-1] = curr

			if tm != r.TradingSession.LastTime {
				price, _ := GetPenultimatePrice(r, symbol)
				r.TradingSession.Close[len(r.TradingSession.Close)-1] = price
				r.TradingSession.Close = append(r.TradingSession.Close, curr)
				r.TradingSession.Close = r.TradingSession.Close[1:]
			}
			r.TradingSession.LastTime = tm
			result, _ := r.StrategyFunc(r.TradingSession.Close)
			result = true

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

	// fmt.Println(r.TradingSession.Active)
	// fmt.Println(r.TradingSession.BuyValue)
	// fmt.Println(r.TradingSession.OpenPrice)
	// fmt.Println(r.TradingSession.StopLossValue)
	// fmt.Println(r.TradingSession.TakeProfitValue)

	// main trading session
	if r.TradingSession.Active {
		for r.TradingSession.Active {
			prev = curr
			// fmt.Println("Waiting for change 0")
			curr = <-close
			<-t
			// fmt.Println("Waiting for change 1")
			priceChange := r.TradingSession.Quantity*curr - r.TradingSession.BuyValue
			fmt.Printf("Price change from beginning: %.4f\n", priceChange)
			if curr <= r.TradingSession.StopLossValue {
				// [todo] sell with loss
				r.TradingSession.Active = false
				r.TradingSession.Result.EndTime = time.Now()
				r.TradingSession.Result.Profit = priceChange
				log.Println("Deal is closed with loss: ", priceChange)
				break
			} else if curr >= r.TradingSession.TakeProfitValue { // [todo] sell with profit
				r.TradingSession.Active = false
				r.TradingSession.Result.EndTime = time.Now()
				r.TradingSession.Result.Profit = priceChange
				log.Println("Deal is closed with profit: ", priceChange)
				break
			}

			diff := curr - r.TradingSession.LastPriceForSLChange
			prev = curr
			fmt.Println("Price: ", curr)
			fmt.Printf("Diff: %.5f\n", diff)

			if trailing.RecalcTrailingStop(r, diff) {
				// change OCO order
				fmt.Printf("SL: %f, TP:%f\n", r.TradingSession.StopLossValue, r.TradingSession.TakeProfitValue)
				r.TradingSession.LastPriceForSLChange = curr
			}
			fmt.Println()
		}

	}
}
