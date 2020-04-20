package input

import (
	"bufio"
	"fmt"
	"os"
)

var (
	Print  = make(chan bool)
	Pause  = make(chan bool)
	Resume = make(chan bool)
	Stop   = make(chan bool)
)

func init() {
	go read()
}

func read() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		switch stdin.Text() {
		case "PRINT":
			Print <- true
		case "PAUSE":
			Pause <- true
		case "RESUME":
			Resume <- true
		case "STOP":
			Stop <- true
		default:
			fmt.Fprintln(os.Stderr, "invalid operation: " + stdin.Text())
		}
	}
	if err := stdin.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin: " + err.Error())
	}
}
