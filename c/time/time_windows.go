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

package time

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const LLGoPackage = "link"

type TimeT int64

// Tm matches the Universal CRT struct tm. The Unix-only tm_gmtoff and tm_zone
// extensions are intentionally absent on Windows.
type Tm struct {
	Sec   c.Int
	Min   c.Int
	Hour  c.Int
	Mday  c.Int
	Mon   c.Int
	Year  c.Int
	Wday  c.Int
	Yday  c.Int
	Isdst c.Int
}

//go:linkname Time C._time64
func Time(timer *TimeT) TimeT

//go:linkname Mktime C._mktime64
func Mktime(timer *Tm) TimeT

//go:linkname Ctime C._ctime64
func Ctime(timer *TimeT) *c.Char

//go:linkname Difftime C._difftime64
func Difftime(end, start TimeT) float64

//go:linkname Gmtime C._gmtime64
func Gmtime(timer *TimeT) *Tm

//go:linkname Localtime C._localtime64
func Localtime(timer *TimeT) *Tm

//go:linkname Strftime C.strftime
func Strftime(buf *c.Char, bufSize uintptr, format *c.Char, timeptr *Tm) uintptr

type ClockT c.Long

//go:linkname Clock C.clock
func Clock() ClockT
