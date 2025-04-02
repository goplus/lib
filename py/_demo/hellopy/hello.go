package main

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/py"
)

func main() {
	py.Initialize()
	py.SetProgramName(*c.Argv)
	py.RunSimpleString(c.Str(`print('Hello, World!')`))
	py.Finalize()
}
