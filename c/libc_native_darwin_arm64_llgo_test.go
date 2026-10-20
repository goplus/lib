//go:build llgo && darwin && arm64 && !ios

package c_test

import (
	"github.com/goplus/lib/c"
	"testing"
)

type nativeChar = c.Char
type nativeFilePosition = c.LongLong

const nativeFilePositionSize = 8

func TestGeneratedDarwinLongDoubleResult(t *testing.T) {
	var end *c.Char
	value := c.Strtold(c.Str("1.5"), &end)
	if value != 1.5 || end == nil || *end != 0 {
		t.Fatalf("strtold value/end = %v/%p", value, end)
	}
}
