package action

import (
	"fmt"
)

func Info(msg ...any) {
	fmt.Printf("\033[0;34m%s\033[0m", fmt.Sprintln(msg...))
}

func Debug(msg ...any) {
	fmt.Printf("\033[0;36m%s\033[0m", fmt.Sprintln(msg...))
}
