# C networking bindings

This package binds `arpa/inet.h`, `net/if.h`, `netdb.h`, `netinet/in.h`, `sys/socket.h`, and `sys/un.h`, together with the C types those headers require. Each target is generated as one common libc translation unit, then classified into the existing C packages, with unmodified llcppg commit `4724a68a99e7f8872ae0ba41bce2a5938acd7f9c`.

| Target | Header source | Build constraint |
| --- | --- | --- |
| macOS arm64 | macOS 14.0 SDK | `darwin && arm64 && !ios` |
| Linux amd64 | Bootlin stable 2017.05, glibc 2.24 | `linux && amd64 && !android` |
| Linux arm64 | Bootlin stable 2017.05, glibc 2.24 | `linux && arm64 && !android` |

Linux support assumes glibc. Existing implementations remain selected for other targets.

`legacy_native.go` retains the existing Unix function names and signatures, including `Bind`, `Accept`, `GetHostByName`, `SetSockOpt`, and `Freeaddrinfo`, along with the five untyped socket constants. Address types retain the existing field names for corresponding C fields. Linux address families use a 16-bit field without the BSD length byte, and Linux `AddrInfo` places `Addr` before `CanOnName`. `In6Addr.U6Addr` remains a byte array with native IPv6 alignment. `Iovec` aliases the shared `c.Iovec` definition.

`llcppg.cfg` and its Linux variants retain local header ownership and dependency metadata for downstream bindings. The current snapshots come from unified generation followed by classification, rather than running these directory configurations independently. `Dir` and `Stdlib` refer to locally prepared headers under ignored `.sysroot` directories. `Deps` contains directories such as `..`; the local caller resolves them through the existing `Config.IncludeDir` callback. Each target is kept in one `libc_<os>_<arch>.go` file. Generation inputs, scripts, and header caches are kept locally and are not included in these package changes. `llcppg.pub` lists only types retained in the available target outputs.

The input preserves existing handwritten functions and skips unsupported declarations. `cmsghdr` is excluded on every target. Records requiring packing, explicit alignment, bit fields, flexible arrays, unsupported scalar types, or unresolved callback-field layouts are excluded together with declarations that depend on them. The existing `if_nameindex` name collision is also skipped. This is a selection of common native APIs, rather than a complete libc binding.

Validation compares C and Go sizes, alignments, field offsets, and field widths using the original target headers. LLGo runtime checks exercise sockets, address resolution, IPv6 conversion, and interface names in both GC and nogc modes on all three targets. Linux runtime checks ran on newer glibc installations; the generation baseline is glibc 2.24.
