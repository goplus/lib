package main

import (
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/clang"
)

const (
	LLGoCFlags = "-I$(llvm-config --includedir)"
)

func main() {
	c.Printf(c.Str("sizeof(clang.Cursor) = %lu\n"), unsafe.Sizeof(clang.Cursor{}))
}
