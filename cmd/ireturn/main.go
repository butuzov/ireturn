package main

import (
	"github.com/butuzov/ireturn"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(ireturn.NewAnalyzer())
}
