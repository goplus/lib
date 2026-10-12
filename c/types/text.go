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

import "math/bits"

const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// String returns the decimal representation of u.
func (u Uint128) String() string { return u.Text(10) }

// Text returns u in the given base, which must be between 2 and 62.
func (u Uint128) Text(base int) string {
	if base < 2 || base > len(digits) {
		panic("invalid base")
	}
	var buf [128]byte
	n := len(buf)
	for {
		var rem uint64
		u.hi, rem = bits.Div64(0, u.hi, uint64(base))
		u.lo, rem = bits.Div64(rem, u.lo, uint64(base))
		n--
		buf[n] = digits[rem]
		if u.IsZero() {
			return string(buf[n:])
		}
	}
}

// String returns the decimal representation of i.
func (i Int128) String() string { return i.Text(10) }

// Text returns i in the given base, which must be between 2 and 62.
func (i Int128) Text(base int) string {
	if i.Sign() < 0 {
		// Two's-complement negation also gives the unsigned magnitude of MinInt128.
		return "-" + Uint128(i.XGo_Neg()).Text(base)
	}
	return Uint128(i).Text(base)
}

func FormatInt128(i Int128, base int) string   { return i.Text(base) }
func FormatUint128(u Uint128, base int) string { return u.Text(base) }

type parseError struct {
	kind, text string
	rangeError bool
}

func (e *parseError) Error() string {
	if e.rangeError {
		return "string outside " + e.kind + " range: " + e.text
	}
	return "invalid " + e.kind + " string: " + e.text
}

// ParseUint128 parses an optionally signed integer. Base must be 0 or 2..62.
// Base 0 recognizes 0b, 0o, 0x and leading-zero octal prefixes and permits
// underscores between digits or immediately after a base prefix. Bases up to
// 36 ignore letter case; larger bases use a-z for 10..35 and A-Z for 36..61.
// Syntax errors return zero. Out-of-range values return the nearest endpoint
// and an error, including zero for a negative value other than negative zero.
func ParseUint128(s string, base int) (Uint128, error) {
	u, neg, overflow, valid := parseMagnitude(s, base)
	if !valid {
		return Uint128{}, &parseError{kind: "uint128", text: s}
	}
	if neg && (overflow || !u.IsZero()) {
		return Uint128{}, &parseError{kind: "uint128", text: s, rangeError: true}
	}
	if overflow {
		return Uint128{maxUint64, maxUint64}, &parseError{kind: "uint128", text: s, rangeError: true}
	}
	return u, nil
}

// ParseInt128 uses the same syntax as ParseUint128. Out-of-range values return
// MinInt128 or MaxInt128 and an error; syntax errors return zero and an error.
func ParseInt128(s string, base int) (Int128, error) {
	u, neg, overflow, valid := parseMagnitude(s, base)
	if !valid {
		return Int128{}, &parseError{kind: "int128", text: s}
	}
	limit := Uint128{maxInt64, maxUint64}
	if neg {
		limit = Uint128{signBit, 0}
	}
	if overflow || u.Cmp__1(limit) > 0 {
		if neg {
			return Int128{signBit, 0}, &parseError{kind: "int128", text: s, rangeError: true}
		}
		return Int128{maxInt64, maxUint64}, &parseError{kind: "int128", text: s, rangeError: true}
	}
	if neg {
		return Int128(u).XGo_Neg(), nil
	}
	return Int128(u), nil
}

func parseMagnitude(s string, base int) (u Uint128, neg, overflow, valid bool) {
	if base != 0 && (base < 2 || base > len(digits)) {
		panic("invalid base")
	}
	if len(s) == 0 {
		return
	}
	if s[0] == '+' || s[0] == '-' {
		neg = s[0] == '-'
		s = s[1:]
	}
	auto := base == 0
	previousDigit, haveDigit := false, false
	if auto {
		base = 10
		if len(s) > 0 && s[0] == '0' {
			base = 8
			previousDigit, haveDigit = true, true
			s = s[1:]
			if len(s) > 0 {
				explicitPrefix := true
				switch s[0] {
				case 'b', 'B':
					base = 2
				case 'o', 'O':
					base = 8
				case 'x', 'X':
					base = 16
				default:
					// The initial zero is an octal digit as well as a prefix.
					explicitPrefix = false
				}
				if explicitPrefix {
					s = s[1:]
					haveDigit = false
				}
			}
		}
	}
	for j := 0; j < len(s); j++ {
		ch := s[j]
		if ch == '_' && auto && previousDigit {
			previousDigit = false
			continue
		}
		var d int
		switch {
		case ch >= '0' && ch <= '9':
			d = int(ch - '0')
		case ch >= 'a' && ch <= 'z':
			d = int(ch-'a') + 10
		case ch >= 'A' && ch <= 'Z':
			d = int(ch-'A') + 10
			if base > 36 {
				d += 26
			}
		default:
			return
		}
		if d >= base {
			return
		}
		previousDigit, haveDigit = true, true
		if !overflow {
			carry, lo := bits.Mul64(u.lo, uint64(base))
			hiOverflow, hi := bits.Mul64(u.hi, uint64(base))
			hi, c1 := bits.Add64(hi, carry, 0)
			lo, c2 := bits.Add64(lo, uint64(d), 0)
			hi, c3 := bits.Add64(hi, 0, c2)
			u = Uint128{hi, lo}
			overflow = hiOverflow != 0 || c1 != 0 || c3 != 0
		}
	}
	valid = haveDigit && previousDigit
	return
}
