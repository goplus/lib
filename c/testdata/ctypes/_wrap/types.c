// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

#include <stddef.h>
#include <stdint.h>

#ifdef _WIN32
#include <basetsd.h>
typedef SSIZE_T lib_ctypes_ssize_t;
#else
#include <sys/types.h>
typedef ssize_t lib_ctypes_ssize_t;
#endif

unsigned int lib_ctypes_size_bytes(void) { return sizeof(size_t); }
unsigned int lib_ctypes_long_bytes(void) { return sizeof(long); }
unsigned int lib_ctypes_wchar_bytes(void) { return sizeof(wchar_t); }
unsigned long long lib_ctypes_char_minus_one(void) { return (char)-1; }
unsigned long long lib_ctypes_wchar_minus_one(void) { return (wchar_t)-1; }
size_t lib_ctypes_echo_size(size_t value) { return value; }
lib_ctypes_ssize_t lib_ctypes_echo_signed(lib_ctypes_ssize_t value) { return value; }
long lib_ctypes_echo_long(long value) { return value; }
ptrdiff_t lib_ctypes_echo_ptrdiff(ptrdiff_t value) { return value; }
intptr_t lib_ctypes_echo_intptr(intptr_t value) { return value; }
uintptr_t lib_ctypes_echo_uintptr(uintptr_t value) { return value; }
wchar_t lib_ctypes_echo_wchar(wchar_t value) { return value; }

int lib_ctypes_check_words(size_t *words, size_t stride) {
  if (stride != sizeof(size_t)) return 0;
  size_t value = 7;
  if (sizeof(size_t) == 8) value += (uint64_t)1 << 40;
  if (words[0] != value || words[1] != value + 2) return 0;
  words[1] = value + 4;
  return 1;
}
