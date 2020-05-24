package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// progress represents a progress (current/total)
type progress struct {
	Current uint64 `json:"current"`
	Total uint64 `json:"total"`
}

// globalProgress represents the progress of the whole operation as a whole
type globalProgress struct {
	progress
	Name string `json:"name"`
}

// partialProgress represents the progress of a step of the operation
type partialProgress struct {
	progress
	Details string `json:"details"`
}

var (
	Verbose bool
	status = struct {
		Global globalProgress `json:"global"`
		Partial partialProgress `json:"partial"`
		m sync.Mutex
	}{}

	ErrProcessStopped = errors.New("process stopped")
)

// Setup sets the global number of steps
func Setup(totalGlobalSteps uint64) {
	status.m.Lock()
	status.Global.Total = totalGlobalSteps
	status.m.Unlock()
}

// Sets a new global step. It resets all partial steps.
func NewGlobalStep(name string, totalPartialSteps uint64) {
	status.m.Lock()

	status.Global.Name = name
	status.Global.Current++
	status.Partial.Total = totalPartialSteps
	status.Partial.Current = 0
	status.Partial.Details = ""
	printNoLock()

	status.m.Unlock()
}

// Sets a new partial step.
func NewPartialStep(details string) {
	status.m.Lock()
	status.Partial.Current++
	status.Partial.Details = details
	if Verbose {
		printNoLock()
	}
	status.m.Unlock()
}

// Print prints the current status
func Print() {
	status.m.Lock()
	printNoLock()
	status.m.Unlock()
}

// printNoLock prints the current status without locking
func printNoLock() {
	data, _ := json.Marshal(&status)
	fmt.Println(string(data))
}

// PrintErrorf prints an error formating it like fmt.Printf
func PrintErrorf(format string, a ...interface{}) {
	PrintError(fmt.Sprintf(format, a...))
}

// PrintError prints an error line
func PrintError(err interface{}) {
	fmt.Fprintln(os.Stderr, err)
}
