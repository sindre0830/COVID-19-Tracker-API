package notification

import (
	"fmt"
	"time"
)

var Ticker *time.Ticker

func Schedule() {
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case t := <- Ticker.C:
			fmt.Println("Tick at", t)
		}
	}
}
