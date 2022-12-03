package main

import (
	"fmt"

	"github.com/go-binance-robot/internal/binance"
	"github.com/go-binance-robot/internal/robot"
)

var _ = fmt.Sprintf("%d", 0)

func main() {
	robot := robot.New()
	robot.Print()

	// binance.WebSocketGetClose(robot)
	binance.WebSocketRun(robot, "ADA", 500)

}
