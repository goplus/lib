//go:build amd64 && !windows

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

// VaList is the System V x86-64 va_list element (__va_list_tag), used on
// Linux, macOS and the BSDs. C declares va_list as `__va_list_tag[1]`, so as
// a function parameter it decays to a pointer: a C parameter of type va_list
// maps to *VaList, and a local va_list variable maps to VaList.
//
// Windows x64 uses a plain pointer; see VaList in valist_ptr.go.
type VaList struct {
	GpOffset        uint32         // offset of the next general-purpose register slot in RegSaveArea
	FpOffset        uint32         // offset of the next SSE register slot in RegSaveArea
	OverflowArgArea unsafe.Pointer // next argument passed on the stack
	RegSaveArea     unsafe.Pointer // register save area
}
