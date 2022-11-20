// Prints a mermaid graph to stdout from a go module. Inspect the result at
// https://mermaid.live
// Example output:
//
// flowchart BT
//         A --> B
//         C --> B
//         D --> B
//         E --> B
//         F --> B
//         G --> B
//         B --> H
//         C --> A
//         G --> D
//         A[get-doer/config]
//         B[get-doer/api]
//         C[get-doer/db]
//         D[get-doer/slack]
//         E[get-doer/slackController]
//         F[get-doer/zd]
//         G[get-doer/internal]
//         H[get-doer/MAIN]

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

const mode packages.LoadMode = packages.NeedName |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule

var ignore = []string{}

func main() {
	pattern := flag.String("pattern", "./...", "Go package pattern")
	ignoreString := flag.String("ignore", "", "List of packages to ignore, comma separated")
	ignorePrefix := flag.String("ignorePrefix", "", "Prefix to ignore from each package name")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Expecting a single argument: directory of module")
	}

	ignore = strings.Split(*ignoreString, ",")

	var fset = token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: flag.Args()[0]}
	pkgs, err := packages.Load(cfg, *pattern)
	if err != nil {
		log.Fatal(err)
	}

	imports := make(map[string]map[string]bool)

	for _, pkg := range pkgs {
		if ignorePgk(pkg.PkgPath) {
			continue
		}

		pkgPath := strings.ReplaceAll(pkg.PkgPath, *ignorePrefix, "")

		if pkg.Name == "main" {
			pkgPath += "/MAIN"
		}

		imports[pkgPath] = make(map[string]bool)

		for _, file := range pkg.Syntax {

			ast.Inspect(file, func(n ast.Node) bool {
				is, ok := n.(*ast.ImportSpec)
				if !ok {
					return true
				}

				importPath := strings.ReplaceAll(is.Path.Value, "\"", "")

				if ignorePgk(importPath) {
					return true
				}

				importPath = strings.ReplaceAll(importPath, *ignorePrefix, "")

				if strings.Contains(is.Path.Value, pkg.Module.Path) {
					imports[pkgPath][importPath] = true
				}

				return true
			})
		}
	}

	letters := NewLetters()

	fmt.Println("flowchart BT")
	for importer, imports := range imports {
		for importee := range imports {
			fmt.Printf(
				"\t%s --> %s\n",
				letters.Key(importee),
				letters.Key(importer),
			)
		}
	}

	for val, key := range letters.Keys {
		fmt.Printf("\t%s[%v]\n", key, val)
	}
}

func NewLetters() *Letters {
	return &Letters{
		current: 64,
		Keys:    make(map[string]string),
	}
}

type Letters struct {
	current int
	Keys    map[string]string
}

func (l *Letters) next() string {
	l.current++

	return string(l.current)
}

func (l *Letters) Key(val string) string {
	key, ok := l.Keys[val]
	if ok {
		return key
	}

	key = l.next()
	l.Keys[val] = key

	return key
}

func ignorePgk(pkg string) bool {
	for _, i := range ignore {
		if i != "" && strings.Contains(pkg, i) {
			return true
		}
	}

	return false
}
