//go:build js && wasm && llgo.wasm.emscripten.memory64
// +build js,wasm,llgo.wasm.emscripten.memory64

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c

// Emscripten Memory64 uses the wasm64 LP64 C data model.
type (
	Long  int64
	Ulong uint64
)

type (
	SizeT    Ulong
	SsizeT   Long
	PtrdiffT Long
	IntptrT  Long
	UintptrT Ulong
	WcharT   int32
)
