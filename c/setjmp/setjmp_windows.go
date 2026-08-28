//go:build windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package setjmp

// #include <setjmp.h>
import "C"

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const LLGoPackage = "decl"

type JmpBuf = C.jmp_buf

// Windows has no signal-mask variant of setjmp. These LLGo intrinsics lower
// at the caller, which is required because setjmp cannot safely be hidden in a
// wrapper frame that returns before longjmp restores it.

//go:linkname Setjmp llgo.setjmp
func Setjmp(env *JmpBuf) c.Int

//go:linkname Longjmp llgo.longjmp
func Longjmp(env *JmpBuf, val c.Int)
