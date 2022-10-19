//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/grafana/pkg/codegen"
	"github.com/grafana/grafana/pkg/cuectx"
	"github.com/grafana/grafana/pkg/framework/kind"
)

// Core kinds code generator. Produces all generated code in grafana/grafana
// that derives from raw and structured core kinds.

// All the single-kind generators to be run for core kinds.
var singles = []codegen.KindGenerator{}

// All the aggregate generators to be run for core kinds.
var multis = []codegen.AggregateKindGenerator{}

const sep = string(filepath.Separator)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "plugin thema code generator does not currently accept any arguments\n, got %q", os.Args)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get working directory: %s", err)
		os.Exit(1)
	}
	grootp := strings.Split(cwd, sep)
	groot := filepath.Join(sep, filepath.Join(grootp[:len(grootp)-1]...))

	wd := codegen.NewWriteDiffer()
	rt := cuectx.GrafanaThemaRuntime()
	var all []*codegen.DeclForGen

	// structured kinds first
	f := os.DirFS(filepath.Join(groot, "kind", "structured"))
	ents := elsedie(fs.ReadDir(f, "."))("error reading structured fs root directory")
	for _, ent := range ents {
		rel := filepath.Join("kind", "structured", ent.Name())
		sub := elsedie(fs.Sub(f, ent.Name()))(fmt.Sprintf("error creating subfs for path %s", rel))
		decl, err := kind.LoadCoreKindFS[kind.CoreStructuredMeta](sub, rel, rt.Context())
		if err != nil {
			die(fmt.Errorf("core structured kind at %s is invalid: %w", rel, err))
		}
		if decl.Meta.Name != ent.Name() {
			die(fmt.Errorf("%s: kind name (%s) must equal parent dir name (%s)", rel, decl.Meta.Name, ent.Name()))
		}

		all = append(all, elsedie(codegen.ForGen(rt, decl.Some()))(rel))
	}

	// now raw kinds
	f = os.DirFS(filepath.Join(groot, "kind", "raw"))
	ents = elsedie(fs.ReadDir(f, "."))("error reading raw fs root directory")
	for _, ent := range ents {
		rel := filepath.Join("kind", "raw", ent.Name())
		sub := elsedie(fs.Sub(f, ent.Name()))(fmt.Sprintf("error creating subfs for path %s", rel))
		pk, err := kind.LoadCoreKindFS[kind.RawMeta](sub, rel, rt.Context())
		if err != nil {
			die(fmt.Errorf("raw kind at %s is invalid: %w", rel, err))
		}
		if pk.Meta.Name != ent.Name() {
			die(fmt.Errorf("%s: kind name (%s) must equal parent dir name (%s)", rel, pk.Meta.Name, ent.Name()))
		}
		dfg, _ := codegen.ForGen(nil, pk.Some())
		all = append(all, dfg)
	}

	// Sort em real good
	sort.Slice(all, func(i, j int) bool {
		return nameFor(all[i].Meta) < nameFor(all[j].Meta)
	})

	// Run all single generators
	for _, gen := range singles {
		for _, decl := range all {
			gf, err := gen.Generate(decl)
			if err != nil {
				die(fmt.Errorf("%s: %w", err))
			}
			wd[filepath.Join(groot, gf.RelativePath)] = gf.Data
		}
	}

	// Run all multi generators
	for _, gen := range multis {
		gf, err := gen.Generate(all)
		if err != nil {
			die(fmt.Errorf("%s: %w", err))
		}
		wd[filepath.Join(groot, gf.RelativePath)] = gf.Data
	}

	if _, set := os.LookupEnv("CODEGEN_VERIFY"); set {
		err = wd.Verify()
		if err != nil {
			die(fmt.Errorf("generated code is out of sync with inputs:\n%s\nrun `make gen-cue` to regenerate\n\n", err))
		}
	} else {
		err = wd.Write()
		if err != nil {
			die(fmt.Errorf("error while writing generated code to disk:\n%s\n", err))
		}
	}
}

func nameFor(m kind.SomeKindMeta) string {
	switch x := m.(type) {
	case kind.RawMeta:
		return x.Name
	case kind.CoreStructuredMeta:
		return x.Name
	case kind.CustomStructuredMeta:
		return x.Name
	case kind.SlotImplMeta:
		return x.Name
	default:
		// unreachable so long as all the possibilities in KindMetas have switch branches
		panic("unreachable")
	}
}

func elsedie[T any](t T, err error) func(msg string) T {
	if err != nil {
		return func(msg string) T {
			fmt.Fprintf(os.Stderr, "%s: %s", msg, err)
			os.Exit(1)
			return t
		}
	}
	return func(msg string) T {
		return t
	}
}

func die(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}
