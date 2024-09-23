package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/sudeeya/metrics-harvester/internal/staticlint"
)

func main() {
	analyzers := staticlint.NewAnalyzers()
	multichecker.Main(
		analyzers...,
	)
}
