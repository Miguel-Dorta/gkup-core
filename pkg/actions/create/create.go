package create

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/settings"
	"os"
)

func Run(path string, hashAlgorithm string) {
	output.Setup(1)
	output.NewGlobalStep("Create repository", 0)
	if err := repo.Create(path, &settings.Settings{HashAlgorithm: hashAlgorithm}); err != nil {
		output.PrintErrorf("error creating repository: %s", err)
		os.Exit(1)
	}
}
