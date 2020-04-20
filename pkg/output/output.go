package output

import (
	"encoding/json"
	"fmt"
	"sync"
)

type progress struct {
	Current uint64 `json:"current"`
	Total uint64 `json:"total"`
}

type globalProgress struct {
	progress
	Name string `json:"name"`
}

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
)

func Setup(totalGlobalSteps uint64) {
	status.m.Lock()
	status.Global.Total = totalGlobalSteps
	status.m.Unlock()
}

func NewGlobalStep(name string, totalPartialSteps uint64) {
	status.m.Lock()

	status.Global.Name = name
	status.Global.Current++
	status.Partial.Total = totalPartialSteps
	status.Partial.Current = 0
	status.Partial.Details = ""
	if Verbose {
		printNoLock()
	}

	status.m.Unlock()
}

func NewPartialStep(details string) {
	status.m.Lock()
	status.Partial.Current++
	status.Partial.Details = details
	if Verbose {
		printNoLock()
	}
	status.m.Unlock()
}

func Print() {
	status.m.Lock()
	printNoLock()
	status.m.Unlock()
}

func printNoLock() {
	data, _ := json.Marshal(&status)
	fmt.Println(string(data))
}
