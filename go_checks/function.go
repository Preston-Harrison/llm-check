package go_checks

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"llm-check/comments"
	"llm-check/common"
	"log"
)

var _ common.Function = (*GoFunction)(nil)

type GoFunction struct {
	pkg  string
	fn   *ast.FuncDecl
	proj *Project
}

type FunctionWithCheck struct {
	GoFunction
	check common.Check
}

func (fwc FunctionWithCheck) Check() common.Check {
	return fwc.check
}

func NewGoFunction(fn *ast.FuncDecl, proj *Project, pkg string) GoFunction {
	return GoFunction{pkg, fn, proj}
}

func (gf GoFunction) Name() string {
	return gf.fn.Name.Name
}

func (gf GoFunction) Module() string {
	return gf.pkg
}

func (gf GoFunction) Location() (token.Position, token.Position) {
	return gf.proj.Fset.Position(gf.fn.Pos()), gf.proj.Fset.Position(gf.fn.End())
}

func (gf GoFunction) Docs() []string {
	var docs []string
	if gf.fn.Doc == nil {
		return docs
	}
	for _, comment := range gf.fn.Doc.List {
		docs = append(docs, comment.Text)
	}

	return docs
}

func (gf GoFunction) String() string {
	var buf bytes.Buffer
	printer := &printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	if err := printer.Fprint(&buf, gf.proj.Fset, gf.fn); err != nil {
		log.Fatalf("failed to print function declaration: %v", err)
	}
	return buf.String()
}

func (gf GoFunction) CallSites() []common.Function {
	var fns []common.Function
	// var fileContext *ast.File

	funcFilter := func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			typefunc, err := gf.proj.GetFuncObject(gf.fn)
			if err != nil {
				panic(err)
			}
			result, err := gf.proj.IsSameFunction(typefunc, call)
			if err != nil {
				panic(err)
			}
			return result
		}

		return false
	}

	funcProcessor := func(n []ast.Node) {
		if len(n) < 1 {
			return
		}
		// Remove the called function
		n = n[:len(n)-2]

		// Find next function call higher in tree
		for i := len(n) - 1; i >= 0; i-- {
			if fn, ok := n[i].(*ast.FuncDecl); ok {
				pkg := n[0].(*ast.File)
				fns = append(fns, NewGoFunction(fn, gf.proj, pkg.Name.String()))
			}
		}
	}

	siteTracker := newVisitor(funcFilter, funcProcessor)

	for _, file := range gf.proj.Files() {
		ast.Walk(siteTracker, file)
	}

	return fns
}

func GetFnsWithCheck(proj *Project) []common.FunctionWithCheck {
	var fns []common.FunctionWithCheck

	for _, file := range proj.Files() {
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Doc != nil {
				newFn := NewGoFunction(fn, proj, file.Name.String())
				checks := comments.ParseChecksFromComments(newFn.Docs())
				if len(checks) != 0 {
					// TODO: look at using more than just the first check
					fns = append(fns, FunctionWithCheck{
						newFn,
						checks[0],
					})
				}
			}
			return true
		})
	}

	return fns
}
