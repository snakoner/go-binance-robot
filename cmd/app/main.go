package main

import (
	"fmt"

	"github.com/go-binance-robot/internal/binance"
	"github.com/go-binance-robot/internal/robot"
)

var _ = fmt.Sprintf("%d", 0)

func main() {
	symbols := []string{"ADA", "AVAX", "ETH"}
	for _, s := range symbols {
		robot := robot.New()
		robot.Print()

		// binance.WebSocketGetClose(robot)
		go binance.WebSocketRun(robot, s, 500)
	}

	for {

	}
}
