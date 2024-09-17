package staticlint

import (
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/sudeeya/metrics-harvester/internal/exitcheck"
)

func NewAnalyzers() []*analysis.Analyzer {
	analyzers := make([]*analysis.Analyzer, 0)

	analysysChecks := listAnalysis()
	analyzers = append(analyzers, analysysChecks...)

	staticcheckChecks := listStaticcheck()
	analyzers = append(analyzers, staticcheckChecks...)

	analyzers = append(analyzers, errcheck.Analyzer)

	analyzers = append(analyzers, exitcheck.Analyzer)

	return analyzers
}
