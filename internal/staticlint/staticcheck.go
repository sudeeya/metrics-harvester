package staticlint

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func ListStaticcheck() []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		if v.Analyzer.Name == "S1000" {
			analyzers = append(analyzers, v.Analyzer)
			break
		}
	}

	for _, v := range stylecheck.Analyzers {
		if v.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, v.Analyzer)
			break
		}
	}

	for _, v := range quickfix.Analyzers {
		if v.Analyzer.Name == "QF1001" {
			analyzers = append(analyzers, v.Analyzer)
			break
		}
	}

	return analyzers
}
