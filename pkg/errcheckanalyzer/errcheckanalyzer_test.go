package errcheckanalyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitCheckAnalyzer(t *testing.T) {
	// Custom analizer.
	var osExitCheckAnalyzer = &analysis.Analyzer{
		Name: "osExit",
		Doc:  "Check if os.Exit() call exist in the main package for the main function.",
		Run:  run,
	}
	analysistest.Run(t, analysistest.TestData(), osExitCheckAnalyzer, "./...")
}
