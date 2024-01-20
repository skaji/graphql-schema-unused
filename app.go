package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type Type struct {
	Name       string
	Kind       ast.DefinitionKind
	BuiltIn    bool
	SourceFile string
	SourceLine int
	Union      []string
	Interfaces []string
	Fields     []*Field
}

func (t *Type) KindName() string {
	switch t.Kind {
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

type Field struct {
	Name      string
	Type      string
	Arguments []*Argument
}

type Argument struct {
	Name string
	Type string
}

type Types []*Type

func (ts Types) Get(name string) *Type {
	for _, t := range ts {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type App struct {
	types Types
}

func (a *App) Load(paths ...string) error {
	sources := make([]*ast.Source, len(paths))
	for i, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		sources[i] = &ast.Source{
			Name:    path,
			Input:   string(b),
			BuiltIn: false,
		}
	}
	schema, err := gqlparser.LoadSchema(sources...)
	if err != nil {
		return err
	}

	ds := make([]*ast.Definition, 0, len(schema.Types))
	for _, d := range schema.Types {
		ds = append(ds, d)
	}
	types := make(Types, len(ds))
	for i, d := range ds {
		types[i] = &Type{
			Name:       d.Name,
			Kind:       d.Kind,
			Union:      d.Types,
			Interfaces: d.Interfaces,
			BuiltIn:    d.BuiltIn,
			SourceLine: d.Position.Line,
			SourceFile: d.Position.Src.Name,
		}
	}
	for i, d := range ds {
		fs := make([]*Field, len(d.Fields))
		for j, f := range d.Fields {
			args := make([]*Argument, len(f.Arguments))
			for k, a := range f.Arguments {
				args[k] = &Argument{
					Name: a.Name,
					Type: a.Type.Name(),
				}
			}
			slices.SortFunc(args, func(a, b *Argument) int {
				return strings.Compare(a.Name, b.Name)
			})
			fs[j] = &Field{
				Name:      f.Name,
				Type:      f.Type.Name(),
				Arguments: args,
			}
		}
		slices.SortFunc(fs, func(a, b *Field) int {
			return strings.Compare(a.Name, b.Name)
		})
		types[i].Fields = fs
	}
	slices.SortFunc(types, func(a, b *Type) int {
		return strings.Compare(a.Name, b.Name)
	})
	a.types = types
	return nil
}

func (a *App) DetectUnused() []*Type {
	seen := map[string]bool{}
	var walk func(t *Type)
	walk = func(t *Type) {
		seen[t.Name] = true
		for _, u := range t.Union {
			if !seen[u] {
				walk(a.types.Get(u))
			}
		}
		for _, i := range t.Interfaces {
			seen[i] = true // XXX
		}
		for _, f := range t.Fields {
			if strings.HasPrefix(f.Type, "__") {
				continue
			}
			if !seen[f.Type] {
				walk(a.types.Get(f.Type))
			}
			for _, arg := range f.Arguments {
				if !seen[arg.Type] {
					walk(a.types.Get(arg.Type))
				}
			}
		}
	}

	for _, name := range []string{"Query", "Mutation", "Subscription"} {
		t := a.types.Get(name)
		if t == nil {
			continue
		}
		walk(t)
	}
	var unused Types
	for _, t := range a.types {
		if !seen[t.Name] && !t.BuiltIn {
			unused = append(unused, t)
		}
	}
	return unused
}
