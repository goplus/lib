//go:build s390x

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

// VaList is the s390x ELF va_list element (__va_list_tag). C declares
// va_list as `__va_list_tag[1]`, so as a function parameter it decays to a
// pointer: a C parameter of type va_list maps to *VaList, and a local va_list
// variable maps to VaList.
type VaList struct {
	Gpr             int64          // number of general-purpose argument registers already consumed
	Fpr             int64          // number of floating-point argument registers already consumed
	OverflowArgArea unsafe.Pointer // next argument passed on the stack
	RegSaveArea     unsafe.Pointer // register save area
}
