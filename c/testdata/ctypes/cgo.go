//go:build !llgo && cgo

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package ctypes

/*
#include "_wrap/types.c"
*/
import "C"

import (
	"unsafe"

	"github.com/goplus/lib/c"
)

func SizeBytes() c.Uint  { return c.Uint(C.lib_ctypes_size_bytes()) }
func LongBytes() c.Uint  { return c.Uint(C.lib_ctypes_long_bytes()) }
func WcharBytes() c.Uint { return c.Uint(C.lib_ctypes_wchar_bytes()) }

func CharMinusOne() uint64  { return uint64(C.lib_ctypes_char_minus_one()) }
func WcharMinusOne() uint64 { return uint64(C.lib_ctypes_wchar_minus_one()) }

func EchoSize(value c.SizeT) c.SizeT {
	return c.SizeT(C.lib_ctypes_echo_size(C.size_t(value)))
}

func EchoSigned(value c.SsizeT) c.SsizeT {
	return c.SsizeT(C.lib_ctypes_echo_signed(C.lib_ctypes_ssize_t(value)))
}

func EchoLong(value c.Long) c.Long {
	return c.Long(C.lib_ctypes_echo_long(C.long(value)))
}

func EchoPtrdiff(value c.PtrdiffT) c.PtrdiffT {
	return c.PtrdiffT(C.lib_ctypes_echo_ptrdiff(C.ptrdiff_t(value)))
}

func EchoIntptr(value c.IntptrT) c.IntptrT {
	return c.IntptrT(C.lib_ctypes_echo_intptr(C.intptr_t(value)))
}

func EchoUintptr(value c.UintptrT) c.UintptrT {
	return c.UintptrT(C.lib_ctypes_echo_uintptr(C.uintptr_t(value)))
}

func EchoWchar(value c.WcharT) c.WcharT {
	return c.WcharT(C.lib_ctypes_echo_wchar(C.wchar_t(value)))
}

func CheckWords(words *[2]c.SizeT, stride uintptr) c.Int {
	return c.Int(C.lib_ctypes_check_words((*C.size_t)(unsafe.Pointer(words)), C.size_t(stride)))
}
