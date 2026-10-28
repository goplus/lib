//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
)

type nativeABISlot struct {
	size, align uintptr
	class       int
}

func TestNativeLibcLayouts(t *testing.T) {
	t.Run("native function prototypes", func(t *testing.T) {
		checkSlot := func(t *testing.T, name string, got reflect.Type, want nativeABISlot) {
			if got.Size() != want.size || uintptr(got.Align()) != want.align {
				t.Errorf("%s size/alignment = %d/%d, C = %d/%d", name, got.Size(), got.Align(), want.size, want.align)
			}
			class := 0
			switch got.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				class = 1
			case reflect.Ptr, reflect.UnsafePointer, reflect.Func:
				class = 5
			case reflect.Float32, reflect.Float64:
				class = 8
			case reflect.Struct:
				class = 12
			}
			if class != want.class {
				t.Errorf("%s ABI class = %d (%v), C = %d", name, class, got.Kind(), want.class)
			}
		}
		for _, tc := range generatedFunctionPrototypes {
			t.Run(tc.name, func(t *testing.T) {
				wantInputs := len(tc.inputs)
				if tc.variadic {
					wantInputs++
				}
				if tc.signature.Kind() != reflect.Func || tc.signature.NumIn() != wantInputs || tc.signature.IsVariadic() != tc.variadic {
					t.Fatalf("Go parameters/variadic = %v, C fixed parameters = %d, variadic = %v", tc.signature, len(tc.inputs), tc.variadic)
				}
				for i, slot := range tc.inputs {
					checkSlot(t, "parameter", tc.signature.In(i), slot)
				}
				wantOutputs := 0
				if tc.output != nil {
					wantOutputs = 1
				}
				if tc.signature.NumOut() != wantOutputs {
					t.Fatalf("Go result count = %d, C = %d", tc.signature.NumOut(), wantOutputs)
				}
				if tc.output != nil {
					checkSlot(t, "result", tc.signature.Out(0), *tc.output)
				}
			})
		}
	})
	t.Run("symbol inventory", func(t *testing.T) {
		filename := "libc_" + runtime.GOOS + "_" + runtime.GOARCH + ".go"
		source, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		functions, types, constants, variables := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
		prototypes := map[string]bool{}
		for _, prototype := range generatedFunctionPrototypes {
			prototypes[prototype.name] = true
		}
		for _, symbol := range generatedFunctionSymbols {
			functions[symbol.name] = true
			if !prototypes[symbol.name] {
				t.Errorf("generated C function %s has no native prototype test", symbol.name)
			}
		}
		for _, name := range generatedAccessorSymbols {
			functions[name] = true
		}
		for _, tc := range generatedTypeLayouts {
			types[tc.symbol] = true
		}
		for _, tc := range generatedIntegerConstants {
			constants[tc.name] = true
		}
		for _, tc := range generatedStringConstants {
			constants[tc.name] = true
		}
		for _, symbol := range generatedVariableSymbols {
			variables[symbol.name] = true
		}
		for _, declaration := range source.Decls {
			switch decl := declaration.(type) {
			case *ast.FuncDecl:
				name := decl.Name.Name
				if !ast.IsExported(name) {
					continue
				}
				if decl.Recv != nil {
					receiver, prefix := decl.Recv.List[0].Type, ""
					if pointer, ok := receiver.(*ast.StarExpr); ok {
						receiver, prefix = pointer.X, "*"
					}
					id := receiver.(*ast.Ident).Name
					if !strings.HasPrefix(id, "_llcppg_anon_") {
						name = prefix + id + "." + name
					}
				}
				if !functions[name] {
					t.Errorf("generated function %s has no symbol test", name)
				}
			case *ast.GenDecl:
				for _, specification := range decl.Specs {
					switch spec := specification.(type) {
					case *ast.TypeSpec:
						if ast.IsExported(spec.Name.Name) && !types[spec.Name.Name] {
							t.Errorf("generated type %s has no native ABI test", spec.Name.Name)
						}
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							if !ast.IsExported(name.Name) {
								continue
							}
							if decl.Tok == token.CONST && name.Name == "LLGoPackage" {
								if c.LLGoPackage != "link" {
									t.Errorf("generated Go accessors require LLGoPackage=link, got %q", c.LLGoPackage)
								}
								continue
							}
							if decl.Tok == token.CONST && !constants[name.Name] {
								t.Errorf("generated constant %s has no native value test", name.Name)
							}
							if decl.Tok == token.VAR && !variables[name.Name] {
								t.Errorf("generated variable %s has no symbol test", name.Name)
							}
						}
					}
				}
			}
		}
	})
	t.Run("native constants", func(t *testing.T) {
		for _, tc := range generatedIntegerConstants {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go constant = %d, C headers = %d", tc.got, tc.want)
				}
			})
		}
		for _, tc := range generatedStringConstants {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go constant = %q, C headers = %q", tc.got, tc.want)
				}
			})
		}
	})
	t.Run("native types", func(t *testing.T) {
		for _, tc := range generatedTypeLayouts {
			t.Run(tc.symbol+"/"+tc.property, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go layout = %d, C headers = %d", tc.got, tc.want)
				}
			})
		}
	})
	t.Run("function symbols", func(t *testing.T) {
		for _, symbol := range generatedFunctionSymbols {
			t.Run(symbol.name, func(t *testing.T) {
				if symbol.address == nil {
					t.Fatal("C function address is nil")
				}
				t.Logf("C function address: %p", symbol.address)
			})
		}
	})
	t.Run("variable symbols", func(t *testing.T) {
		for _, symbol := range generatedVariableSymbols {
			t.Run(symbol.name, func(t *testing.T) {
				if symbol.address == nil {
					t.Fatal("C variable address is nil")
				}
				t.Logf("C variable address: %p", symbol.address)
			})
		}
	})
	var iovec c.Iovec
	var div c.DivT
	var ldiv c.LdivT
	var lldiv c.LldivT
	var imaxdiv c.ImaxdivT
	var filePosition nativeFilePosition
	// These values come from the original macOS 14 SDK and glibc 2.24 headers.
	cases := []struct {
		name                             string
		size, align, wantSize, wantAlign uintptr
	}{
		{"iovec", unsafe.Sizeof(iovec), unsafe.Alignof(iovec), 16, 8},
		{"div_t", unsafe.Sizeof(div), unsafe.Alignof(div), 8, 4},
		{"ldiv_t", unsafe.Sizeof(ldiv), unsafe.Alignof(ldiv), 16, 8},
		{"lldiv_t", unsafe.Sizeof(lldiv), unsafe.Alignof(lldiv), 16, 8},
		{"imaxdiv_t", unsafe.Sizeof(imaxdiv), unsafe.Alignof(imaxdiv), 16, 8},
		{"fpos_t", unsafe.Sizeof(filePosition), unsafe.Alignof(filePosition), nativeFilePositionSize, 8},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.size != tc.wantSize || tc.align != tc.wantAlign {
				t.Fatalf("size/alignment = %d/%d, want %d/%d", tc.size, tc.align, tc.wantSize, tc.wantAlign)
			}
		})
	}
	if unsafe.Offsetof(iovec.IovBase) != 0 || unsafe.Offsetof(iovec.IovLen) != 8 || unsafe.Sizeof(iovec.IovLen) != 8 {
		t.Fatalf("iovec fields: base=%d len=%d width=%d", unsafe.Offsetof(iovec.IovBase), unsafe.Offsetof(iovec.IovLen), unsafe.Sizeof(iovec.IovLen))
	}
	if unsafe.Offsetof(imaxdiv.Quot) != 0 || unsafe.Offsetof(imaxdiv.Rem) != 8 {
		t.Fatalf("imaxdiv_t field offsets = %d/%d, want 0/8", unsafe.Offsetof(imaxdiv.Quot), unsafe.Offsetof(imaxdiv.Rem))
	}
}

