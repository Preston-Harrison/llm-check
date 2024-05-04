package go_checks

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Project struct {
	Pkgs []*packages.Package
	Fset *token.FileSet
}

func (p *Project) Files() []*ast.File {
	var allFiles []*ast.File
	for _, pkg := range p.Pkgs {
		allFiles = append(allFiles, pkg.Syntax...) // Append all AST files from each package
	}
	return allFiles
}

func NewProject(modulePath string) (*Project, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Dir:   modulePath,
		Fset:  token.NewFileSet(),
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	return &Project{pkgs, cfg.Fset}, nil
}

// NewProjectFromSource creates a new project from a source string.
func NewProjectFromSource(source string) (*Project, error) {
	fset := token.NewFileSet()

	// Parse the source string to create an AST
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	files := []*ast.File{file}

	// Create a new type checker and type information container
	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	// Perform type-checking on the parsed file
	pkg, err := conf.Check("main", fset, files, info)
	if err != nil {
		return nil, err
	}

	// Package the type-checked package into the format used by our project
	loadedPkg := &packages.Package{
		Fset:      fset,
		ID:        "main",
		PkgPath:   "main",
		Syntax:    files,
		Types:     pkg,
		TypesInfo: info,
	}

	return &Project{Pkgs: []*packages.Package{loadedPkg}, Fset: fset}, nil
}

func (p *Project) IsSameFunction(funcDecl *types.Func, callExpr *ast.CallExpr) (bool, error) {
	// Find the called function in the call expression
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		// If it's not a simple identifier, it could be a selector (e.g., pkg.Func)
		sel, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			// Stuff like {}byte("hello world")
			return false, nil
		}
		ident = sel.Sel
	}

	// Use type information to determine the called function's object
	var calledFuncObj types.Object
	for _, pkg := range p.Pkgs {
		if obj, exists := pkg.TypesInfo.Uses[ident]; exists {
			calledFuncObj = obj
			break
		}
	}

	if calledFuncObj == nil {
		return false, fmt.Errorf("function called in expression not found in type info")
	}

	// Check if the function from the CallExpr is the same as funcDecl
	return calledFuncObj == funcDecl, nil
}

// GetFuncObject retrieves the *types.Func corresponding to the given *ast.FuncDecl.
func (p *Project) GetFuncObject(funcDecl *ast.FuncDecl) (*types.Func, error) {
	if funcDecl.Name == nil {
		return nil, fmt.Errorf("function declaration has no name")
	}

	// Loop through each package's type info
	for _, pkg := range p.Pkgs {
		if obj, ok := pkg.TypesInfo.Defs[funcDecl.Name]; ok {
			if tFunc, ok := obj.(*types.Func); ok {
				return tFunc, nil
			} else {
				return nil, fmt.Errorf("identified object is not a function")
			}
		}
	}

	return nil, fmt.Errorf("function not found in any package")
}
