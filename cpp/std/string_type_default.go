//go:build !windows

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package std

import "unsafe"

// String represents a libc++ std::string object.
type String struct {
	Unused [3 * unsafe.Sizeof(0)]byte
}
