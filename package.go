//
// parser.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package enumer

import (
	"go/ast"
	"go/constant"
	"go/importer"
	"go/token"
	"go/types"
	"log"
)

type Package struct {
	name     string
	files    []PackageFile
	defs     map[*ast.Ident]types.Object
	typesPkg *types.Package
}

type PackageFile struct {
	name      string
	ast       *ast.File
	constants []Constant
	pkg       *Package
}

// Constant represents a string like constant in go file
// i.e constant defined like that:
//
//     const Value ConstantType = "value of this constant"
//
type Constant struct {
	Name  string
	Type  string
	Value string
}

type OutputFile struct {
	Package string
	Cmd     string
}

func (f *PackageFile) constantsOfType(typeName string) func(ast.Node) bool {
	return func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			// We only care about const declarations.
			return true
		}
		typ := ""
		for _, spec := range decl.Specs {
			vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
			if vspec.Type == nil && len(vspec.Values) > 0 {
				// "X = 1". With no type but a value, the constant is untyped.
				// Skip this vspec and reset the remembered type.
				typ = ""
				continue
			}
			if vspec.Type != nil {
				// "X T". We have a type. Remember it.
				ident, ok := vspec.Type.(*ast.Ident)
				if !ok {
					continue
				}
				typ = ident.Name
			}
			if typ != typeName {
				// This is not the type we're looking for.
				continue
			}
			// We now have a list of names (from one line of source code) all being
			// declared with the desired type.
			// Grab their names and actual values and store them in f.values.
			for _, name := range vspec.Names {
				if name.Name == "_" {
					continue
				}
				// This dance lets the type checker find the values for us. It's a
				// bit tricky: look up the object declared by the name, find its
				// types.Const, and extract its value.
				obj, ok := f.pkg.defs[name]
				if !ok {
					log.Fatalf("no value for constant %s", name)
				}
				info := obj.Type().Underlying().(*types.Basic).Info()
				if info&types.IsString == 0 {
					log.Fatalf("can't handle non-string constant type %s", typ)
				}
				value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
				if value.Kind() != constant.String {
					log.Fatalf("can't happen: constant is not a string %s", name)
				}
				c := Constant{
					Name:  name.Name,
					Type:  typeName,
					Value: value.String(),
				}
				f.constants = append(f.constants, c)
			}
		}
		return false
	}
}

// check type-checks the package. The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{Defs: pkg.defs}
	astFiles := make([]*ast.File, len(pkg.files))
	for i, f := range pkg.files {
		astFiles[i] = f.ast
	}
	typesPkg, err := config.Check(".", fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}
