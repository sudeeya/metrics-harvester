package staticlint

import (
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/sudeeya/metrics-harvester/internal/exitcheck"
)

// NewAnalyzers returns all used analyzers.
func NewAnalyzers() []*analysis.Analyzer {
	analyzers := make([]*analysis.Analyzer, 0)

	analysysChecks := listAnalysis()
	analyzers = append(analyzers, analysysChecks...) // add go/analysis analyzers

	staticcheckChecks := listStaticcheck()
	analyzers = append(analyzers, staticcheckChecks...) // add staticcheck analyzers

	analyzers = append(analyzers, errcheck.Analyzer) // add errcheck analyzer

	analyzers = append(analyzers, exitcheck.Analyzer) // add exitcheck analyzer

	return analyzers
}
