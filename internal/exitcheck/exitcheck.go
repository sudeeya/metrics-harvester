// Package exitcheck provides an analyzer that checks for absence of os.Exit calls in main function.
package exitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for absence of os.Exit calls in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if funcDecl, ok := node.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "main" {
					for _, stmt := range funcDecl.Body.List {
						if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
							if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
								if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
									if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
										pass.Reportf(callExpr.Pos(), "os.Exit call in main")
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
