// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

#include <stddef.h>
#include <stdint.h>
#include <sys/types.h>

unsigned int lib_ctypes_size_bytes(void) { return sizeof(size_t); }
size_t lib_ctypes_echo_size(size_t value) { return value; }
ssize_t lib_ctypes_echo_signed(ssize_t value) { return value; }
long lib_ctypes_echo_long(long value) { return value; }
ptrdiff_t lib_ctypes_echo_ptrdiff(ptrdiff_t value) { return value; }
intptr_t lib_ctypes_echo_intptr(intptr_t value) { return value; }
uintptr_t lib_ctypes_echo_uintptr(uintptr_t value) { return value; }
wchar_t lib_ctypes_echo_wchar(wchar_t value) { return value; }

int lib_ctypes_check_words(size_t *words) {
  size_t value = 7;
#if __SIZEOF_SIZE_T__ == 8
  value += (size_t)1 << 40;
#endif
  if (words[0] != value || words[1] != value + 2) return 0;
  words[1] = value + 4;
  return 1;
}
