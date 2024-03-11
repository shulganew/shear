package main

import (
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func docgen(analyzers ...*analysis.Analyzer) {
	var alist strings.Builder
	intro := `
	My Check (mychech) is a custom generator, witch consist of:
	- passes analysers (golang.org/x/tools/go/analysis/passes)
	- staticcheck analyser (all SA* packages and "QF1006", "QF1010", "QF1007")
	- custom osExit checker


	How to use:
	mycheck [packages]
	Example:
	./cmd/staticlint/mycheck ./cmd/... ./internal/...
	
	Detail description of all linters in set:
	`
	alist.WriteString(intro)

	for _, a := range analyzers {
		alist.WriteString(">>>>        ")
		alist.WriteString(a.Name)
		alist.WriteString("        <<<<\n\n")
		alist.WriteString(a.Doc)
		alist.WriteString("\n\n")
	}

	// add comments slashes to all lines
	doc := strings.Replace(alist.String(), "\n", "\n // ", -1)


	file, err := os.OpenFile("./doc.go", os.O_CREATE|os.O_RDWR, 0775)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	file.Write([]byte(doc))
	file.Write([]byte("\n"))
	file.Write([]byte("package main"))
	file.Close()
}
