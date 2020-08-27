package main

import (
	"log"
	"time"
)

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println("Tick at ", t)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true
	log.Println("Ticker executed successfullyl")
}
