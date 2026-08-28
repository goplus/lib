//go:build windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llvm

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

// LLVM's official Windows binaries use the MSVC standard library, whose C++
// symbol names differ from the libc++ names used on Unix. Keep those names
// behind the C ABI implemented by demangle_windows.cpp.

//go:linkname itaniumDemangle C.llgoLLVMItaniumDemangle
func itaniumDemangle(data *c.Char, size uintptr, parseParams c.Int) *c.Char

func ItaniumDemangle(mangledName StringView, parseParams bool) *c.Char {
	parse := c.Int(0)
	if parseParams {
		parse = 1
	}
	return itaniumDemangle(c.GoStringData(mangledName), uintptr(len(mangledName)), parse)
}

//go:linkname microsoftDemangle C.llgoLLVMMicrosoftDemangle
func microsoftDemangle(data *c.Char, size uintptr, nRead *uintptr, status *c.Int, flags MSDemangleFlags) *c.Char

func MicrosoftDemangle(mangledName StringView, nRead *uintptr, status *c.Int, flags MSDemangleFlags) *c.Char {
	return microsoftDemangle(c.GoStringData(mangledName), uintptr(len(mangledName)), nRead, status, flags)
}

//go:linkname rustDemangle C.llgoLLVMRustDemangle
func rustDemangle(data *c.Char, size uintptr) *c.Char

func RustDemangle(mangledName StringView) *c.Char {
	return rustDemangle(c.GoStringData(mangledName), uintptr(len(mangledName)))
}
