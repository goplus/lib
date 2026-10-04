//go:build (windows && !(mingw && (amd64 || 386))) || (darwin && arm64) || arm || mips || mipsle || ((ppc64 || ppc64le) && !linux)

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
//   - Windows MSVC ABI (all architectures), and windows/arm64 (MSVC or MinGW)
//   - darwin/arm64 and ios/arm64
//   - 32-bit ARM, MIPS o32
//   - non-Linux ppc64 (AIX, FreeBSD, OpenBSD: 64-bit long double)
//
// Build with -tags mingw to select the MinGW x86 layout on windows/amd64
// and windows/386 instead.
type LongDouble float64

// LongDoubleSize is sizeof(long double) in bytes on this target.
const LongDoubleSize = 8
