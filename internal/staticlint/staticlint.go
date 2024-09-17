package staticlint

import (
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
)

func NewAnalyzers() []*analysis.Analyzer {
	analyzers := make([]*analysis.Analyzer, 0)

	analysysChecks := listAnalysis()
	analyzers = append(analyzers, analysysChecks...)

	staticcheckChecks := listStaticcheck()
	analyzers = append(analyzers, staticcheckChecks...)

	analyzers = append(analyzers, errcheck.Analyzer)

	return analyzers
}
