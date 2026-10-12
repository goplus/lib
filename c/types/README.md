# Lightweight 128-bit integers

`c/types` provides signed and unsigned 128-bit integers for `c.Int128T` and
`c.Uint128T`. Its only direct dependency is `math/bits`. Importing the core
package does not pull in `fmt`, `math/big`, `reflect`, or `os`, and performs no
big-integer initialization.

Arithmetic, comparison, bit operations, ordinary integer conversions, and the
existing constructors keep their names and representations. `String`, `Text`,
`ParseInt128`, `ParseUint128`, `FormatInt128`, and `FormatUint128` operate directly
on two `uint64` words.

Text conversion supports bases 2 through 62 for every value (including small
unsigned values, which previously delegated to `strconv` and rejected bases
above 36). Parsing preserves `math/big.Int.SetString` syntax: base 0 recognizes
binary, octal, and hexadecimal prefixes and valid underscore separators;
explicit bases do not accept separators. Overflow returns the nearest endpoint
and an error. Invalid syntax returns zero and an error. Invalid bases panic.
Error messages no longer depend on `fmt`; callers should check the error rather
than match its text.

## Migrating optional APIs

`c/types/extra` imports `fmt` and `math/big` only when explicitly requested.
The following APIs have moved out of the core package:

| Previous API | Replacement |
| --- | --- |
| `types.Int128_Cast__8(b)` / `Int128_Cast__9(b)` | `extra.Int128FromBig(b)` returns `(value, inRange)` |
| `types.Uint128_Cast__8(b)` / `Uint128_Cast__9(b)` | `extra.Uint128FromBig(b)` returns `(value, inRange)` |
| `i.BigInt()` | `extra.Int128ToBig(new(big.Int), i)` |
| `u.BigInt()` | `extra.Uint128ToBig(new(big.Int), u)` |
| `i.ToBigInt(dst)` / `u.ToBigInt(dst)` | `extra.Int128ToBig(dst, i)` / `extra.Uint128ToBig(dst, u)` |
| Custom `fmt` formatting/scanning on core values | `extra.Int128{Int128: i}` / `extra.Uint128{Uint128: u}` |

The wrappers embed the core value, provide `BigInt`, `ToBigInt`, `Format`, and
`Scan`, and promote core methods. Arithmetic returns core values; wrap the result
again when custom formatting is needed. A scan updates the embedded field only
on success. Scanning retains the original decimal-token behavior. Conversions
from `big.Int` clamp out-of-range values and do not modify their input;
conversions to `big.Int` fully replace the destination, including its old sign.

```go
import (
    "fmt"
    "math/big"

    "github.com/goplus/lib/c/types"
    "github.com/goplus/lib/c/types/extra"
)

value, err := types.ParseInt128("-123456789012345678901234567890", 10)
if err != nil {
    panic(err)
}
fmt.Printf("%#x\n", extra.Int128{Int128: value})
b := extra.Int128ToBig(new(big.Int), value)
restored, inRange := extra.Int128FromBig(b)
_, _ = restored, inRange

var scanned extra.Int128
_, err = fmt.Sscan("42", &scanned)
_, _ = scanned.Int128, err
```

This is an intentional API migration. Core values still implement `String`, but
no longer implement `fmt.Formatter` or `fmt.Scanner`. Callers using numeric
format verbs or scanning must use the wrappers. Platform build tags do not
change the public API.

## Validation

Run `go test ./c/types/...`. Tests compare parsing and formatting against
`math/big`, cover signed/unsigned limits, destination reuse, optional adapters,
and check the core's transitive dependency graph. `math/big` is used only as a
test oracle and by the optional extension package.

Run `go test ./c/types -run '^$' -fuzz '^FuzzParse$' -fuzztime=30s` for additional
parser differential testing. LLGo integration validation uses a local module
replacement and its complete `test/buildcache/test.sh` suite, including the
ESP32-C3 target that exposed the original dependency regression.
