package staticlint

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

func ListStaticcheck() []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}
	return analyzers
}
