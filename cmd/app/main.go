package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-binance-robot/internal/trade"
)

var _ = fmt.Sprintf("%d", 0)

func RunInfinitely(r *trade.Robot) {
	mu := sync.Mutex{}
	wg := new(sync.WaitGroup)

	wg.Add(len(r.Tokens))
	for i := range r.Tokens {
		mu.Lock()
		trade.GetPrecision(&r.TradingSession[i])
		mu.Unlock()
		time.Sleep(2 * time.Second)

		go func(ts *trade.Trade) {
			trade.WebSocketRun(ts, 500)
			log.Printf("%s socket finished\n", ts.Token)
			wg.Done()
		}(&r.TradingSession[i])
	}
	wg.Wait()
}

func main() {
	r := trade.New()
	RunInfinitely(r)
}
