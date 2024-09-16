package staticlint

import "golang.org/x/tools/go/analysis"

func ListAnalyzers() []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0)

	analysysChecks := ListAnalysis()
	checks = append(checks, analysysChecks...)

	staticcheckChecks := ListStaticcheck()
	checks = append(checks, staticcheckChecks...)

	return checks
}
