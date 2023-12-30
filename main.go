package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"golang.org/x/term"
)

func formatTime(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

var terminalFd int = int(os.Stdout.Fd())
var notf *notificator.Notificator

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
	percent = 1 - percent // make it go from left to right
	width := terminalWidth()
	barWidth := width - 5 // 5 from time and 1 from space
	progressWidth := int(float64(barWidth) * percent)
	return fmt.Sprintf("%s%s", strings.Repeat("=", progressWidth), strings.Repeat(" ", barWidth-progressWidth))
}

func setUpNotification() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	notf = notificator.New(notificator.Options{
		AppName:     "go-focus",
		DefaultIcon: dir + "/icon/default.png",
	})
}

func sendNotification(message string) {
	notf.Push("Go Focus", message, "", notificator.UR_NORMAL)
}

func main() {
	setUpNotification()
	// creating a 1 second ticker
	ticker := time.NewTicker(time.Second * 1)
	twentyfiveminutes := 3

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
