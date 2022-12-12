package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-binance-robot/internal/binance"
	"github.com/go-binance-robot/internal/robot"
)

var _ = fmt.Sprintf("%d", 0)


func main() {
	r := robot.New()
	robot.Run(r)
}
