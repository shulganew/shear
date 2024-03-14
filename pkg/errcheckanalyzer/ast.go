package errcheckanalyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Check AST systax tree for os.Exit call in the main func.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
		// Finding main package.
		ast.Inspect(file, func(node ast.Node) bool {
			if f, ok := node.(*ast.File); ok {
				if f.Name.Name == "main" {
					// Finding main method.
					ast.Inspect(node, func(pmnode ast.Node) bool {
						if fd, ok := pmnode.(*ast.FuncDecl); ok {
							if fd.Name.Name == "main" {
								// Inspect main method.
								ast.Inspect(pmnode, func(mnode ast.Node) bool {
									// if os.Exit in go statement
									if _, ok := mnode.(*ast.GoStmt); ok {
										return false
									}
									if sl, ok := mnode.(*ast.SelectorExpr); ok {
										if ind, ok := sl.X.(*ast.Ident); ok {
											if ind.Name == "os" && sl.Sel.Name == "Exit" {
												pass.Reportf(sl.Pos(), "strait call os.Exit in the main method prohibited")
											}
										}
									}
									return true
								})
								return false
							}
						}
						return true
					})
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}
