//go:build windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package std

// String reserves the release-mode MSVC std::string layout: a 16-byte small
// string buffer followed by size and capacity fields.
type String struct {
	storage [16]byte
	words   [2]uintptr
}
