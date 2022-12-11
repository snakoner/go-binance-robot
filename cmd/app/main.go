package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-binance-robot/internal/binance"
	"github.com/go-binance-robot/internal/robot"
)

var _ = fmt.Sprintf("%d", 0)

func Run() {
	wg := &sync.WaitGroup{}
	r := robot.New()
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

func main() {
	// symbols := []string{"ADA", "AVAX", "ETH", "MATIC", "SOL", "FTM"}
	Run()
}
