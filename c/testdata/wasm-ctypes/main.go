//go:build llgo && wasm

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package main

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoFiles = "_wrap/types.c"

//go:linkname sizeBytes C.lib_ctypes_size_bytes
func sizeBytes() c.Uint

//go:linkname echoSize C.lib_ctypes_echo_size
func echoSize(c.SizeT) c.SizeT

//go:linkname echoSigned C.lib_ctypes_echo_signed
func echoSigned(c.SsizeT) c.SsizeT

//go:linkname echoLong C.lib_ctypes_echo_long
func echoLong(c.Long) c.Long

//go:linkname echoPtrdiff C.lib_ctypes_echo_ptrdiff
func echoPtrdiff(c.PtrdiffT) c.PtrdiffT

//go:linkname echoIntptr C.lib_ctypes_echo_intptr
func echoIntptr(c.IntptrT) c.IntptrT

//go:linkname echoUintptr C.lib_ctypes_echo_uintptr
func echoUintptr(c.UintptrT) c.UintptrT

//go:linkname echoWchar C.lib_ctypes_echo_wchar
func echoWchar(c.WcharT) c.WcharT

//go:linkname checkWords C.lib_ctypes_check_words
func checkWords(*[2]c.SizeT) c.Int

func main() {
	value := uint64(7)
	if sizeBytes() == 8 {
		value += uint64(1) << 40
	}
	if uint64(echoSize(c.SizeT(value))) != value {
		panic("size_t round trip")
	}
	if int64(echoSigned(c.SsizeT(-int64(value)))) != -int64(value) {
		panic("ssize_t round trip")
	}
	if int64(echoLong(c.Long(-int64(value)))) != -int64(value) {
		panic("long round trip")
	}
	if int64(echoPtrdiff(c.PtrdiffT(-int64(value)))) != -int64(value) {
		panic("ptrdiff_t round trip")
	}
	if int64(echoIntptr(c.IntptrT(-int64(value)))) != -int64(value) {
		panic("intptr_t round trip")
	}
	if uint64(echoUintptr(c.UintptrT(value))) != value {
		panic("uintptr_t round trip")
	}
	if echoWchar(c.WcharT(0x10ffff)) != c.WcharT(0x10ffff) {
		panic("wchar_t round trip")
	}
	words := [2]c.SizeT{c.SizeT(value), c.SizeT(value + 2)}
	if checkWords(&words) != 1 || uint64(words[1]) != value+4 {
		panic("size_t array storage")
	}
	println("C ABI round trip OK", sizeBytes(), value)
}
