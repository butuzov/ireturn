package ireturn

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const nolintPrefix = "//nolint" // used for dissallow comments

const name string = "ireturn"

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     name,
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

		// 02. Does it return any results ?
		if f.Type == nil || f.Type.Results == nil {
			return
		}

		// 02. Is it allowed to be checked?
		if hasDisallowDirective(f.Doc) {
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

func hasDisallowDirective(cg *ast.CommentGroup) bool {
	if cg == nil {
		return false
	}

	return disallowDirectiveFound(cg)
}

func disallowDirectiveFound(cg *ast.CommentGroup) bool {
	for i := len(cg.List) - 1; i >= 0; i-- {
		comment := cg.List[i]
		if !strings.HasPrefix(comment.Text, nolintPrefix) {
			continue
		}

		startingIdx := len(nolintPrefix)
		for {
			idx := strings.Index(comment.Text[startingIdx:], name)
			if idx == -1 {
				break
			}

			if len(comment.Text[startingIdx+idx:]) == len(name) {
				return true
			}

			c := comment.Text[startingIdx+idx+len(name)]
			if c == '.' || c == ',' || c == ' ' || c == '	' {
				return true
			}
			startingIdx += idx + 1
		}

	}

	return false
}
