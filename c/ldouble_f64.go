//go:build (windows && (arm64 || msvc)) || (darwin && arm64) || arm || mips || mipsle || ((ppc64 || ppc64le) && !linux)

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

// LongDouble is identical to double on these targets (8 bytes, IEEE binary64):
//   - windows/arm64 (MSVC or MinGW)
//   - windows/amd64 and windows/386 when built with -tags msvc (MSVC ABI)
//   - darwin/arm64 and ios/arm64
//   - 32-bit ARM, MIPS o32
//   - non-Linux ppc64 (AIX, FreeBSD, OpenBSD: 64-bit long double)
//
// On windows/amd64 and windows/386 the default is the MinGW layout (the C
// toolchain cgo uses on Windows, and what GitHub's windows runners provide),
// where long double is the x87 80-bit format. Build with -tags msvc to select
// the 8-byte MSVC layout instead.
type LongDouble float64

// LongDoubleSize is sizeof(long double) in bytes on this target.
const LongDoubleSize = 8
