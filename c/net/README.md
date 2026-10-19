# C networking bindings

This package contains native bindings for `arpa/inet.h`, `net/if.h`, `netdb.h`, `netinet/in.h`, `sys/socket.h`, and `sys/un.h`, with the types those headers require.

The generated Unix bindings are extracted from the [libc bindings](https://github.com/MeteorsLiu/libc/commit/7385771b1366f8ba181672ea0eef9f4951d6b284), generated with llcppg commit `0204f71b7441956a16511c0389dae8c8a2234060`:

| Target | Header source | Build constraint |
| --- | --- | --- |
| macOS arm64 | macOS 14.0 SDK | `darwin && arm64 && !ios` |
| Linux amd64 | Bootlin glibc 2.44 | `linux && amd64 && !android` |

Linux support assumes glibc. The existing Windows implementation and the existing fallback for other targets remain selected by their build constraints.

`legacy_native.go` retains the existing Unix function names and signatures, including `Bind`, `Accept`, `GetHostByName`, `SetSockOpt`, and `Freeaddrinfo`. It also retains the five untyped socket constants so existing `Socket(..., SOCK_STREAM, ...)` calls compile. The generated address types use the existing field names for corresponding C fields; their layout and constants come from each platform's headers. Linux address families use a 16-bit field without the BSD length byte, and Linux `AddrInfo` places `Addr` before `CanOnName`. `In6Addr.U6Addr` remains a byte array with the native IPv6 union alignment.

`Iovec` is an alias of the shared `c.Iovec` definition, retaining its existing fields and their platform-specific types.

`llcppg.pub` lists the available network type mappings for dependent bindings. Availability follows the target platform; the manifest combines the platform mappings.

The extraction keeps networking declarations and their required type declarations, omits duplicate C functions already covered by the retained API, and applies the platform build constraints.

Known generation omissions include the inherited function skips, macOS packed structures whose generated alignment differs from C, and the Linux `cmsghdr` flexible array and its dependent helper. Opaque forward declarations remain pointer-only types.

Validation compares C and Go sizes, alignments, and key field offsets, and exercises sockets, address resolution, IP address conversion, and interface names with LLGo on both generated targets.
