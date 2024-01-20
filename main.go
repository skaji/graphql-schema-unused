package main

import (
	"flag"
	"fmt"
	"os"
)

var version string = "dev"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s schema1.graphql schema2.graphql ...\n", os.Args[0])
		flag.PrintDefaults()
	}
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

	app := &App{}
	if err := app.Load(flag.Args()...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if unused := app.DetectUnused(); len(unused) > 0 {
		for _, t := range unused {
			fmt.Printf("unused %s %s at %s line %d\n",
				t.KindName(), t.Name, t.SourceFile, t.SourceLine)
		}
		os.Exit(1)
	}
}
