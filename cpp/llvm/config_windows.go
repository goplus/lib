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

package llvm

const (
	// LLVM may itself be built with either the MSVC or libc++ C++ ABI. Compile
	// the bridge for the ABI selected by its host clang++, then expose only C
	// functions to the LLGo target. This avoids leaking that C++ ABI across the
	// package boundary while retaining the target's native C ABI.
	LLGoFiles = "$(pkg-config --cflags llvm-19) -std=c++17 --target=$(llvm-config --host-target): _wrap/demangle_windows.cpp"
	// The llvm-19 package preserves LLVM's library order and required Windows
	// system libraries.
	LLGoPackage = "link: $(pkg-config --libs llvm-19)"
)
