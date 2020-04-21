package input

import (
	"bufio"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"os"
)

var (
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
			output.Print()
		case "PAUSE":
			Pause <- true
		case "RESUME":
			Resume <- true
		case "STOP":
			Stop <- true
		default:
			output.PrintErrorf("invalid operation: %s", stdin.Text())
		}
	}
	if err := stdin.Err(); err != nil {
		output.PrintErrorf("error reading stdin: %s", err)
		os.Exit(1)
	}
}
