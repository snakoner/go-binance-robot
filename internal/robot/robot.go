package robot

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-binance-robot/internal/indicators"
	"github.com/go-binance-robot/internal/strategy"
	"github.com/go-binance-robot/internal/binance"
	"github.com/joho/godotenv"
)

type Robot struct {
	Name              string   /* name of robot */
	BinanceApiPrivate string   /* binance api_secret_key */
	BinanceApiPublic  string   /* binance api_public_key */
	StableCurrency    string   /* stable_currency : USDT or BUSD */
	Tokens            []string /* tokens : ETH,BTC,... */
	StrategyNames     []string /* strategy : divergence, rsi, envelope */
	CheckLong         bool     /* open long position */
	CheckShort        bool     /* open short position */
	Timeframe         string
	IsFutures         bool    /* futures of spot trading */
	StartBalance      float64 /* value to buy on */
	TakeProfit        float64 /* take profit value in percent */
	StopLoss          float64 /* stop loss value in percent */
	TradingSession    []Trade
	MaxTokensTrack    int
}

type Trade struct {
	Root                 *Robot
	Strategy             *strategy.Strategy
	Active               bool    /* deal is active or not */
	Token                string  /* token name */
	BuyValue             float64 /* value bought on */
	OpenPrice            float64 /* deal open price */
	Quantity             float64
	LastTime             time.Time /* time of last price */
	Close                []float64 /* actual candle close counts */
	StopLossValue        float64
	TakeProfitValue      float64
	LastPriceForSLChange float64
	Result               TradingResult
}

type TradingResult struct {
	ProfitPerc  float64
	ProfitPrice float64
	StartTime   time.Time
	EndTime     time.Time
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
	robot.StrategyNames = strings.Split(os.Getenv("strategy"), ",")
	robot.CheckLong, _ = strconv.ParseBool(os.Getenv("check_long"))
	robot.CheckShort, _ = strconv.ParseBool(os.Getenv("check_short"))
	robot.Timeframe = os.Getenv("timeframe")
	robot.IsFutures, _ = strconv.ParseBool(os.Getenv("futures"))
	robot.StartBalance, _ = strconv.ParseFloat(os.Getenv("start_balance"), 64)
	robot.TakeProfit, _ = strconv.ParseFloat(os.Getenv("take_profit"), 64)
	robot.StopLoss, _ = strconv.ParseFloat(os.Getenv("stop_loss"), 64)
	robot.MaxTokensTrack, _ = strconv.Atoi(os.Getenv("max_tokens_track"))
	// robot.Strategy = strategy.New()

	for i, symbol := range robot.Tokens {
		ts := Trade{
			Token: symbol,
		}
		robot.TradingSession = append(robot.TradingSession, ts)
		robot.TradingSession[i].Root = robot
		robot.TradingSession[i].Strategy = new(strategy.Strategy)

		for _, name := range robot.StrategyNames {
			switch name {
			case "envelope":
				robot.TradingSession[i].Strategy.Add(indicators.Envelope, name)
				break
			case "divergence":
				robot.TradingSession[i].Strategy.Add(indicators.Divergence, name)
				break
			case "rsi":
				//robot.Strategy.Add(indicators.Rsi, name)
				break
			default:
				log.Fatalf("Unknown indicator: %s\n", name)
			}
		}
	}

	return robot
}

func (r *Robot) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(len(r.Tokens))
	for i := range r.Tokens {
		go func(wg *sync.WaitGroup, ts *robot.Trade) {
			binance.WebSocketRun(ts, 500)
			wg.Done()
			log.Printf("%s socket finished\n", ts.Token)
		}(wg, &r.TradingSession[i])
	}
	wg.Wait()
}

func (r *Robot) Print() {
	fmt.Println(strings.ToUpper("name: "), r.Name)                              /* name of robot */
	fmt.Println(strings.ToUpper("api_private: "), r.BinanceApiPrivate[:10]+"*") /* binance api_secret_key */
	fmt.Println(strings.ToUpper("api_public: "), r.BinanceApiPublic[:10]+"*")   /* binance api_public_key */
	fmt.Println(strings.ToUpper("stable_currency: "), r.StableCurrency)         /* stable_currency : USDT or BUSD */
	fmt.Println(strings.ToUpper("tokens: "), r.Tokens)                          /* tokens : ETH,BTC,... */
	fmt.Println(strings.ToUpper("strategy: "), r.StrategyNames)                 /* strategy : divergence, rsi, envelope */
	fmt.Println(strings.ToUpper("long: "), r.CheckLong)                         /* open long position */
	fmt.Println(strings.ToUpper("short: "), r.CheckShort)
	fmt.Println(strings.ToUpper("futures: "), r.IsFutures)
	fmt.Println(strings.ToUpper("start_balance: "), r.StartBalance)
	fmt.Println(strings.ToUpper("stop_loss: "), r.StopLoss)
	fmt.Println(strings.ToUpper("take_profit: "), r.TakeProfit)

	//fmt.Println(r.Strategy.GetName())

	fmt.Println()
}
