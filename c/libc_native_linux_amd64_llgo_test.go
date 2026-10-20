//go:build llgo && linux && amd64 && !android

package c_test

import "github.com/goplus/lib/c"

type nativeChar = c.Char
type nativeFilePosition = c.FposT

const nativeFilePositionSize = 16
