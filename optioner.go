// Copyright 2014 AKUALAB INC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
optioner is a tool to generate functional options. Intended to be
used with go generate; see the invocation in example/example.go.

To learn about functional options, see:
http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
http://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html

This code was adapted from the stringer cmd (golang.org/x/tools/cmd/stringer/).

optioner will generate a file with code of the form:

  // N sets a value for instances of type Example.
  func N(o int) optExample {
	  return func(t *Example) optExample {
		  previous := t.N
		  t.N = o
		  return N(previous)
	  }
  }

The file is created in the package where the code is generated.

optioner will create options for all the fields in the struct except those that include
the tag `opt:"-"`.

For example, given this snippet,

  package example

  //go:generate optioner -type Example
  type Example struct {
	  N      int
	  FSlice []float64 `json:"float_slice"`
	  Map    map[string]int
	  Name   string `opt:"-" json:"name"`
	  ff     func(int) int
  }

  func NewExample(name string, options ...optExample) *Example {

	  // Set required values and initialize optional fields with default values.
	  ex := &Example{
		  Name:   name,
		  N:      10,
		  FSlice: make([]float64, 0, 100),
		  Map:    make(map[string]int),
		  ff:     func(n int) int { return n },
	  }

	  // Set options.
	  ex.Option(options...)
  }

go generate will generate option functions for fields N, FSlice, Map, and ff. Your package
users can now set options as follows:

  myFunc := func(n int) int {return 2 * n}
  ex := example.NewExample("test", example.N(22), example.Ff(myFunc))

the new struct "ex" will use default values for "FSlice" and "Map", and custom values for
"N" and "ff". Note that the argument "name" in NewExample() is required. For this reason
the struct field "name" is excluded using a tag.

To temporarily modify a value, do the following:

  prev := ex.Option(N(5)) // previous value is stored in prev.
  // do something...
  ex.Option(prev) // restores the previous value.

struct fields don't need to be exported, however, the corresponding option will be exported by
capitalizing the first letter. Documentation for options is auto-generated in the source file
to make it available to package users in godoc format.

It is possible to create options for various types in the same package by using various annotations

    //go:generate optioner -type Type1

    //go:generate optioner -type Type2

However, to keep function names short, there is no namespaces so you must use different names for different
option functions in the same package.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/go/types"
)

var (
	typeNameArg  = flag.String("type", "", "type name of the options struct; must be set")
	typeNameMngl = flag.String("m", "", "type name of the option; defaults to opt<type>")
	output       = flag.String("output", "", "output file name; default srcdir/<type>_gen_opt.go")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\toptioner [flags] -type T\n")
	fmt.Fprintf(os.Stderr, "Use struct tag ```opt:\"-\"``` to exclude fields\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttp://github.com/akualab/optioner\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("optioner: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNameArg) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Parse the package once.
	var (
		g Generator
	)

	if *typeNameMngl == "" {
		g.optName = "opt" + *typeNameArg
	} else {
		g.optName = *typeNameMngl
	}
	g.typeName = *typeNameArg
	if g.optName == g.typeName {
		log.Fatal("option and type names must be different")
	}
	g.options = []option{}

	g.parsePackage()

	// Print the header and package clause.
	g.Printf("// generated by optioner %s; DO NOT EDIT\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf(header)
	g.Printf("package %s", g.packageName)
	g.Printf("\n")
	g.Printf("// %s type is used to set options in %s.\n", g.optName, g.typeName)
	//	g.Printf("type option func(*%s) option\n", g.typeName)
	g.Printf("type %s func(*%s) %[1]s\n", g.optName, g.typeName)
	g.Printf("\n")
	g.Printf("// Option method sets the options. Returns previous option for last arg.\n")
	//	g.Printf("func (t *%s) Option(options ...option) (previous option) {\n", g.typeName)
	g.Printf("func (t *%[1]s) Option(options ...%s) (previous %[2]s) {\n", g.typeName, g.optName)
	g.Printf("for _, opt := range options {\n")
	g.Printf("previous = opt(t)\n")
	g.Printf("}\n")
	g.Printf("return previous\n")
	g.Printf("}\n")
	g.Printf("\n")

	for _, opt := range g.options {
		tname := strings.Title(opt.optName)
		g.Printf("// %s sets a value for instances of type %s.\n", tname, g.typeName)
		g.Printf("func %s(o %s) %s {\n", tname, opt.typ, g.optName)
		g.Printf("return func(t *%s) %s {\n", g.typeName, g.optName)
		g.Printf("previous := t.%s\n", opt.name)
		g.Printf("t.%s = o\n", opt.name)
		g.Printf("return %s(previous)\n", tname)
		g.Printf("}\n")
		g.Printf("}\n")
		g.Printf("\n")
	}

	// Format the output.
	src := g.format()

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_gen_opt.go", g.typeName)
		outputName = strings.ToLower(baseName)
	}
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Generator) parsePackage() {

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatalf("cannot process directory: %s", err)
	}

	for _, f := range pkg.GoFiles {
		if found := g.parseFields(f); found {
			break
		}
	}
}

func (g *Generator) parseFields(fn string) bool {

	log.Println("parse file: ", fn)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fn, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	g.packageName = f.Name.Name
	for _, v := range f.Decls {
		if genDecl, ok := v.(*ast.GenDecl); ok {
			if genDecl.Tok != token.TYPE {
				continue
			}
			for _, s := range genDecl.Specs {
				if typeSpec, ok := s.(*ast.TypeSpec); ok {
					if typeSpec.Name.String() != g.typeName {
						continue
					}
					if structDecl, ok := typeSpec.Type.(*ast.StructType); ok {
						fields := structDecl.Fields.List
						for _, field := range fields {
							id := field.Names[0]

							var optName string
							// check struct tags to exclude fields
							if field.Tag != nil {
								s := strings.Replace(field.Tag.Value, "`", "", -1)
								tag := reflect.StructTag(s).Get("opt")
								switch tag {
								case "-":
									continue
								default:
									optName = tag
								}
							}
							if optName == "" {
								optName = id.Name
							}
							typeExpr := field.Type
							typ := types.ExprString(typeExpr)
							g.options = append(g.options, option{name: id.Name, typ: typ, optName: optName})
							log.Printf("generating option %q for field %q of type %q)", optName, id.Name, typ)
						}
						return true
					} else {
						log.Fatal("target type is not a struct")
					}
				}
			}
		}
	}
	return false
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf         bytes.Buffer // Accumulated output.
	options     []option
	typeName    string
	optName     string
	packageName string
}

type option struct {
	name    string
	typ     string
	optName string
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

const header = `
// Please report issues and submit contributions at:
// http://github.com/akualab/optioner
// optioner is a project of AKUALAB INC.

`