func TestNativeLibcIntegers(t *testing.T) {
	for _, tc := range generatedLifecycleCalls {
		t.Run(tc.name, tc.check)
	}
	t.Run("integer and floating conversions", func(t *testing.T) {
		signed := (*nativeChar)(unsafe.Pointer(c.Str("-4294967297tail")))
		unsigned := (*nativeChar)(unsafe.Pointer(c.Str("100000001tail")))
		var end *nativeChar
		for _, tc := range []struct {
			name      string
			got, want int64
		}{
			{"Abs", int64(c.Abs(-91)), 91},
			{"Labs", int64(c.Labs(-4294967297)), 4294967297},
			{"Llabs", int64(c.Llabs(-4294967297)), 4294967297},
			{"Imaxabs", int64(c.Imaxabs(-4294967297)), 4294967297},
			{"Atol", int64(c.Atol(signed)), -4294967297},
			{"Atoll", int64(c.Atoll(signed)), -4294967297},
			{"Strtol", int64(c.Strtol(signed, &end, 10)), -4294967297},
			{"Strtoll", int64(c.Strtoll(signed, &end, 10)), -4294967297},
			{"Strtoq", int64(c.Strtoq(signed, &end, 10)), -4294967297},
			{"Strtoul", int64(c.Strtoul(unsigned, &end, 16)), 4294967297},
			{"Strtoull", int64(c.Strtoull(unsigned, &end, 16)), 4294967297},
			{"Strtouq", int64(c.Strtouq(unsigned, &end, 16)), 4294967297},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("result = %d, want %d", tc.got, tc.want)
				}
			})
		}
		if end == nil || *end != 't' {
			t.Fatal("integer conversion end pointer differs")
		}
		floating := (*nativeChar)(unsafe.Pointer(c.Str("1.25tail")))
		if c.Atof(floating) != 1.25 || c.Strtod(floating, &end) != 1.25 || end == nil || *end != 't' {
			t.Fatal("double conversion differs")
		}
		if c.Strtof(floating, &end) != 1.25 || end == nil || *end != 't' {
			t.Fatal("float conversion differs")
		}
		wideSigned := [6]nativeWide{'-', '1', '2', '3', 'x', 0}
		wideUnsigned := [5]nativeWide{'1', '2', '3', 'x', 0}
		var wideEnd *nativeWide
		if c.Wcstoimax(&wideSigned[0], &wideEnd, 10) != -123 || wideEnd != &wideSigned[4] {
			t.Fatal("wcstoimax value/end differ")
		}
		if c.Wcstoumax(&wideUnsigned[0], &wideEnd, 10) != 123 || wideEnd != &wideUnsigned[3] {
			t.Fatal("wcstoumax value/end differ")
		}
	})
	t.Run("locale and multibyte calls", func(t *testing.T) {
		previous := c.Setlocale(c.LC_ALL, nil)
		if previous == nil {
			t.Fatal("setlocale query returned nil")
		}
		saved := c.GoString((*c.Char)(unsafe.Pointer(previous)))
		defer func() { c.Setlocale(c.LC_ALL, (*nativeChar)(unsafe.Pointer(c.AllocaCStr(saved)))) }()
		if c.Setlocale(c.LC_ALL, (*nativeChar)(unsafe.Pointer(c.Str("C")))) == nil {
			t.Fatal("C locale is unavailable")
		}
		locale := c.Localeconv()
		if locale == nil || c.GoString((*c.Char)(unsafe.Pointer(locale.DecimalPoint))) != "." {
			t.Fatal("localeconv decimal point differs")
		}
		input := [3]nativeChar{'a', 'b', 0}
		if c.Mblen(&input[0], 1) != 1 || c.Mblen(nil, 0) != 0 {
			t.Fatal("mblen result differs")
		}
		var character c.WcharT
		if c.Mbtowc(&character, &input[0], 1) != 1 || character != 'a' {
			t.Fatal("mbtowc result differs")
		}
		var bytes [8]nativeChar
		if c.Wctomb(&bytes[0], 'b') != 1 || bytes[0] != 'b' {
			t.Fatal("wctomb result differs")
		}
		var wide [4]c.WcharT
		if c.Mbstowcs(&wide[0], &input[0], 4) != 2 || wide != [4]c.WcharT{'a', 'b', 0, 0} {
			t.Fatalf("mbstowcs = %v", wide)
		}
		if c.Wcstombs(&bytes[0], &wide[0], 8) != 2 || bytes[0] != 'a' || bytes[1] != 'b' || bytes[2] != 0 {
			t.Fatal("wcstombs result differs")
		}
		if c.Strcoll(&input[0], (*nativeChar)(unsafe.Pointer(c.Str("ac")))) >= 0 {
			t.Fatal("strcoll ordering differs")
		}
		if c.Strxfrm(&bytes[0], &input[0], 8) != 2 || bytes[0] != 'a' || bytes[1] != 'b' || bytes[2] != 0 {
			t.Fatal("C locale strxfrm differs")
		}
	})
	t.Run("new memory and string calls", func(t *testing.T) {
		input := [8]nativeChar{'a', 'b', 'c', ':', 'd', 0}
		colon := (*nativeChar)(unsafe.Pointer(c.Str(":")))
		if c.Strcspn(&input[0], colon) != 3 || c.Strspn(&input[0], (*nativeChar)(unsafe.Pointer(c.Str("abc")))) != 3 || c.Strpbrk(&input[0], colon) != &input[3] {
			t.Fatal("string span results differ")
		}
		var output [16]nativeChar
		if c.Memccpy(unsafe.Pointer(&output[0]), unsafe.Pointer(&input[0]), ':', 6) != unsafe.Pointer(&output[4]) || output[3] != ':' || output[4] != 0 {
			t.Fatal("memccpy stop position differs")
		}
		c.Bzero(unsafe.Pointer(&output[0]), 16)
		c.Bcopy(unsafe.Pointer(&input[0]), unsafe.Pointer(&output[0]), 6)
		if c.Bcmp(unsafe.Pointer(&output[0]), unsafe.Pointer(&input[0]), 6) != 0 {
			t.Fatal("bcopy/bcmp differ")
		}
		word := (*nativeChar)(unsafe.Pointer(c.Str("abc")))
		if c.Stpcpy(&output[0], word) != &output[3] || output[0] != 'a' || output[3] != 0 {
			t.Fatal("stpcpy end pointer differs")
		}
		if c.Stpncpy(&output[0], word, 5) != &output[3] || output[3] != 0 || output[4] != 0 {
			t.Fatal("stpncpy padding differs")
		}
		if c.Strnlen(word, 2) != 2 || c.Strnlen(word, 8) != 3 {
			t.Fatal("strnlen bounds differ")
		}
		if c.Rindex(&input[0], ':') != &input[3] || c.Strcasecmp((*nativeChar)(unsafe.Pointer(c.Str("AbC"))), word) != 0 || c.Strncasecmp((*nativeChar)(unsafe.Pointer(c.Str("ABCx"))), word, 3) != 0 {
			t.Fatal("string search/case results differ")
		}
		if c.Ffs(0) != 0 || c.Ffs(8) != 4 {
			t.Fatal("ffs bit positions differ")
		}
		var errorText [128]nativeChar
		if c.StrerrorR(2, &errorText[0], 128) != 0 || errorText[0] == 0 {
			t.Fatal("strerror_r result differs")
		}
		if signalText := c.Strsignal(15); signalText == nil || *signalText == 0 {
			t.Fatal("strsignal returned no description")
		}
		path := (*nativeChar)(unsafe.Pointer(c.AllocaCStr("/tmp/leaf")))
		if c.GoString((*c.Char)(unsafe.Pointer(c.Dirname(path)))) != "/tmp" {
			t.Fatal("dirname result differs")
		}
		tokenInput := [5]nativeChar{'a', ':', 'b', 0}
		var save *nativeChar
		if c.StrtokR(&tokenInput[0], colon, &save) != &tokenInput[0] || tokenInput[1] != 0 || c.StrtokR(nil, colon, &save) != &tokenInput[2] || c.StrtokR(nil, colon, &save) != nil {
			t.Fatal("strtok_r state differs")
		}
		separated := [5]nativeChar{'a', ':', 'b', 0}
		cursor := &separated[0]
		if c.Strsep(&cursor, colon) != &separated[0] || cursor != &separated[2] || c.Strsep(&cursor, colon) != &separated[2] || cursor != nil {
			t.Fatal("strsep cursor differs")
		}
	})
	t.Run("allocation and callback calls", func(t *testing.T) {
		page := c.Valloc(64)
		if page == nil {
			t.Fatal("valloc returned nil")
		}
		defer c.Free(page)
		if uintptr(page)%uintptr(os.Getpagesize()) != 0 {
			t.Fatal("valloc is not page aligned")
		}
		aligned := c.AlignedAlloc(64, 128)
		if aligned == nil {
			t.Fatal("aligned_alloc returned nil")
		}
		defer c.Free(aligned)
		if uintptr(aligned)%64 != 0 {
			t.Fatal("aligned_alloc alignment differs")
		}
		var pointer unsafe.Pointer
		if result := c.PosixMemalign(&pointer, 64, 128); result != 0 || pointer == nil {
			t.Fatalf("posix_memalign = %d/%p", result, pointer)
		}
		defer c.Free(pointer)
		if uintptr(pointer)%64 != 0 {
			t.Fatal("posix_memalign alignment differs")
		}
		items := [4]c.Int{2, 5, 7, 11}
		key := c.Int(7)
		found := c.Bsearch(unsafe.Pointer(&key), unsafe.Pointer(&items[0]), 4, unsafe.Sizeof(key), func(a, b unsafe.Pointer) c.Int { return *(*c.Int)(a) - *(*c.Int)(b) })
		if found != unsafe.Pointer(&items[2]) {
			t.Fatalf("bsearch returned %p, want %p", found, &items[2])
		}
	})
	t.Run("structure returns", func(t *testing.T) {
		cases := []struct{ numerator, denominator, quotient, remainder c.Long }{
			{29, 6, 4, 5},
			{-29, 6, -4, -5},
			{29, -6, -4, 5},
			{-29, -6, 4, -5},
			{(1 << 40) + 7, 16, 1 << 36, 7},
		}
		for _, tc := range cases {
			got := c.Imaxdiv(tc.numerator, tc.denominator)
			if got.Quot != tc.quotient || got.Rem != tc.remainder {
				t.Fatalf("imaxdiv(%d,%d) = %d,%d, want %d,%d", tc.numerator, tc.denominator, got.Quot, got.Rem, tc.quotient, tc.remainder)
			}
		}
		div := c.Div(29, 6)
		ldiv := c.Ldiv(-29, 6)
		lldiv := c.Lldiv((1<<40)+7, 16)
		if div.Quot != 4 || div.Rem != 5 || ldiv.Quot != -4 || ldiv.Rem != -5 || lldiv.Quot != 1<<36 || lldiv.Rem != 7 {
			t.Fatalf("division structure returns: div=%v ldiv=%v lldiv=%v", div, ldiv, lldiv)
		}
	})
	t.Run("signed parsing", func(t *testing.T) {
		input := (*nativeChar)(unsafe.Pointer(c.AllocaCStr("-1099511627776tail")))
		var end *nativeChar
		if got := c.Strtoimax(input, &end, 10); got != -(1 << 40) {
			t.Fatalf("strtoimax = %d, want %d", got, -(1 << 40))
		}
		if got := c.GoString((*c.Char)(unsafe.Pointer(end))); got != "tail" {
			t.Fatalf("strtoimax end pointer = %q, want tail", got)
		}
	})
	t.Run("unsigned parsing", func(t *testing.T) {
		input := (*nativeChar)(unsafe.Pointer(c.AllocaCStr("0x10000000000rest")))
		var end *nativeChar
		if got := c.Strtoumax(input, &end, 0); got != 1<<40 {
			t.Fatalf("strtoumax = %d, want %d", got, 1<<40)
		}
		if got := c.GoString((*c.Char)(unsafe.Pointer(end))); got != "rest" {
			t.Fatalf("strtoumax end pointer = %q, want rest", got)
		}
	})
	t.Run("variadic pointer arguments", func(t *testing.T) {
		var first, second c.Int
		input := (*nativeChar)(unsafe.Pointer(c.AllocaCStr("17 -23")))
		format := (*nativeChar)(unsafe.Pointer(c.AllocaCStr("%d %d")))
		if count := c.Sscanf(input, format, &first, &second); count != 2 || first != 17 || second != -23 {
			t.Fatalf("sscanf = %d (%d,%d), want 2 (17,-23)", count, first, second)
		}
	})
}

