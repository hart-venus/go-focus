package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"golang.org/x/term"
)

type Config struct {
	ShortBreakLength   int    `json:"shortBreakLength"`
	LongBreakLength    int    `json:"longBreakLength"`
	PomodoroLength     int    `json:"pomodoroLength"`
	BreakInterval      int    `json:"breakInterval"`
	PauseAfterBreak    bool   `json:"pauseAfterBreak"`
	PauseAfterPomodoro bool   `json:"pauseAfterPomodoro"`
	ShortBreakMessage  string `json:"shortBreakMessage"`
	LongBreakMessage   string `json:"longBreakMessage"`
	PomodoroMessage    string `json:"pomodoroMessage"`
}

var defaultConfig = Config{
	ShortBreakLength:   (5 * 60),
	LongBreakLength:    (15 * 60),
	PomodoroLength:     (25 * 60),
	BreakInterval:      4,
	PauseAfterBreak:    false,
	PauseAfterPomodoro: true,
	ShortBreakMessage:  "Time for a short break!",
	LongBreakMessage:   "Time for a long break!",
	PomodoroMessage:    "Back to work!",
}

func loadConfig(path string) (Config, error) {
	config := defaultConfig // set defaults
	file, err := os.Open(path)
	if err != nil {
		// if file doesn't exist, just return the defaults
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}

	// Unmarshal the JSON encoding of Config into config.
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

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
	barWidth := width - 7 // 5 from time and 1 from space
	progressWidth := int(float64(barWidth) * percent)
	progressSide := strings.Repeat("\033[42m \033[0m", progressWidth)
	if percent != 1 {
		progressSide += "\033[32m\033[0m"
	}
	return fmt.Sprintf("%s%s", progressSide, strings.Repeat(" ", barWidth-progressWidth))
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

func createTimer(seconds int) {
	fullSeconds := seconds
	ticker := time.NewTicker(time.Second * 1) // creating a 1 second ticker
	done := make(chan bool)                   // channel to receive signal when to stop
	go func() {
		displayTime(seconds, float64(seconds)/float64(fullSeconds))
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				seconds--
				displayTime(seconds, float64(seconds)/float64(fullSeconds))
				if seconds == 0 {
					done <- true
				}
			}
		}
	}()
	<-done // block here until we receive the done signal
	ticker.Stop()
}

func displayTime(seconds int, percent float64) {
	removeLastLine()
	fmt.Print(formatTime(seconds))
	fmt.Print(" ")
	fmt.Println(getProgressBar(percent))
}

func pause() {
	fmt.Println("Press enter to continue...")
	fmt.Scanln()
}

func main() {
	// grab path from args
	path := "config.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	config, err := loadConfig(path)
	if err != nil {
		panic(err)
	}
	setUpNotification()
	fmt.Println()
	shortBreakLength := config.ShortBreakLength
	longBreakLength := config.LongBreakLength
	pomodoroLength := config.PomodoroLength
	breaks := 0
	breakinterval := config.BreakInterval

	pauseAfterBreak := config.PauseAfterBreak
	pauseAfterPomodoro := config.PauseAfterPomodoro

	shortBreakMessage := config.ShortBreakMessage
	longBreakMessage := config.LongBreakMessage
	pomodoroMessage := config.PomodoroMessage

	createTimer(pomodoroLength)
	for {
		if breaks == (breakinterval - 1) {
			sendNotification(longBreakMessage)
			if pauseAfterPomodoro {
				pause()
			}
			breaks = 0
			createTimer(longBreakLength)
		} else {
			sendNotification(shortBreakMessage)
			if pauseAfterPomodoro {
				pause()
			}
			breaks++
			createTimer(shortBreakLength)
		}
		sendNotification(pomodoroMessage)
		if pauseAfterBreak {
			pause()
		}
		createTimer(pomodoroLength)
	}
}
