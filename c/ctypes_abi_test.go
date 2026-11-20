//go:build llgo || cgo

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c_test

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/testdata/ctypes"
)

func cTypeTestValue(bytes c.Uint) uint64 {
	if bytes == 8 {
		return uint64(1)<<40 | 7
	}
	return 7
}

// The same assertions run through cgo and LLGo. See testdata/ctypes/README.md
// for the Wasm runners; ordinary native go test and llgo test discover this test.
func TestCTypeABI(t *testing.T) {
	value := cTypeTestValue(ctypes.SizeBytes())
	if got := uint64(ctypes.EchoSize(c.SizeT(value))); got != value {
		t.Fatalf("size_t round trip: got %d, want %d", got, value)
	}
	if got := int64(ctypes.EchoSigned(c.SsizeT(-int64(value)))); got != -int64(value) {
		t.Errorf("ssize_t round trip: got %d, want %d", got, -int64(value))
	}
	longValue := -int64(cTypeTestValue(ctypes.LongBytes()))
	if got := int64(ctypes.EchoLong(c.Long(longValue))); got != longValue {
		t.Errorf("long round trip: got %d, want %d", got, longValue)
	}
	if got := int64(ctypes.EchoPtrdiff(c.PtrdiffT(-int64(value)))); got != -int64(value) {
		t.Errorf("ptrdiff_t round trip: got %d, want %d", got, -int64(value))
	}
	if got := int64(ctypes.EchoIntptr(c.IntptrT(-int64(value)))); got != -int64(value) {
		t.Errorf("intptr_t round trip: got %d, want %d", got, -int64(value))
	}
	if got := uint64(ctypes.EchoUintptr(c.UintptrT(value))); got != value {
		t.Errorf("uintptr_t round trip: got %d, want %d", got, value)
	}

	wcharValue := uint32(0x10ffff)
	if ctypes.WcharBytes() == 2 {
		wcharValue = 0xfffd
	}
	if got := uint32(ctypes.EchoWchar(c.WcharT(wcharValue))); got != wcharValue {
		t.Errorf("wchar_t round trip: got %d, want %d", got, wcharValue)
	}
	var ch c.Char
	ch--
	if got, want := uint64(ch), ctypes.CharMinusOne(); got != want {
		t.Errorf("char signedness: got %d, want %d", got, want)
	}
	var wc c.WcharT
	wc--
	if got, want := uint64(wc), ctypes.WcharMinusOne(); got != want {
		t.Errorf("wchar_t signedness: got %d, want %d", got, want)
	}

	words := [2]c.SizeT{c.SizeT(value), c.SizeT(value + 2)}
	// Observe the generated storage stride, independently of unsafe.Sizeof's
	// frontend constant folding. The C helper checks it before accessing memory.
	stride := uintptr(unsafe.Pointer(&words[1])) - uintptr(unsafe.Pointer(&words[0]))
	if ctypes.CheckWords(&words, stride) != 1 || uint64(words[1]) != value+4 {
		t.Errorf("size_t array storage: got %v", words)
	}
}
