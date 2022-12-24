package trade

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

const (
	decreaseQuanitiyPercent = 20.0
)

type FuturesOrder struct {
	Order       *futures.CreateOrderResponse
	OrderT      futures.OrderType // buy or sell
	SideT       futures.SideType  // market, limit ..
	TimeInForce futures.TimeInForceType
	Quantity    float64
	Leverage    string
}

// func GetSymbolInfo()

func roundFloat(val float64, precision float64) float64 {

	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func getExchangeInfo(ts *Trade) *futures.ExchangeInfo {
	futuresClient := binance.NewFuturesClient(ts.Root.BinanceApiPublic, ts.Root.BinanceApiPrivate) // USDT-M Futures
	exInfo, err := futuresClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return exInfo
}

func GetPrecision(ts *Trade) error {
	exchangeInfo := getExchangeInfo(ts)
	if exchangeInfo == nil {
		log.Fatal("Cant get precision")
		return errors.New("Cant get exchange info")
	}

	for _, s := range exchangeInfo.Symbols {
		if s.Symbol == ts.Token+ts.Root.StableCurrency {
			symbolInfo := &s
			stepSize, _ := strconv.ParseFloat(symbolInfo.LotSizeFilter().StepSize, 64)
			ts.Precision = -math.Log10(stepSize)
			break
		}
	}

	return nil
}

func CreateFuturesOrder(ts *Trade) error {
	fOrder := &FuturesOrder{
		OrderT:      futures.OrderTypeMarket,
		SideT:       futures.SideTypeBuy,
		TimeInForce: futures.TimeInForceTypeGTC,
		// Quantity:    ,
	}

	futuresClient := binance.NewFuturesClient(ts.Root.BinanceApiPublic, ts.Root.BinanceApiPrivate) // USDT-M Futures
	quantity := ts.BuyValue * (1 - 0.1) / ts.OpenPrice
	quantity = roundFloat(quantity, ts.Precision)
	fOrder.Quantity = quantity

	order, err := futuresClient.NewCreateOrderService().Symbol(ts.Token + ts.Root.StableCurrency).
		Side(fOrder.SideT).Type(fOrder.OrderT).
		Quantity(fmt.Sprintf("%f", quantity)).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return err
	}

	fOrder.Order = order

	return nil
}

func GetFuturesOrder(ts *Trade) (*futures.Order, error) {
	if ts.FuturesOrder == nil {
		return nil, errors.New("No active futures order for Trading session")
	}

	futuresClient := binance.NewFuturesClient(ts.Root.BinanceApiPublic, ts.Root.BinanceApiPrivate) // USDT-M Futures
	order, err := futuresClient.NewGetOrderService().Symbol(ts.Token + ts.Root.StableCurrency).
		OrderID(ts.FuturesOrder.Order.OrderID).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return order, nil
}
