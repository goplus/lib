# GC-safe Thread Local Storage

The `c/tls` package implements a small helper around pthread TLS that keeps the
slot visible to the Boehm garbage collector when it is enabled. It is useful for
runtime features (such as defer stacks) that keep GC-managed pointers inside the
pthread-specific area, which Boehm does not scan by default.

## Key ideas

- The package exposes a generic `tls.Handle[T]` type. `tls.Alloc` creates a
  handle and registers the slot with Boehm (`GC_add_roots`) in GC builds, or
  falls back to a plain `pthread_setspecific` slot in `-tags nogc` builds.
- A slot stores a value of type `T` directly, so even non-pointer data (for
  example counters) can be stored without an additional allocation.
- An optional destructor can be supplied when allocating the handle. It is
  invoked whenever the owning thread exits, allowing callers to release
  resources or reset global state.

## Example

```go
package main

import (
    "log"

    "github.com/goplus/lib/c/tls"
)

type stats struct {
    handled int
}

var threadStats = tls.Alloc[*stats](func(s **stats) {
    // Destructor runs when the thread exits.
    if *s != nil {
        log.Printf("thread handled %d requests", (*s).handled)
    }
})

func recordRequest() {
    cur := threadStats.Get()
    if cur == nil {
        cur = &stats{}
    }
    cur.handled++
    threadStats.Set(cur)
}

func resetStats() {
    threadStats.Clear()
}
```

In GC builds the slot is automatically registered as a root, which prevents
Boehm from reclaiming any GC-managed memory reachable from the `stats`
instances. In builds that use `-tags nogc`, the helper still relies on pthread
TLS but simply skips the root registration step.
