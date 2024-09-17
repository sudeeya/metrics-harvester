package main

import (
	"github.com/sudeeya/metrics-harvester/internal/staticlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := staticlint.NewAnalyzers()
	multichecker.Main(
		analyzers...,
	)
}
