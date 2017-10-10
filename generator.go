//
// generator.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package enumer

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
)

type Generator struct {
	pkg *Package // Package we are scanning.
}

// files collects all .go file names in directory dir,
// *_test.go files are ignored
func files(dir string) []string {
	pkg, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		log.Fatalf("cannot process %q: %s", dir, err)
	}
	if dir == "." {
		return pkg.GoFiles
	}
	fl := make([]string, len(pkg.GoFiles))
	for i := range pkg.GoFiles {
		fl[i] = dir + "/" + pkg.GoFiles[i]
	}
	return fl
}

func (g *Generator) parse(dir string) {
	g.pkg = new(Package)
	if g.pkg.files == nil {
		g.pkg.files = []PackageFile{}
	}
	fs := token.NewFileSet()
	for _, name := range files(dir) {
		f, err := parser.ParseFile(fs, name, nil, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		g.pkg.files = append(
			g.pkg.files,
			PackageFile{
				pkg:       g.pkg,
				name:      name,
				ast:       f,
				constants: []Constant{},
			},
		)
	}
	if len(g.pkg.files) == 0 {
		log.Fatalf("no buildable Go files in %q", dir)
	}
	g.pkg.name = g.pkg.files[0].ast.Name.Name
	g.pkg.check(fs)
}

func (g *Generator) generate(typeName string) {
	constants := make([]Constant, 0, 100)
	for _, file := range g.pkg.files {
		file.constants = []Constant{}
		if file.ast != nil {
			ast.Inspect(file.ast, file.constantsOfType(typeName))
			constants = append(constants, file.constants...)
		}
	}
	if len(constants) == 0 {
		log.Fatalf("no constant values defined for type %s", typeName)
	}
	fmt.Fprintln(os.Stderr, constants)
}
