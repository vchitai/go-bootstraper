package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	filePath := "assets/plasma/internal/services/server_name_lower_gen.go"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		panic(err)
	}

	err = ast.Print(fset, f)
	if err != nil {
		panic(err)
	}
}