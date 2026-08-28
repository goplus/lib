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

package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const (
	LLGoPackage = "link"
	PATH_MAX    = 260
)

type (
	ModeT c.Int
	OffT  int64
	PidT  c.Int
)

const (
	O_RDONLY = 0x0000
	O_WRONLY = 0x0001
	O_RDWR   = 0x0002
	O_APPEND = 0x0008
	O_CREAT  = 0x0100
	O_TRUNC  = 0x0200
	O_EXCL   = 0x0400
	O_TEXT   = 0x4000
	O_BINARY = 0x8000

	EAGAIN = 11
)

//go:linkname errno C._errno
func errno() *c.Int

func Errno() c.Int {
	return *errno()
}

//go:linkname Remove C.remove
func Remove(path *c.Char) c.Int

//go:linkname Rename C.rename
func Rename(oldpath *c.Char, newpath *c.Char) c.Int

//go:linkname Chdir C._chdir
func Chdir(path *c.Char) c.Int

//go:linkname Getenv C.getenv
func Getenv(name *c.Char) *c.Char

//go:linkname Putenv C._putenv
func Putenv(env *c.Char) c.Int

//go:linkname Getcwd C._getcwd
func Getcwd(buffer c.Pointer, size c.Int) *c.Char

//go:linkname Open C._open
func Open(path *c.Char, flags c.Int, __llgo_va_list ...any) c.Int

//go:linkname Creat C._creat
func Creat(path *c.Char, mode ModeT) c.Int

//go:linkname Dup C._dup
func Dup(fd c.Int) c.Int

//go:linkname Dup2 C._dup2
func Dup2(oldfd c.Int, newfd c.Int) c.Int

//go:linkname Close C._close
func Close(fd c.Int) c.Int

//go:linkname Read C._read
func Read(fd c.Int, buf c.Pointer, count c.Uint) c.Int

//go:linkname Write C._write
func Write(fd c.Int, buf c.Pointer, count c.Uint) c.Int

//go:linkname Lseek C._lseeki64
func Lseek(fd c.Int, offset OffT, whence c.Int) OffT

//go:linkname Isatty C._isatty
func Isatty(fd c.Int) c.Int

//go:linkname Execl C._execl
func Execl(path *c.Char, arg0 *c.Char, __llgo_va_list ...any) c.Int

//go:linkname Execle C._execle
func Execle(path *c.Char, arg0 *c.Char, __llgo_va_list ...any) c.Int

//go:linkname Execlp C._execlp
func Execlp(file *c.Char, arg0 *c.Char, __llgo_va_list ...any) c.Int

//go:linkname Execv C._execv
func Execv(path *c.Char, argv **c.Char) c.Int

//go:linkname Execve C._execve
func Execve(path *c.Char, argv **c.Char, envp **c.Char) c.Int

//go:linkname Execvp C._execvp
func Execvp(file *c.Char, argv **c.Char) c.Int

//go:linkname Getpid C._getpid
func Getpid() PidT

//go:linkname Exit C.exit
func Exit(status c.Int)
