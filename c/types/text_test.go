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

package types

import (
	"math/big"
	"math/rand"
	"strings"
	"testing"
)

func unsignedBig(u Uint128) *big.Int {
	b := new(big.Int).SetUint64(u.hi)
	b.Lsh(b, 64)
	return b.Or(b, new(big.Int).SetUint64(u.lo))
}

func signedBig(i Int128) *big.Int {
	b := unsignedBig(Uint128(i))
	if i.Sign() < 0 {
		b.Sub(b, new(big.Int).Lsh(big.NewInt(1), 128))
	}
	return b
}

func TestTextRoundTrip(t *testing.T) {
	values := []Uint128{{}, {0, 1}, {0, maxUint64}, {1, 0}, {maxInt64, maxUint64}, {signBit, 0}, {maxUint64, maxUint64}}
	rng := rand.New(rand.NewSource(128))
	for n := 0; n < 256; n++ {
		values = append(values, Uint128{rng.Uint64(), rng.Uint64()})
	}
	for _, u := range values {
		for base := 2; base <= 62; base++ {
			want := unsignedBig(u).Text(base)
			if got := u.Text(base); got != want {
				t.Fatalf("unsigned base %d: got %s want %s", base, got, want)
			}
			if got, err := ParseUint128(want, base); err != nil || got != u {
				t.Fatalf("unsigned round trip base %d: %v", base, err)
			}
			i := Int128(u)
			want = signedBig(i).Text(base)
			if got := i.Text(base); got != want {
				t.Fatalf("signed base %d: got %s want %s", base, got, want)
			}
			if got, err := ParseInt128(want, base); err != nil || got != i {
				t.Fatalf("signed round trip base %d: %v", base, err)
			}
		}
		if u.String() != FormatUint128(u, 10) || Int128(u).String() != FormatInt128(Int128(u), 10) {
			t.Fatal("String/Format mismatch")
		}
	}
}

func checkParse(t *testing.T, s string, base int) {
	t.Helper()
	b, valid := new(big.Int).SetString(s, base)
	gotU, errU := ParseUint128(s, base)
	gotI, errI := ParseInt128(s, base)
	if !valid {
		if errU == nil || errI == nil || !gotU.IsZero() || !gotI.IsZero() {
			t.Fatalf("accepted invalid input %q base %d", s, base)
		}
		return
	}
	maxU := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
	maxI := new(big.Int).Rsh(new(big.Int).Set(maxU), 1)
	minI := new(big.Int).Neg(new(big.Int).Add(new(big.Int).Set(maxI), big.NewInt(1)))
	wantU := new(big.Int).Set(b)
	validU := b.Sign() >= 0 && b.Cmp(maxU) <= 0
	if b.Sign() < 0 {
		wantU.SetInt64(0)
	} else if b.Cmp(maxU) > 0 {
		wantU.Set(maxU)
	}
	wantI := new(big.Int).Set(b)
	validI := b.Cmp(minI) >= 0 && b.Cmp(maxI) <= 0
	if b.Cmp(minI) < 0 {
		wantI.Set(minI)
	} else if b.Cmp(maxI) > 0 {
		wantI.Set(maxI)
	}
	if (errU == nil) != validU || unsignedBig(gotU).Cmp(wantU) != 0 {
		t.Fatalf("uint parse %q base %d: got %s, %v; want %s valid %v", s, base, gotU, errU, wantU, validU)
	}
	if (errI == nil) != validI || signedBig(gotI).Cmp(wantI) != 0 {
		t.Fatalf("int parse %q base %d: got %s, %v; want %s valid %v", s, base, gotI, errI, wantI, validI)
	}
}

func TestParseSyntaxAndLimits(t *testing.T) {
	inputs := []string{"", "+", "-", "0", "-0", "+0", "00", "08", "0_7", "0_", "0__7", "0b", "0b_", "0b_10", "0x_FF", "0Xff", "0o_77", "_1", "1_", "1__2", "1_2", " 12", "12 ", "++1", "-0x80", "Z", "aZ", "0xg", "１２", "\x00", strings.Repeat("9", 200), strings.Repeat("9", 200) + "!"}
	for bits := uint(0); bits <= 130; bits++ {
		b := new(big.Int).Lsh(big.NewInt(1), bits)
		for delta := int64(-1); delta <= 1; delta++ {
			n := new(big.Int).Add(b, big.NewInt(delta))
			inputs = append(inputs, n.String(), new(big.Int).Neg(n).String())
		}
	}
	for _, s := range inputs {
		for _, base := range []int{0, 2, 8, 10, 16, 36, 37, 62} {
			checkParse(t, s, base)
		}
	}
}

func TestInvalidBases(t *testing.T) {
	for _, base := range []int{-1, 1, 63} {
		for _, f := range []func(){func() { ParseUint128("1", base) }, func() { ParseInt128("1", base) }, func() { (Uint128{}).Text(base) }, func() { (Int128{}).Text(base) }} {
			func() {
				defer func() {
					if recover() == nil {
						t.Errorf("base %d did not panic", base)
					}
				}()
				f()
			}()
		}
	}
}

func FuzzParse(f *testing.F) {
	for _, s := range []string{"", "0x_FF", "-0", "0__7", "340282366920938463463374607431768211455", "-170141183460469231731687303715884105728"} {
		f.Add(s, uint8(0))
	}
	f.Fuzz(func(t *testing.T, s string, b uint8) {
		base := int(b) % 63
		if base == 1 {
			base = 0
		}
		checkParse(t, s, base)
	})
}
