package main

import (
	"flag"
	"fmt"
	"go/ast"
	"strings"

	"4d63.com/gochecknoglobals/checknoglobals"
	"github.com/butuzov/ireturn/analyzer"
	"github.com/kyoh86/exportloopref"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

//go:generate go run ./... -d
func main() {

	dgen := flag.Bool("d", false, "flag for doc genetation")
	flag.Parse()

	var analyzers []*analysis.Analyzer

	// Analysis package, passes group.
	passes := []*analysis.Analyzer{appends.Analyzer, asmdecl.Analyzer, httpmux.Analyzer, assign.Analyzer, atomic.Analyzer, atomicalign.Analyzer, bools.Analyzer, buildssa.Analyzer, buildtag.Analyzer, cgocall.Analyzer, composite.Analyzer, copylock.Analyzer, ctrlflow.Analyzer, deepequalerrors.Analyzer, defers.Analyzer, directive.Analyzer, errorsas.Analyzer, fieldalignment.Analyzer, findcall.Analyzer, framepointer.Analyzer, httpresponse.Analyzer, ifaceassert.Analyzer, inspect.Analyzer, loopclosure.Analyzer, lostcancel.Analyzer, nilfunc.Analyzer, nilness.Analyzer, pkgfact.Analyzer, printf.Analyzer, reflectvaluecompare.Analyzer, shadow.Analyzer, shift.Analyzer, sigchanyzer.Analyzer, slog.Analyzer, sortslice.Analyzer, stdmethods.Analyzer, stringintconv.Analyzer, structtag.Analyzer, testinggoroutine.Analyzer, tests.Analyzer, timeformat.Analyzer, unmarshal.Analyzer, unreachable.Analyzer, unsafeptr.Analyzer, unusedresult.Analyzer, unusedwrite.Analyzer, usesgenerics.Analyzer}

	// Staticcheck package, SA* group.
	var analyserSA []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyserSA = append(analyserSA, v.Analyzer)
		}
	}

	// Staticcheck package, QF* group.
	check := map[string]bool{
		"QF1006": true,
		"QF1010": true,
		"QF1007": true,
	}

	var quick []*analysis.Analyzer
	for _, v := range quickfix.Analyzers {
		if check[v.Analyzer.Name] {
			quick = append(quick, v.Analyzer)
		}
	}

	// Public Analizers.
	var pub []*analysis.Analyzer

	// Custom analizer.
	var osExitCheckAnalyzer = &analysis.Analyzer{
		Name: "osExit",
		Doc:  "Check if os.Exit call exist in the main package for the main function.",
		Run:  run,
	}

	pub = append(pub, exportloopref.Analyzer)
	pub = append(pub, checknoglobals.Analyzer())
	pub = append(pub, analyzer.NewAnalyzer())

	analyzers = append(analyzers, passes...)
	analyzers = append(analyzers, analyserSA...)
	analyzers = append(analyzers, quick...)
	analyzers = append(analyzers, pub...)
	analyzers = append(analyzers, osExitCheckAnalyzer)

	if !*dgen {
		// Create custom multichecker.
		fmt.Println("Analizers from go/analysis/passes: ", len(passes))
		fmt.Println("Staticcheck analyzers SA*: ", len(analyserSA))
		fmt.Println("Analizers quick checks: ", len(quick))
		fmt.Println("Analizers public checks: ", len(pub))
		multichecker.Main(analyzers...)
	} else {
		// Generate doc.go.
		fmt.Println("Generate user's documentation: ")
		docgen(analyzers...)
	}

}

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
												pass.Reportf(sl.Pos(), "strait call in the main method prohibited")
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
