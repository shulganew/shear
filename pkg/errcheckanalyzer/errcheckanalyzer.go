package errcheckanalyzer

import (
	"golang.org/x/tools/go/analysis"
)

// Create osExit analizer. Check if os.Exit call exist in the main package for the main function.
func NewAnalyzer() (osExitCheckAnalyzer *analysis.Analyzer) {
	// Custom analizer.
	osExitCheckAnalyzer = &analysis.Analyzer{
		Name: "osExit",
		Doc:  "Check if os.Exit call exist in the main package for the main function.",
		Run:  run,
	}
	return
}
