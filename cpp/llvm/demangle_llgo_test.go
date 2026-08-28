//go:build llgo

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

package llvm_test

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/cpp/llvm"
)

func demangledString(t *testing.T, value *c.Char) string {
	t.Helper()
	if value == nil {
		t.Fatal("demangler returned nil")
	}
	defer c.Free(unsafe.Pointer(value))
	return c.GoString(value)
}

func TestItaniumDemangle(t *testing.T) {
	if got := demangledString(t, llvm.ItaniumDemangle("_Z3foov", true)); got != "foo()" {
		t.Fatalf("ItaniumDemangle = %q, want %q", got, "foo()")
	}
}

func TestMicrosoftDemangle(t *testing.T) {
	var nRead uintptr
	var status c.Int
	got := demangledString(t, llvm.MicrosoftDemangle("?foo@@YAHH@Z", &nRead, &status, 0))
	if status != 0 || nRead == 0 || !strings.Contains(got, "foo") {
		t.Fatalf("MicrosoftDemangle = %q, nRead %d, status %d", got, nRead, status)
	}
}

func TestRustDemangle(t *testing.T) {
	if got := demangledString(t, llvm.RustDemangle("_RNvC6_123foo3bar")); got != "123foo::bar" {
		t.Fatalf("RustDemangle = %q, want %q", got, "123foo::bar")
	}
}
