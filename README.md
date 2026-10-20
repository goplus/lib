LLGo Standard Libraries
=====

[![Build Status](https://github.com/goplus/lib/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/lib/actions/workflows/go.yml)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/lib.svg?label=release)](https://github.com/goplus/lib/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/lib.svg)](https://pkg.go.dev/github.com/goplus/lib)
[![XGo](https://img.shields.io/badge/project-XGo-blue.svg)](https://github.com/goplus/xgo)
<!--
[![Coverage Status](https://codecov.io/gh/goplus/lib/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/lib)
-->

Native C bindings use the `C.` linkname prefix for both functions and variables.
This revision requires LLGo with
[C-prefixed variable support](https://github.com/xgo-dev/llgo/pull/2882);
older compilers do not recognize the prefix on variables.
