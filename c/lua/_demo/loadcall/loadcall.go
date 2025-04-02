package main

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/lua"
)

func main() {
	L := lua.Newstate__1()
	defer L.Close()

	L.Openlibs()
	if res := L.Loadstring(c.Str("print('hello world')")); res != lua.OK {
		println("error")
	}
	if res := L.Pcall(0, 0, 0); res != lua.OK {
		println("error")
	}

}

/* Expected output:
hello world
*/
