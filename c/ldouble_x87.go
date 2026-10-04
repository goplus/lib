//go:build amd64 && (!windows || mingw)

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

// LongDouble is the x87 80-bit extended precision format, stored in 16 bytes
// (10 significant bytes + 6 bytes padding). Applies to x86-64 on
// Linux/macOS/BSD and MinGW on Windows.
//
// The C ABI requires 16-byte alignment, which Go cannot express (Go caps
// alignment at 8); the zero-size field below gives 8-byte alignment at most.
// When embedding LongDouble in a mirrored C struct, add explicit padding so
// the field offset is a multiple of 16.
type LongDouble struct {
	_    [0]uint64
	Data [16]byte
}

// LongDoubleSize is sizeof(long double) in bytes on this target.
const LongDoubleSize = 16
