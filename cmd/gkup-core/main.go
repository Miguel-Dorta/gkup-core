package main

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/input"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
)

func main() {
	<- input.Print
	output.Print()
}
