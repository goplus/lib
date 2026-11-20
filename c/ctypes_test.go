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
		{"WASI ignores Emscripten Memory64 tag", "wasip1", "llgo.wasm.emscripten.memory64", types.Int32, types.Uint32},
	} {
		t.Run(test.name, func(t *testing.T) {
			checkCTypes(t, test.goos, "wasm", []string{"llgo", test.tag}, map[string]types.BasicKind{
				"Long": test.signed, "Ulong": test.unsigned,
				"SizeT": test.unsigned, "SsizeT": test.signed,
				"PtrdiffT": test.signed, "IntptrT": test.signed, "UintptrT": test.unsigned,
				"WcharT": types.Int32, "Char": types.Int8,
			})
		})
	}
}

func checkCTypes(t *testing.T, goos, arch string, tags []string, wantTypes map[string]types.BasicKind) {
	t.Helper()
	ctx := build.Default
	ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled = goos, arch, false
	ctx.BuildTags = tags
	selected, err := ctx.ImportDir(".", 0)
	if err != nil {
		t.Fatal(err)
	}
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
	for name, want := range wantTypes {
		obj := pkg.Scope().Lookup(name)
		if obj == nil {
			t.Fatalf("missing %s in %v", name, selected.GoFiles)
		}
		if got := obj.Type().Underlying(); got != types.Typ[want] {
			t.Errorf("%s has underlying type %s, want %s", name, got, types.Typ[want])
		}
	}
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
		{"linux", "mipsle", types.Int8, types.Int32},
		{"linux", "mips64", types.Int8, types.Int32},
		{"linux", "mips64le", types.Int8, types.Int32},
		{"freebsd", "386", types.Int8, types.Int32},
		{"freebsd", "amd64", types.Int8, types.Int32},
		{"freebsd", "arm", types.Uint8, types.Uint32},
		{"freebsd", "arm64", types.Uint8, types.Uint32},
		{"freebsd", "riscv64", types.Uint8, types.Int32},
		{"netbsd", "386", types.Int8, types.Int32},
		{"netbsd", "amd64", types.Int8, types.Int32},
		{"netbsd", "arm", types.Uint8, types.Int32},
		{"netbsd", "arm64", types.Uint8, types.Int32},
		{"openbsd", "386", types.Int8, types.Int32},
		{"openbsd", "amd64", types.Int8, types.Int32},
		{"openbsd", "arm", types.Uint8, types.Int32},
		{"openbsd", "arm64", types.Uint8, types.Int32},
		{"openbsd", "ppc64", types.Uint8, types.Int32},
		{"openbsd", "riscv64", types.Uint8, types.Int32},
		{"darwin", "amd64", types.Int8, types.Int32},
		{"darwin", "arm64", types.Int8, types.Int32},
		{"ios", "arm64", types.Int8, types.Int32},
		{"android", "386", types.Int8, types.Int32},
		{"android", "amd64", types.Int8, types.Int32},
		{"android", "arm", types.Uint8, types.Uint32},
		{"android", "arm64", types.Uint8, types.Uint32},
		{"solaris", "amd64", types.Int8, types.Int32},
		{"illumos", "amd64", types.Int8, types.Int32},
		{"windows", "386", types.Int8, types.Uint16},
		{"windows", "amd64", types.Int8, types.Uint16},
		{"windows", "arm64", types.Int8, types.Uint16},
	} {
		t.Run(test.goos+"/"+test.arch, func(t *testing.T) {
			// Native Unix targets use ILP32 or LP64; Windows uses LLP64.
			long, ulong := types.Int64, types.Uint64
			if test.goos == "windows" || types.SizesFor("gc", test.arch).Sizeof(types.Typ[types.Uintptr]) == 4 {
				long, ulong = types.Int32, types.Uint32
			}
			size, ssize := ulong, long
			if test.goos == "windows" {
				size, ssize = types.Uintptr, types.Int
			}
			checkCTypes(t, test.goos, test.arch, nil, map[string]types.BasicKind{
				"Char": test.char, "WcharT": test.wchar,
				"PtrdiffT": types.Int, "IntptrT": types.Int, "UintptrT": types.Uintptr,
				"Long": long, "Ulong": ulong, "SizeT": size, "SsizeT": ssize,
			})
		})
	}
}
