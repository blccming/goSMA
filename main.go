package main

import (
	"github.com/blccming/goSMA/internal/api"
)

func main() {
	go api.StartUpdating()
	api.StartAPI()
}
