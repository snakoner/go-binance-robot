package robot

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-binance-robot/internal/indicators"
	"github.com/joho/godotenv"
)

type Robot struct {
	Name              string   /* name of robot */
	BinanceApiPrivate string   /* binance api_secret_key */
	BinanceApiPublic  string   /* binance api_public_key */
	StableCurrency    string   /* stable_currency : USDT or BUSD */
	Tokens            []string /* tokens : ETH,BTC,... */
	Strategy          string   /* strategy : divergence, rsi, envelope */
	CheckLong         bool     /* open long position */
	CheckShort        bool     /* open short position */
	Timeframe         string
	IsFutures         bool    /* futures of spot trading */
	StartBalance      float64 /* value to buy on */
	TakeProfit        float64 /* take profit value in percent */
	StopLoss          float64 /* stop loss value in percent */
	StrategyFunc      func([]float64) (bool, bool)
	TradingSession    Trade
}

type Trade struct {
	Active    bool      /* deal is active or not */
	Token     string    /* token name */
	BuyValue  float64   /* value bought on */
	OpenPrice float64   /* deal open price */
	LastTime  time.Time /* time of last price */
	Close     []float64 /* actual candle close counts */
	Result    TradingResult
}

type TradingResult struct {
	Profit    float64
	StartTime time.Time
	EndTime   time.Time
}

func New() *Robot {
	robot := new(Robot)
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Cant load .env\n")
		return nil
	}

	robot.Name = os.Getenv("name")
	robot.BinanceApiPrivate = os.Getenv("binance_api_private")
	robot.BinanceApiPublic = os.Getenv("binance_api_public")
	robot.StableCurrency = os.Getenv("stable_currency")
	robot.Tokens = strings.Split(os.Getenv("tokens"), ",")
	robot.Strategy = os.Getenv("strategy")
	robot.CheckLong, _ = strconv.ParseBool(os.Getenv("check_long"))
	robot.CheckShort, _ = strconv.ParseBool(os.Getenv("check_short"))
	robot.Timeframe = os.Getenv("timeframe")
	robot.IsFutures, _ = strconv.ParseBool(os.Getenv("futures"))
	robot.StartBalance, _ = strconv.ParseFloat(os.Getenv("start_balance"), 64)
	robot.TakeProfit, _ = strconv.ParseFloat(os.Getenv("take_profit"), 64)
	robot.StopLoss, _ = strconv.ParseFloat(os.Getenv("stop_loss"), 64)

	switch robot.Strategy {
	case "envelope":
		robot.StrategyFunc = indicators.Envelope
		break
	case "divergence":
		robot.StrategyFunc = indicators.Divergence
		break
	default:
		log.Fatal("Unknown strategy. Please, check .env file")
		return nil
	}

	return robot
}

func (r *Robot) Print() {
	fmt.Println(strings.ToUpper("name: "), r.Name)                              /* name of robot */
	fmt.Println(strings.ToUpper("api_private: "), r.BinanceApiPrivate[:10]+"*") /* binance api_secret_key */
	fmt.Println(strings.ToUpper("api_public: "), r.BinanceApiPublic[:10]+"*")   /* binance api_public_key */
	fmt.Println(strings.ToUpper("stable_currency: "), r.StableCurrency)         /* stable_currency : USDT or BUSD */
	fmt.Println(strings.ToUpper("tokens: "), r.Tokens)                          /* tokens : ETH,BTC,... */
	fmt.Println(strings.ToUpper("strategy: "), r.Strategy)                      /* strategy : divergence, rsi, envelope */
	fmt.Println(strings.ToUpper("long: "), r.CheckLong)                         /* open long position */
	fmt.Println(strings.ToUpper("short: "), r.CheckShort)
	fmt.Println(strings.ToUpper("futures: "), r.IsFutures)
	fmt.Println(strings.ToUpper("start_balance: "), r.StartBalance)
	fmt.Println(strings.ToUpper("stop_loss: "), r.StopLoss)
	fmt.Println(strings.ToUpper("take_profit: "), r.TakeProfit)
	fmt.Println()
}
