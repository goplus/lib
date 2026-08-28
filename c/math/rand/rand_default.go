//go:build !windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package rand

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const LLGoPackage = "decl"

//go:linkname RandR C.rand_r
func RandR(*c.Uint) c.Int

//go:linkname Sranddev C.sranddev
func Sranddev()

//go:linkname Random C.random
func Random() c.Long

//go:linkname Srandom C.srandom
func Srandom(c.Uint)

//go:linkname Srandomdev C.srandomdev
func Srandomdev()
