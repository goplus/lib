//go:build llgo && darwin && arm64 && !ios

package c_test

import "github.com/goplus/lib/c"

type nativeChar = c.Char
type nativeFilePosition = c.LongLong

const nativeFilePositionSize = 8
