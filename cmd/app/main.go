package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-binance-robot/internal/binance"
	"github.com/go-binance-robot/internal/robot"
)

var _ = fmt.Sprintf("%d", 0)

func RunInfinitely(r *robot.Robot) {
	tokenFinished := make(map[string]bool)
	mu := sync.Mutex{}
	for _, token := range r.Tokens {
		tokenFinished[token] = true
	}

	for {
		for i := range r.Tokens {
			mu.Lock()
			finished := tokenFinished[r.Tokens[i]]
			mu.Unlock()
			time.Sleep(2 * time.Second)

			if finished {
				go func(ts *robot.Trade, tokenFin *map[string]bool) {
					mu.Lock()
					(*tokenFin)[ts.Token] = false
					mu.Unlock()

					binance.WebSocketRun(ts, 500)
					log.Printf("%s socket finished\n", ts.Token)

					mu.Lock()
					(*tokenFin)[ts.Token] = true
					mu.Unlock()
				}(&r.TradingSession[i], &tokenFinished)
			}
		}
	}
}

func main() {
	r := robot.New()
	RunInfinitely(r)
}
