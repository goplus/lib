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

package extra

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/goplus/lib/c/types"
)

func TestBigConversions(t *testing.T) {
	for _, s := range []string{"0", "1", "-1", "18446744073709551616", "340282366920938463463374607431768211455", "340282366920938463463374607431768211456", "170141183460469231731687303715884105727", "170141183460469231731687303715884105728", "-170141183460469231731687303715884105728", "-170141183460469231731687303715884105729"} {
		b, _ := new(big.Int).SetString(s, 10)
		u, okU := Uint128FromBig(b)
		wantU, errU := types.ParseUint128(s, 10)
		i, okI := Int128FromBig(b)
		wantI, errI := types.ParseInt128(s, 10)
		if u != wantU || okU != (errU == nil) || i != wantI || okI != (errI == nil) {
			t.Fatalf("conversion %s mismatch", s)
		}
		if b.String() != s {
			t.Fatal("source modified")
		}
		// Destination contents and sign must be fully replaced, including zero.
		for _, initial := range []int64{-99, 0, 99} {
			dst := big.NewInt(initial)
			if Uint128ToBig(dst, u) != dst || dst.String() != u.String() {
				t.Fatalf("unsigned destination reuse: %s", s)
			}
			dst.SetInt64(initial)
			if Int128ToBig(dst, i) != dst || dst.String() != i.String() {
				t.Fatalf("signed destination reuse: %s", s)
			}
			(Int128{i}).ToBigInt(dst)
			if dst.String() != i.String() {
				t.Fatal("signed wrapper")
			}
			(Uint128{u}).ToBigInt(dst)
			if dst.String() != u.String() {
				t.Fatal("unsigned wrapper")
			}
		}
	}
}

func TestFormattingAndScanning(t *testing.T) {
	for _, s := range []string{"0", "-123", "170141183460469231731687303715884105727", "-170141183460469231731687303715884105728"} {
		i, _ := types.ParseInt128(s, 10)
		wrapped := Int128{i}
		for _, format := range []string{"%v", "%d", "%+040d", "%#x", "%X", "%b", "%o", "% 45d", "%-45d", "%.50d"} {
			if got, want := fmt.Sprintf(format, wrapped), fmt.Sprintf(format, wrapped.BigInt()); got != want {
				t.Fatalf("%s: got %s want %s", format, got, want)
			}
		}
		var scanned Int128
		if n, err := fmt.Sscan(s, &scanned); err != nil || n != 1 || scanned.Int128 != i {
			t.Fatalf("scan %s: %v", s, err)
		}
	}
	u, _ := types.ParseUint128("340282366920938463463374607431768211455", 10)
	wrapped := Uint128{u}
	if fmt.Sprintf("%#040x", wrapped) != fmt.Sprintf("%#040x", wrapped.BigInt()) {
		t.Fatal("unsigned format")
	}
	var scanned Uint128
	if _, err := fmt.Sscan(u.String(), &scanned); err != nil || scanned.Uint128 != u {
		t.Fatalf("unsigned scan: %v", err)
	}
	for _, s := range []string{"bad", "340282366920938463463374607431768211456", "-1"} {
		before := scanned
		if _, err := fmt.Sscan(s, &scanned); err == nil || scanned != before {
			t.Fatalf("invalid unsigned scan %s", s)
		}
	}
	i := Int128{types.Int128_Cast__2(7)}
	if _, err := fmt.Sscan("bad", &i); err == nil || i.Int128 != types.Int128_Cast__2(7) {
		t.Fatal("invalid signed scan changed value")
	}
}
