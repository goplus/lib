//go:build (amd64 && windows) || 386 || arm || (arm64 && (darwin || windows)) || riscv64 || loong64 || mips || mipsle || mips64 || mips64le || ppc64 || ppc64le || wasm

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

import "unsafe"

// VaList is a single pointer (va_list is `char *` or `void *`) on these targets:
//   - Windows on amd64, 386 and arm64 (MSVC and MinGW)
//   - darwin/arm64 and ios/arm64 (Apple deviates from the AAPCS64 struct)
//   - 386 on every OS
//   - riscv64, loong64
//   - mips, mipsle, mips64, mips64le
//   - ppc64, ppc64le (all OSes; ELF and AIX both use a plain pointer)
//   - 32-bit ARM: C declares `struct __va_list { void *__ap; }`, a struct with
//     a single pointer, which has the same size and is passed in the same way
//     as a pointer
//   - wasm: C wasm32 uses a 4-byte `void *`, while Go's wasm pointers are
//     8 bytes, so VaListSize does not match the C size there
//
// Because va_list is a plain value here, a C parameter of type va_list maps
// to VaList itself (not *VaList).
type VaList unsafe.Pointer
