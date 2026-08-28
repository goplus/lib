//go:build !windows

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

package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

//go:linkname Getcwd C.getcwd
func Getcwd(buffer c.Pointer, size uintptr) *c.Char

//go:linkname Open C.open
func Open(path *c.Char, flags c.Int, __llgo_va_list ...any) c.Int

//go:linkname Fcntl C.fcntl
func Fcntl(a c.Int, b c.Int, __llgo_va_list ...any) c.Int

//go:linkname Close C.close
func Close(fd c.Int) c.Int

//go:linkname Read C.read
func Read(fd c.Int, buf c.Pointer, count uintptr) int

//go:linkname Write C.write
func Write(fd c.Int, buf c.Pointer, count uintptr) int

//go:linkname Execlp C.execlp
func Execlp(file *c.Char, arg0 *c.Char, __llgo_va_list ...any) c.Int
