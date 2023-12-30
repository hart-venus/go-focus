package main

import (
	"fmt"
	"time"
)

func formatTime(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func main() {
	// creating a 1 second ticker
	ticker := time.NewTicker(time.Second * 1)
	twentyfiveminutes := 60 * 25

	done := make(chan bool) // channel to receive signal when to stop
	// every second, print "tick"
	go func() { // create a goroutine to run this function
		// startup
		fmt.Println(formatTime(twentyfiveminutes))
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				twentyfiveminutes--
				fmt.Print("\033[A\033[2K") // clear the previous line
				fmt.Println(formatTime(twentyfiveminutes))
				if twentyfiveminutes == 0 {
					done <- true
				}
			}
		}
	}()

	<-done // block here until we receive the done signal
}
