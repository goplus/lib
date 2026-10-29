#include <ffi.h>

void *llgo_ffi_closure_alloc(void **code) {
    return ffi_closure_alloc(sizeof(ffi_closure), code);
}
