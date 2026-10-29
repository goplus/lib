//go:build sizetest

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

/*
#include <stddef.h>     // size_t, ptrdiff_t, wchar_t
#include <stdint.h>     // intptr_t, uintptr_t, intmax_t, uintmax_t
#include <sys/types.h>  // ssize_t

#if defined(__cplusplus) && __cplusplus >= 201103L
// char16_t / char32_t are keywords
#elif defined(__STDC_VERSION__) && __STDC_VERSION__ >= 202311L
// C23: keywords
#elif defined(__has_include)
#  if __has_include(<uchar.h>)
#    include <uchar.h>
#  endif
#endif

#define SizeChar    sizeof(char)
#define SizeChar8T  sizeof(char8_t)
#define SizeChar16T sizeof(char16_t)
#define SizeChar32T sizeof(char32_t)
#define SizeWcharT  sizeof(wchar_t)

#define SizeFloat  sizeof(float)
#define SizeDouble sizeof(double)
#define SizeLongDouble sizeof(long double)

#define SizePointer sizeof(void*)

#define SizeInt  sizeof(int)
#define SizeUint sizeof(unsigned int)

#define SizeLong  sizeof(long)
#define SizeUlong sizeof(unsigned long)

#define SizeLongLong  sizeof(long long)
#define SizeUlongLong sizeof(unsigned long long)

#define SizeSizeT   sizeof(size_t)
#define SizeSsizeT  sizeof(ssize_t)

#define SizePtrdiffT sizeof(ptrdiff_t)
#define SizeIntptrT  sizeof(intptr_t)
#define SizeUintptrT sizeof(uintptr_t)

#define SizeIntmaxT  sizeof(intmax_t)
#define SizeUintmaxT sizeof(uintmax_t)
*/
import "C"
import (
	"testing"
	"unsafe"
)

func SizeTest(t *testing.T) {
	assertSize(t, "char", unsafe.Sizeof(Char(0)), C.SizeChar)
	assertSize(t, "char8_t", unsafe.Sizeof(Char8T(0)), C.SizeChar8T)
	assertSize(t, "char16_t", unsafe.Sizeof(Char16T(0)), C.SizeChar16T)
	assertSize(t, "char32_t", unsafe.Sizeof(Char32T(0)), C.SizeChar32T)
	assertSize(t, "wchar_t", unsafe.Sizeof(WcharT(0)), C.SizeWcharT)

	assertSize(t, "float", unsafe.Sizeof(Float(0)), C.SizeFloat)
	assertSize(t, "double", unsafe.Sizeof(Double(0)), C.SizeDouble)
	assertSize(t, "long double", unsafe.Sizeof(LongDouble(0)), C.SizeLongDouble)

	assertSize(t, "void*", unsafe.Sizeof(Pointer(nil)), C.SizePointer)

	assertSize(t, "int", unsafe.Sizeof(Int(0)), C.SizeInt)
	assertSize(t, "unsigned int", unsafe.Sizeof(Uint(0)), C.SizeUint)

	assertSize(t, "long", unsafe.Sizeof(Long(0)), C.SizeLong)
	assertSize(t, "unsigned long", unsafe.Sizeof(Ulong(0)), C.SizeUlong)

	assertSize(t, "long long", unsafe.Sizeof(LongLong(0)), C.SizeLongLong)
	assertSize(t, "unsigned long long", unsafe.Sizeof(UlongLong(0)), C.SizeUlongLong)

	assertSize(t, "size_t", unsafe.Sizeof(SizeT(0)), C.SizeSizeT)
	assertSize(t, "ssize_t", unsafe.Sizeof(SsizeT(0)), C.SizeSsizeT)

	assertSize(t, "ptrdiff_t", unsafe.Sizeof(PtrdiffT(0)), C.SizePtrdiffT)
	assertSize(t, "intptr_t", unsafe.Sizeof(IntptrT(0)), C.SizeIntptrT)
	assertSize(t, "uintptr_t", unsafe.Sizeof(UintptrT(0)), C.SizeUintptrT)

	assertSize(t, "intmax_t", unsafe.Sizeof(IntmaxT(0)), C.SizeIntmaxT)
	assertSize(t, "uintmax_t", unsafe.Sizeof(UintmaxT(0)), C.SizeUintmaxT)
}

func assertSize(t *testing.T, name string, got uintptr, want C.size_t) {
	if int(got) != int(want) {
		t.Errorf("Size of %s: got %v, want %v", name, got, want)
	}
}
