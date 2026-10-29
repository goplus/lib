//go:build (amd64 && (!windows || mingw)) || (arm64 && !darwin && !windows) || riscv64 || loong64 || s390x || mips64 || mips64le || wasm || ((ppc64 || ppc64le) && linux)

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

// LongDouble is an opaque 16-byte value on targets where sizeof(long double)
// is 16. The underlying format depends on the target:
//   - x86-64 on Linux/macOS/BSD and MinGW on Windows: x87 80-bit extended
//     precision (10 significant bytes + 6 bytes padding)
//   - arm64 on Linux/Android/BSD, riscv64, loong64, s390x, mips64/mips64le,
//     wasm: IEEE binary128
//   - linux ppc64/ppc64le: IBM double-double or IEEE binary128, depending on
//     the toolchain
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
