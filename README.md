# go-focus
A simple pomodoro timer for the command line.
## Requirements
- Go
- [Powerline Font](https://github.com/powerline/fonts)
## Usage
```bash
go run main.go <config file>
```
## Configuration
If the configuration file is not provided, the program will look for a file named `config.json` in the same directory as the executable. If no configuration file is found, the program will use the default configuration. If any of the values in the configuration file are invalid, the program will use the default value for that field.

The configuration file is a JSON file with the following fields:
- `shortBreakLength` - The length of the short break in seconds
- `longBreakLength` - The length of the long break in seconds
- `pomodoroLength` - The length of the pomodoro in seconds
- `breakInterval` - The number of pomodoros before a long break
- `pauseAfterPomodoro` - Whether or not to pause after a pomodoro
- `pauseAfterBreak` - Whether or not to pause after a break
- `shortBreakMessage` - The message to display when a short break starts
- `longBreakMessage` - The message to display when a long break starts
- `pomodoroMessage` - The message to display when a pomodoro starts
