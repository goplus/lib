# C networking bindings

This package binds `arpa/inet.h`, `net/if.h`, `netdb.h`, `netinet/in.h`, `sys/socket.h`, and `sys/un.h`, together with the C types those headers require. Each target is generated independently for this directory with unmodified llcppg commit `4724a68a99e7f8872ae0ba41bce2a5938acd7f9c`.

| Target | Header source | Build constraint |
| --- | --- | --- |
| macOS arm64 | macOS 14.0 SDK | `darwin && arm64 && !ios` |
| Linux amd64 | Bootlin stable 2017.05, glibc 2.24 | `linux && amd64 && !android` |
| Linux arm64 | Bootlin stable 2017.05, glibc 2.24 | `linux && arm64 && !android` |

Linux support assumes glibc. Existing implementations remain selected for other targets.

`legacy_native.go` retains the existing Unix function names and signatures, including `Bind`, `Accept`, `GetHostByName`, `SetSockOpt`, and `Freeaddrinfo`, along with the five untyped socket constants. Address types retain the existing field names for corresponding C fields. Linux address families use a 16-bit field without the BSD length byte, and Linux `AddrInfo` places `Addr` before `CanOnName`. `In6Addr.U6Addr` remains a byte array with native IPv6 alignment. `Iovec` aliases the shared `c.Iovec` definition.

`llcppg.cfg` selects macOS arm64; `llcppg-linux_amd64.cfg` and `llcppg-linux_arm64.cfg` record the Linux inputs. `Dir` and `Stdlib` refer to locally prepared headers under ignored `.sysroot` directories. `Deps` contains directories such as `..`. Local generation resolves these directories through the existing `Config.IncludeDir` callback before calling `Config.NewPackage`; the llcppg command itself is unchanged. Generated Go files receive target build constraints. `llcppg.pub` combines the available platform type mappings for dependent bindings.

The input configuration preserves existing handwritten functions and skips unsupported declarations, including the flexible array in Linux `cmsghdr` and the `if_nameindex` name collision. macOS declarations requiring unsupported dependency types are also omitted. This is a selection of common native APIs, rather than a complete libc binding.

Validation compares C and Go sizes, alignments, field offsets, and field widths using the original target headers. LLGo runtime checks exercise sockets, address resolution, IPv6 conversion, and interface names in both GC and nogc modes on all three targets. Linux runtime checks ran on newer glibc installations; the generation baseline is glibc 2.24.
