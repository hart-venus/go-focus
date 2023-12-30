package main

import (
	"fmt"
	"time"
)

func main() {
	// creating a 1 second ticker
	ticker := time.NewTicker(time.Second * 1)
	done := make(chan bool) // channel to receive signal when to stop
	i := 0
	// every second, print "tick"
	go func() { // create a goroutine to run this function
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Println("tick", i)
				i++
			}
		}
	}()

	// wait for 5 seconds
	time.Sleep(time.Second * 5)
	ticker.Stop()
	done <- true // set channel to true to stop the goroutine
	fmt.Println("Ticker stopped")
}
