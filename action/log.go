package action

import (
	"fmt"
	"os"
)

func Info(msg ...any) {
	fmt.Println(msg...)
}

func Debug(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
}
