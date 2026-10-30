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

#include <stdlib.h>
#include <stdarg.h>
#include <locale.h>

#ifndef _WIN32
#include <iconv.h>
#include <getopt.h>

#define SizeIconvT  sizeof(iconv_t)
#define SizeOption  sizeof(struct option)
#else

#define SizeIconvT  sizeof(void*)
#define SizeOption  0
#endif

#define SizeChar    sizeof(char)
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

#define SizeDivT    sizeof(div_t)
#define SizeLdivT   sizeof(ldiv_t)
#define SizeLldivT  sizeof(lldiv_t)
#define SizeLocaleT sizeof(locale_t)
#define SizeVaList  sizeof(va_list)
*/
import "C"
import (
	"testing"
	"unsafe"
)

func SizeTest(t *testing.T) {
	assert(t, "sizeof(char)", unsafe.Sizeof(Char(0)), C.SizeChar)
	assert(t, "sizeof(wchar_t)", unsafe.Sizeof(WcharT(0)), C.SizeWcharT)

	var ld LongDouble
	assert(t, "sizeof(float)", unsafe.Sizeof(Float(0)), C.SizeFloat)
	assert(t, "sizeof(double)", unsafe.Sizeof(Double(0)), C.SizeDouble)
	assert(t, "sizeof(long double)", unsafe.Sizeof(ld), C.SizeLongDouble)

	assert(t, "sizeof(void*)", unsafe.Sizeof(Pointer(nil)), C.SizePointer)

	assert(t, "sizeof(int)", unsafe.Sizeof(Int(0)), C.SizeInt)
	assert(t, "sizeof(unsigned int)", unsafe.Sizeof(Uint(0)), C.SizeUint)

	assert(t, "sizeof(long)", unsafe.Sizeof(Long(0)), C.SizeLong)
	assert(t, "sizeof(unsigned long)", unsafe.Sizeof(Ulong(0)), C.SizeUlong)

	assert(t, "sizeof(long long)", unsafe.Sizeof(LongLong(0)), C.SizeLongLong)
	assert(t, "sizeof(unsigned long long)", unsafe.Sizeof(UlongLong(0)), C.SizeUlongLong)

	assert(t, "sizeof(size_t)", unsafe.Sizeof(SizeT(0)), C.SizeSizeT)
	assert(t, "sizeof(ssize_t)", unsafe.Sizeof(SsizeT(0)), C.SizeSsizeT)

	assert(t, "sizeof(ptrdiff_t)", unsafe.Sizeof(PtrdiffT(0)), C.SizePtrdiffT)
	assert(t, "sizeof(intptr_t)", unsafe.Sizeof(IntptrT(0)), C.SizeIntptrT)
	assert(t, "sizeof(uintptr_t)", unsafe.Sizeof(UintptrT(0)), C.SizeUintptrT)

	assert(t, "sizeof(intmax_t)", unsafe.Sizeof(IntmaxT(0)), C.SizeIntmaxT)
	assert(t, "sizeof(uintmax_t)", unsafe.Sizeof(UintmaxT(0)), C.SizeUintmaxT)

	assert(t, "sizeof(div_t)", unsafe.Sizeof(DivT{}), C.SizeDivT)
	assert(t, "sizeof(ldiv_t)", unsafe.Sizeof(LdivT{}), C.SizeLdivT)
	assert(t, "sizeof(lldiv_t)", unsafe.Sizeof(LldivT{}), C.SizeLldivT)
	assert(t, "sizeof(iconv_t)", unsafe.Sizeof(IconvT(nil)), C.SizeIconvT)
	assert(t, "sizeof(locale_t)", unsafe.Sizeof(LocaleT(nil)), C.SizeLocaleT)
	if C.SizeOption != 0 {
		assert(t, "sizeof(option)", unsafe.Sizeof(Option{}), C.SizeOption)
	}
	assert(t, "sizeof(va_list)", unsafe.Sizeof(VaList(nil)), C.SizeVaList)
}

func assert(t *testing.T, name string, got uintptr, want C.size_t) {
	if int(got) != int(want) {
		t.Helper()
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}
