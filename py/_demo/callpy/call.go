package main

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/py"
)

func main() {
	py.Initialize()
	py.SetProgramName(*c.Argv)
	math := py.ImportModule(c.Str("math"))
	sqrt := math.GetAttrString(c.Str("sqrt"))
	sqrt2 := sqrt.CallOneArg(py.Float(2))
	c.Printf(c.Str("sqrt(2) = %f\n"), sqrt2.Float64())
	sqrt2.DecRef()
	sqrt.DecRef()
	math.DecRef()
	py.Finalize()
}
