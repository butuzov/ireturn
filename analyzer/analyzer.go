package analyzer

import (
	"fmt"
	"go/ast"
	gotypes "go/types"
	"strings"

	"github.com/butuzov/ireturn/config"
	"github.com/butuzov/ireturn/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const name string = "ireturn" // linter name

type validator interface {
	IsValid(types.IFace) bool
}

func NewAnalyzerWithConfig(validate validator) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     name,
		Doc:      "Accept Interfaces, Return Concrete Types",
		Run:      run(validate),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func NewAnalyzer() *analysis.Analyzer {
	return NewAnalyzerWithConfig(config.DefaultValidatorConfig())
}

func run(r validator) func(*analysis.Pass) (interface{}, error) {
	if r == nil {
		r = config.DefaultValidatorConfig()
	}

	return func(pass *analysis.Pass) (interface{}, error) {
		var issues []analysis.Diagnostic

		ins, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		ins.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
			// 001. Casting to funcdecl
			f, _ := node.(*ast.FuncDecl)

			// 002. Does it return any results ?
			if f.Type == nil || f.Type.Results == nil {
				return
			}

			// 003. Is it allowed to be checked?
			// TODO(butuzov): add inline comment
			if hasDisallowDirective(f.Doc) {
				return
			}

			// 004. Filtering Results.
			for _, i := range filterInterfaces(pass, f.Type.Results) {

				if r.IsValid(i) {
					continue
				}

				issues = append(issues, analysis.Diagnostic{ //nolint: exhaustivestruct
					Pos:     f.Pos(),
					Message: fmt.Sprintf("%s returns interface (%s)", f.Name.Name, i.Name),
				})
			}
		})

		for i := range issues {
			pass.Report(issues[i])
		}

		return nil, nil
	}
}

func filterInterfaces(pass *analysis.Pass, fl *ast.FieldList) []types.IFace {
	var results []types.IFace

	for pos, el := range fl.List {
		switch v := el.Type.(type) {
		// ----- empty or anonymous interfaces
		case *ast.InterfaceType:

			if len(v.Methods.List) == 0 {
				results = append(results, issue("interface{}", pos, types.EmptyInterface))
				continue
			}

			results = append(results, issue("anonymouse interface", pos, types.AnonInterface))

		case *ast.Ident:

			t1 := pass.TypesInfo.TypeOf(el.Type)
			if !gotypes.IsInterface(t1.Underlying()) {
				continue
			}

			word := t1.String()
			// only build in interface is error
			if obj := gotypes.Universe.Lookup(word); obj != nil {
				results = append(results, issue(obj.Name(), pos, types.ErrorInterface))
				continue
			}

			results = append(results, issue(word, pos, types.NamedInterface))

		case *ast.SelectorExpr:

			t1 := pass.TypesInfo.TypeOf(el.Type)
			if !gotypes.IsInterface(t1.Underlying()) {
				continue
			}

			word := t1.String()
			if isStdLib(word) {
				results = append(results, issue(word, pos, types.NamedStdInterface))
				continue
			}

			results = append(results, issue(word, pos, types.NamedInterface))
		}
	}

	return results
}

// isStdLib will run small checks against pkg to find out if  named interface
// we lookling on comes from a standard library or not.
func isStdLib(named string) bool {
	// find last dot index.
	idx := strings.LastIndex(named, ".")
	if idx == -1 {
		return false
	}

	if _, ok := std[named[0:idx]]; ok {
		return true
	}

	return false
}

// issue is shortcut that creates issue for next filtering
func issue(name string, pos int, interfaceType types.IType) types.IFace {
	return types.IFace{
		Name: name,
		Pos:  pos,
		Type: interfaceType,
	}
}
