package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/vektah/gqlparser/v2/ast"
)

var version string = "dev"

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %s [option] schema1.graphql schema2.graphql ...\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(out, "Examples:\n")
		fmt.Fprintf(out, "  %s schema1.graphql schema2.graphql\n", os.Args[0])
		fmt.Fprintf(out, "  %s -skip '^(Animal|FooScalar)$' schema1.graphql\n", os.Args[0])
	}
	skip := flag.String("skip", "", "skip detection as an unused function")
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		fmt.Println("Need arguments.")
		os.Exit(1)
	}
	var skipRegexp *regexp.Regexp
	if *skip != "" {
		r, err := regexp.Compile(*skip)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		skipRegexp = r
	}

	app := &App{}
	if err := app.Load(flag.Args()...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if unused := app.DetectUnused(skipRegexp); len(unused) > 0 {
		for _, t := range unused {
			fmt.Printf("unused %s %s at %s line %d\n",
				kindName(t.Kind), t.Name, t.SourceFile, t.SourceLine)
		}
		os.Exit(1)
	}
}

func kindName(kind ast.DefinitionKind) string {
	switch kind {
	case ast.Scalar:
		return "scalar"
	case ast.Object:
		return "type"
	case ast.Interface:
		return "interface"
	case ast.Union:
		return "union"
	case ast.Enum:
		return "enum"
	case ast.InputObject:
		return "input"
	}
	panic("unpexected")
}
