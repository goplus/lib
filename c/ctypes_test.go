//go:build !llgo

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c_test

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

// Check all Wasm profiles on the host so ordinary Go CI covers source selection
// and integer widths without relying on LLGo's unsafe.Sizeof constant folding.
func TestWasmCTypes(t *testing.T) {
	for _, test := range []struct {
		name, goos, tag  string
		signed, unsigned types.BasicKind
	}{
		{"GoJS", "js", "", types.Int32, types.Uint32},
		{"Emscripten32", "js", "llgo.wasm.emscripten", types.Int32, types.Uint32},
		{"Emscripten64", "js", "llgo.wasm.emscripten.memory64", types.Int64, types.Uint64},
		{"WASI", "wasip1", "llgo.wasm.wasi", types.Int32, types.Uint32},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := build.Default
			ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled = test.goos, "wasm", false
			ctx.BuildTags = []string{"llgo", test.tag}
			selected, err := ctx.ImportDir(".", 0)
			if err != nil {
				t.Fatal(err)
			}
			pkg := checkCTypes(t, selected, "wasm")
			for name, want := range map[string]types.BasicKind{
				"Long": test.signed, "Ulong": test.unsigned,
				"SizeT": test.unsigned, "SsizeT": test.signed,
				"PtrdiffT": test.signed, "IntptrT": test.signed, "UintptrT": test.unsigned,
				"WcharT": types.Int32, "Char": types.Int8,
			} {
				obj := pkg.Scope().Lookup(name)
				if obj == nil {
					t.Fatalf("missing %s in %v", name, selected.GoFiles)
				}
				if got := obj.Type().Underlying(); got != types.Typ[want] {
					t.Errorf("%s has underlying type %s, want %s", name, got, types.Typ[want])
				}
			}
		})
	}
}

func checkCTypes(t *testing.T, selected *build.Package, arch string) *types.Package {
	t.Helper()
	fset := token.NewFileSet()
	var files []*ast.File
	for _, name := range selected.GoFiles {
		file, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, file)
	}
	conf := types.Config{Importer: importer.Default(), IgnoreFuncBodies: true, Sizes: types.SizesFor("gc", arch)}
	pkg, err := conf.Check("github.com/goplus/lib/c", fset, files, nil)
	if err != nil {
		t.Fatal(err)
	}
	return pkg
}

func TestNativeCScalarTypes(t *testing.T) {
	for _, test := range []struct {
		goos, arch  string
		char, wchar types.BasicKind
	}{
		{"linux", "amd64", types.Int8, types.Int32},
		{"linux", "386", types.Int8, types.Int32},
		{"linux", "arm", types.Uint8, types.Uint32},
		{"linux", "arm64", types.Uint8, types.Uint32},
		{"linux", "riscv64", types.Uint8, types.Int32},
		{"linux", "loong64", types.Int8, types.Int32},
		{"linux", "ppc64", types.Uint8, types.Int32},
		{"linux", "ppc64le", types.Uint8, types.Int32},
		{"linux", "s390x", types.Uint8, types.Int32},
		{"linux", "mips", types.Int8, types.Int32},
		{"linux", "mips64", types.Int8, types.Int32},
		{"freebsd", "arm", types.Uint8, types.Uint32},
		{"freebsd", "arm64", types.Uint8, types.Uint32},
		{"freebsd", "riscv64", types.Uint8, types.Int32},
		{"netbsd", "arm", types.Uint8, types.Int32},
		{"netbsd", "arm64", types.Uint8, types.Int32},
		{"openbsd", "arm", types.Uint8, types.Int32},
		{"openbsd", "arm64", types.Uint8, types.Int32},
		{"darwin", "amd64", types.Int8, types.Int32},
		{"darwin", "arm64", types.Int8, types.Int32},
		{"windows", "386", types.Int8, types.Uint16},
		{"windows", "amd64", types.Int8, types.Uint16},
		{"windows", "arm64", types.Int8, types.Uint16},
	} {
		t.Run(test.goos+"/"+test.arch, func(t *testing.T) {
			ctx := build.Default
			ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled = test.goos, test.arch, false
			ctx.BuildTags = nil
			selected, err := ctx.ImportDir(".", 0)
			if err != nil {
				t.Fatal(err)
			}
			pkg := checkCTypes(t, selected, test.arch)
			for name, want := range map[string]types.BasicKind{
				"Char": test.char, "WcharT": test.wchar,
				"PtrdiffT": types.Int, "IntptrT": types.Int, "UintptrT": types.Uintptr,
			} {
				if got := pkg.Scope().Lookup(name).Type().Underlying(); got != types.Typ[want] {
					t.Errorf("%s has underlying type %s, want %s", name, got, types.Typ[want])
				}
			}
		})
	}
}
