package main

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
)

func main() {
	output.Verbose = false
	repo.Init("/")
}
