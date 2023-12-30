package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

func formatTime(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

var terminalFd int = int(os.Stdout.Fd())

func terminalWidth() int {
	width, _, err := term.GetSize(terminalFd)
	if err != nil {
		panic(err)
	}
	return width
}

func removeLastLine() {
	fmt.Print("\033[1A") // Move cursor up
	fmt.Print("\033[2K") // Delete line
	// finally, print space to move cursor to the beginning of the line
	fmt.Print("\r")
}

func getProgressBar(percent float64) string {
	width := terminalWidth()
	barWidth := width - 7 // 5 from time and 1 from space
	progressWidth := int(float64(barWidth) * percent)
	return fmt.Sprintf("[%s%s]", strings.Repeat("=", progressWidth), strings.Repeat(" ", barWidth-progressWidth))
}

func main() {
	// creating a 1 second ticker
	ticker := time.NewTicker(time.Second * 1)
	twentyfiveminutes := 60

	done := make(chan bool) // channel to receive signal when to stop
	// every second, print "tick"
	go func() { // create a goroutine to run this function
		// startup
		fmt.Print(formatTime(twentyfiveminutes))
		fmt.Println(getProgressBar(float64(twentyfiveminutes) / (60)))
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				removeLastLine()
				twentyfiveminutes--
				fmt.Print(formatTime(twentyfiveminutes))
				fmt.Println(getProgressBar(float64(twentyfiveminutes) / (60)))

				if twentyfiveminutes == 0 {
					done <- true
				}
			}
		}
	}()

	<-done // block here until we receive the done signal
}
