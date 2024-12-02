package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/butuzov/ireturn/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
