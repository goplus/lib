//go:build arm64 && !darwin && !windows

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

// VaList is the AAPCS64 va_list on arm64 for Linux, Android and the BSDs.
// In C it is a struct `__va_list` passed by value, but a 32-byte struct is
// passed by an invisible pointer to a caller-made copy, so at the ABI level a
// C parameter of type va_list maps to *VaList (the callee consumes the copy,
// not the caller's original, unlike the array form on amd64 and s390x).
//
// Apple (darwin/ios) and Windows use a plain pointer on arm64; see VaList in
// valist_ptr.go.
type VaList struct {
	Stack  unsafe.Pointer // next stack argument
	GrTop  unsafe.Pointer // end of the general-purpose register save area
	VrTop  unsafe.Pointer // end of the vector register save area
	GrOffs int32          // negative offset from GrTop to the next GP register slot
	VrOffs int32          // negative offset from VrTop to the next vector register slot
}
