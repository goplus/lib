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
var deferHeadTLS = tls.Alloc[*runtime.Defer](func(head **runtime.Defer) {
    if head != nil {
        *head = nil
    }
})

func pushDefer(head *runtime.Defer) {
    deferHeadTLS.Set(head)
}

func popDefer() *runtime.Defer {
    return deferHeadTLS.Get()
}
```

In GC builds the slot is automatically registered as a root, which prevents
Boehm from prematurely reclaiming the defer nodes. In builds that use
`-tags nogc`, the helper still uses pthread TLS but simply skips the root
registration.
