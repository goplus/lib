# C ABI test helpers

`c/ctypes_abi_test.go` contains the shared `TestCTypeABI` test. These helpers
compile the same C source through cgo or LLGo. Native CI runs the test through
its existing Go and LLGo commands:

```sh
go test -v -run TestCTypeABI ./c
llgo test -v -run TestCTypeABI ./c
```

The test uses C's actual type widths, including Windows LLP64 (`size_t` is
64-bit while `long` is 32-bit), and checks `char`/`wchar_t` signedness and
`size_t` array storage. Its C source uses Windows `SSIZE_T` where POSIX
`ssize_t` is unavailable.

The following commands run all four Wasm provider/profile paths with a current
LLGo main. Set `LLGO_ROOT` to the matching LLGo checkout or installation; the
existing native compatibility CI uses LLGo v1.0.8 instead.

```sh
ctypes_tmp=$(mktemp -d)

# J32 / GoJS
GOOS=js GOARCH=wasm llgo test -O0 -v -run TestCTypeABI ./c

# J32 / Emscripten
llgo test -O0 -c -target emscripten -o "$ctypes_tmp/j32.mjs" ./c
node "$LLGO_ROOT/targets/emscripten-runner.mjs" "$ctypes_tmp/j32.mjs" -test.run=TestCTypeABI -test.v

# J64 / Emscripten Memory64
llgo test -O0 -c -target emscripten-memory64 -o "$ctypes_tmp/j64.mjs" ./c
node "$LLGO_ROOT/targets/emscripten-memory64-runner.mjs" "$ctypes_tmp/j64.mjs" -test.run=TestCTypeABI -test.v

# W32 / WASI
llgo test -O0 -c -target wasi -o "$ctypes_tmp/w32.wasm" ./c
wasmer run --enable-exceptions --enable-simd --stack-size=1048576 \
  --volume="$PWD" --volume="$ctypes_tmp" --env="PWD=$PWD" \
  "$ctypes_tmp/w32.wasm" -- -test.run=TestCTypeABI -test.v
```

Each successful run prints `--- PASS: TestCTypeABI` and `PASS`.
