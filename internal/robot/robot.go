package robot

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	ActiveTrade       Trade
}

type Trade struct {
	Active       bool    /* deal is active or not */
	Token        string  /* token name */
	BuyValue     float64 /* value bought on */
	OpenPrice    float64 /* deal open price */
	CurrentPrice float64 /* [deprecated] current price of token in stablecoin */
	LastTime     time.Time
	Close        []float64 /* actual candle close counts */
}

type Roboter interface {
	OpenTrade() error
	CloseTrade() error
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

	return robot
}

func (r *Robot) Print() {
	fmt.Println("name: ", r.Name)                              /* name of robot */
	fmt.Println("api_private: ", r.BinanceApiPrivate[:10]+"*") /* binance api_secret_key */
	fmt.Println("api_public: ", r.BinanceApiPublic[:10]+"*")   /* binance api_public_key */
	fmt.Println("stable_currency: ", r.StableCurrency)         /* stable_currency : USDT or BUSD */
	fmt.Println("tokens: ", r.Tokens)                          /* tokens : ETH,BTC,... */
	fmt.Println("strategy: ", r.Strategy)                      /* strategy : divergence, rsi, envelope */
	fmt.Println("long: ", r.CheckLong)                         /* open long position */
	fmt.Println("short: ", r.CheckShort)
	fmt.Println("futures: ", r.IsFutures)
	fmt.Println("start_balance: ", r.StartBalance)
	fmt.Println("stop_loss: ", r.StopLoss)
	fmt.Println("take_profit: ", r.TakeProfit)
}

func (this *Robot) OpenTrade() error {
	return nil
}

func (this *Robot) CloseTrade() error {
	return nil
}
