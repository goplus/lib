//go:build llgo && linux && arm64 && !android

package c_test

import "github.com/goplus/lib/c"

type nativeChar = uint8
type nativeFilePosition = c.FposT

const nativeFilePositionSize = 16