func TestNativeLibcFnmatch(t *testing.T) {
	cases := []struct {
		name, pattern, value string
		flags, want          c.Int
	}{
		{"match", "*.go", "binding.go", 0, 0},
		{"mismatch", "*.go", "binding.c", 0, c.FNM_NOMATCH},
		{"pathname", "*", "dir/file", c.FNM_PATHNAME, c.FNM_NOMATCH},
		{"pathname disabled", "*", "dir/file", 0, 0},
		{"leading period", "*", ".private", c.FNM_PERIOD, c.FNM_NOMATCH},
		{"escaped wildcard", "\\*", "*", 0, 0},
		{"escaping disabled", "\\*", "*", c.FNM_NOESCAPE, c.FNM_NOMATCH},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pattern := (*nativeChar)(unsafe.Pointer(c.AllocaCStr(tc.pattern)))
			value := (*nativeChar)(unsafe.Pointer(c.AllocaCStr(tc.value)))
			if got := c.Fnmatch(pattern, value, tc.flags); got != tc.want {
				t.Fatalf("fnmatch(%q,%q,%d) = %d, want %d", tc.pattern, tc.value, tc.flags, got, tc.want)
			}
		})
	}
}

func TestNativeLibcStdio(t *testing.T) {
	stream := c.Tmpfile()
	if stream == nil {
		t.Fatal("tmpfile returned nil")
	}
	defer func() {
		if got := c.Fclose(stream); got != 0 {
			t.Errorf("fclose = %d, want 0", got)
		}
	}()
	payload := []byte{'a', 0, 'b', '\n'}
	if count := c.Fwrite(unsafe.Pointer(&payload[0]), 1, uintptr(len(payload)), stream); count != uintptr(len(payload)) {
		t.Fatalf("fwrite = %d, want %d", count, len(payload))
	}
	if got := c.Ftell(stream); got != c.Long(len(payload)) {
		t.Fatalf("ftell = %d, want %d", got, len(payload))
	}
	// SEEK_SET is zero in each of the supported target headers.
	if got := c.Fseek(stream, 0, 0); got != 0 {
		t.Fatalf("fseek = %d, want 0", got)
	}
	var position nativeFilePosition
	if got := c.Fgetpos(stream, &position); got != 0 {
		t.Fatalf("fgetpos = %d, want 0", got)
	}
	if got := c.Fgetc(stream); got != 'a' {
		t.Fatalf("first fgetc = %d, want a", got)
	}
	if got := c.Fsetpos(stream, &position); got != 0 {
		t.Fatalf("fsetpos = %d, want 0", got)
	}
	if got := c.Fgetc(stream); got != 'a' {
		t.Fatalf("restored fgetc = %d, want a", got)
	}
	c.Rewind(stream)
	got := make([]byte, len(payload))
	if count := c.Fread(unsafe.Pointer(&got[0]), 1, uintptr(len(got)), stream); count != uintptr(len(got)) {
		t.Fatalf("fread = %d, want %d", count, len(got))
	}
	if string(got) != string(payload) {
		t.Fatalf("read payload = %v, want %v", got, payload)
	}
	// EOF is -1 in each of the supported target headers.
	if value := c.Fgetc(stream); value != -1 {
		t.Fatalf("fgetc at EOF = %d, want -1", value)
	}
	if c.Feof(stream) == 0 || c.Ferror(stream) != 0 {
		t.Fatalf("EOF/error indicators = %d/%d", c.Feof(stream), c.Ferror(stream))
	}
	c.Clearerr(stream)
	if c.Feof(stream) != 0 || c.Ferror(stream) != 0 {
		t.Fatalf("clearerr indicators = %d/%d", c.Feof(stream), c.Ferror(stream))
	}
	c.Rewind(stream)
	if value := c.Fgetc(stream); value != 'a' {
		t.Fatalf("rewound fgetc = %d, want a", value)
	}
}
