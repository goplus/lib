//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c_test

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
)

func TestNativeLibcLayouts(t *testing.T) {
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
