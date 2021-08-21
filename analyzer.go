package ireturn

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "ireturn",
		Doc:      "Accept Interfaces, Return Concrete Types",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, wrapChecker(pass))

	return nil, nil
}

func wrapChecker(pass *analysis.Pass) func(ast.Node) {
	return func(node ast.Node) {
		//nolint: forcetypeassert // (that's ok, we filtering out non func in the inspect.Preorder).
		f := node.(*ast.FuncDecl)

		if f.Type == nil || f.Type.Results == nil {
			return
		}

		for _, i := range filterInterfaces(pass, f.Type.Results) {
			pass.Reportf(f.Pos(), "%s returns interface (%s)", f.Name.Name, i.name)
		}
	}
}

func filterInterfaces(pass *analysis.Pass, fl *ast.FieldList) []iface {
	var results []iface

	for pos, el := range fl.List {
		switch v := el.Type.(type) {
		// -----
		case *ast.InterfaceType:

			if len(v.Methods.List) == 0 {
				results = append(results, iface{
					name: "interface{}",
					pos:  pos,
					t:    typeEmptyInterface,
				})
				continue
			}

			results = append(results, iface{
				name: "anonymouse interface",
				pos:  pos,
				t:    typeAnonInterface,
			})

		// -----
		default:
			_ = v

		}
	}

	return results
}
