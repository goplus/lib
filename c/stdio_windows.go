//go:build windows
// +build windows

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

package c

import _ "unsafe"

// Standard-stream pointers require package initialization; declaration-only
// and link-only packages intentionally skip init functions. The LLGo Windows
// linker supplies legacy_stdio_definitions.lib for the out-of-line formatted
// stdio symbols shared with the other platforms.
const LLGoPackage = true

// The Universal CRT exposes standard streams through __acrt_iob_func rather
// than the Unix stdin/stdout/stderr data symbols.

//go:linkname acrtIobFunc C.__acrt_iob_func
func acrtIobFunc(index Uint) FilePtr

var (
	Stdin  = acrtIobFunc(0)
	Stdout = acrtIobFunc(1)
	Stderr = acrtIobFunc(2)
)
