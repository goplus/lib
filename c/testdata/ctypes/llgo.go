//go:build llgo

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package ctypes

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoFiles = "_wrap/types.c"

//go:linkname SizeBytes C.lib_ctypes_size_bytes
func SizeBytes() c.Uint

//go:linkname LongBytes C.lib_ctypes_long_bytes
func LongBytes() c.Uint

//go:linkname WcharBytes C.lib_ctypes_wchar_bytes
func WcharBytes() c.Uint

//go:linkname CharMinusOne C.lib_ctypes_char_minus_one
func CharMinusOne() uint64

//go:linkname WcharMinusOne C.lib_ctypes_wchar_minus_one
func WcharMinusOne() uint64

//go:linkname EchoSize C.lib_ctypes_echo_size
func EchoSize(c.SizeT) c.SizeT

//go:linkname EchoSigned C.lib_ctypes_echo_signed
func EchoSigned(c.SsizeT) c.SsizeT

//go:linkname EchoLong C.lib_ctypes_echo_long
func EchoLong(c.Long) c.Long

//go:linkname EchoPtrdiff C.lib_ctypes_echo_ptrdiff
func EchoPtrdiff(c.PtrdiffT) c.PtrdiffT

//go:linkname EchoIntptr C.lib_ctypes_echo_intptr
func EchoIntptr(c.IntptrT) c.IntptrT

//go:linkname EchoUintptr C.lib_ctypes_echo_uintptr
func EchoUintptr(c.UintptrT) c.UintptrT

//go:linkname EchoWchar C.lib_ctypes_echo_wchar
func EchoWchar(c.WcharT) c.WcharT

//go:linkname CheckWords C.lib_ctypes_check_words
func CheckWords(*[2]c.SizeT, uintptr) c.Int
