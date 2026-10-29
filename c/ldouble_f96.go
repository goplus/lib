//go:build 386 && (!windows || mingw)

/*
 * Copyright (c) 2026 The XGo Authors (xgo.dev). All rights reserved.
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

package c

// LongDouble is the x87 80-bit extended precision format on 32-bit x86,
// stored in 12 bytes (10 significant bytes + 2 bytes padding) with 4-byte
// alignment (i386 System V ABI and MinGW i386). Applies to linux/386,
// the BSDs on 386, and windows/386 when built with -tags mingw.
//
// Windows MSVC on 386 uses a plain 8-byte double; see LongDouble in
// ldouble_float64.go.
type LongDouble struct {
	_    [0]uint32
	Data [12]byte
}

// LongDoubleSize is sizeof(long double) in bytes on this target.
const LongDoubleSize = 12
