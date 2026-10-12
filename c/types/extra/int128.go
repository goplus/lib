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

// Package extra provides optional math/big and fmt adapters for c/types.
// Importing c/types alone does not load these dependencies.
package extra

import (
	"fmt"
	"math/big"

	"github.com/goplus/lib/c/types"
)

// Int128 wraps a core integer with fmt.Formatter and fmt.Scanner support.
// Construct it with Int128{Int128: value}; the embedded field is the core value.
// Arithmetic methods promoted from the core type return core values.
type Int128 struct{ types.Int128 }

// Uint128 is the unsigned counterpart of Int128.
type Uint128 struct{ types.Uint128 }

var (
	_ fmt.Formatter = Int128{}
	_ fmt.Formatter = Uint128{}
	_ fmt.Scanner   = (*Int128)(nil)
	_ fmt.Scanner   = (*Uint128)(nil)
)

// Uint128FromBig converts v, clamping to zero or MaxUint128 on overflow.
// The boolean reports whether v was in range. v must not be nil.
func Uint128FromBig(v *big.Int) (types.Uint128, bool) {
	if v.Sign() < 0 {
		return types.Uint128{}, false
	}
	if v.BitLen() > 128 {
		return types.Uint128_Cast__b(^uint64(0), ^uint64(0)), false
	}
	var high big.Int
	high.Rsh(v, 64)
	return types.Uint128_Cast__b(high.Uint64(), v.Uint64()), true
}

// Int128FromBig converts v, clamping to MinInt128 or MaxInt128 on overflow.
// The boolean reports whether v was in range. v must not be nil.
func Int128FromBig(v *big.Int) (types.Int128, bool) {
	var magnitude big.Int
	magnitude.Abs(v)
	u, ok := Uint128FromBig(&magnitude)
	limit := types.Uint128_Cast__b(1<<63-1, ^uint64(0))
	negative := v.Sign() < 0
	if negative {
		limit = types.Uint128_Cast__b(1<<63, 0)
	}
	if !ok || u.Cmp__1(limit) > 0 {
		return types.Int128_Cast__7(limit), false
	}
	i := types.Int128_Cast__7(u)
	if negative {
		i = i.XGo_Neg()
	}
	return i, true
}

// Uint128ToBig replaces dst with u and returns dst. dst must not be nil.
func Uint128ToBig(dst *big.Int, u types.Uint128) *big.Int {
	dst.SetUint64(u.XGo_Rsh(64).XGo_Rcast__1())
	dst.Lsh(dst, 64)
	var low big.Int
	low.SetUint64(u.XGo_Rcast__1())
	return dst.Or(dst, &low)
}

// Int128ToBig replaces dst with i and returns dst. dst must not be nil.
func Int128ToBig(dst *big.Int, i types.Int128) *big.Int {
	negative := i.Sign() < 0
	if negative {
		i = i.XGo_Neg()
	}
	Uint128ToBig(dst, i.XGo_Rcast__0())
	if negative {
		dst.Neg(dst)
	}
	return dst
}

func (i Int128) BigInt() *big.Int               { return Int128ToBig(new(big.Int), i.Int128) }
func (u Uint128) BigInt() *big.Int              { return Uint128ToBig(new(big.Int), u.Uint128) }
func (i Int128) ToBigInt(dst *big.Int)          { Int128ToBig(dst, i.Int128) }
func (u Uint128) ToBigInt(dst *big.Int)         { Uint128ToBig(dst, u.Uint128) }
func (i Int128) Format(s fmt.State, verb rune)  { i.BigInt().Format(s, verb) }
func (u Uint128) Format(s fmt.State, verb rune) { u.BigInt().Format(s, verb) }

// Scan retains the original decimal-token scanning behavior. Failed scans
// leave the destination unchanged.
func (i *Int128) Scan(s fmt.ScanState, verb rune) error {
	token, err := s.Token(true, nil)
	if err != nil {
		return err
	}
	v, err := types.ParseInt128(string(token), 10)
	if err == nil {
		i.Int128 = v
	}
	return err
}

func (u *Uint128) Scan(s fmt.ScanState, verb rune) error {
	token, err := s.Token(true, nil)
	if err != nil {
		return err
	}
	v, err := types.ParseUint128(string(token), 10)
	if err == nil {
		u.Uint128 = v
	}
	return err
}
