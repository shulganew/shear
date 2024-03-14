package main

import (
	"github.com/shulganew/shear.git/pkg/staticlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(staticlint.NewAnalyzer()...)
}
