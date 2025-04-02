package main

import (
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/cpp/std"
)

func main() {
	c.Printf(c.Str("sizeof(std::string) = %lu\n"), unsafe.Sizeof(std.String{}))
}
