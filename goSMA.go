package main

import (
	"fmt"

	"github.com/blccming/goSMA/internal/metrics"
)

func main() {
	fmt.Println(metrics.CPU())
	fmt.Println(metrics.System())
}
