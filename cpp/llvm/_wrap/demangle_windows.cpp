/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#if defined(_WIN32) && defined(__clang__) && __clang_major__ < 20
#define _ALLOW_COMPILER_AND_STL_VERSION_MISMATCH
#endif

#include <llvm/Demangle/Demangle.h>

#include <cstddef>
#include <string_view>

static std::string_view llgoLLVMStringView(const char *data, size_t size) {
    return size == 0 ? std::string_view() : std::string_view(data, size);
}

extern "C" {

char *llgoLLVMItaniumDemangle(const char *data, size_t size, int parseParams) {
    return llvm::itaniumDemangle(llgoLLVMStringView(data, size), parseParams != 0);
}

char *llgoLLVMMicrosoftDemangle(const char *data, size_t size, size_t *nRead,
                                int *status, int flags) {
    return llvm::microsoftDemangle(
        llgoLLVMStringView(data, size), nRead, status,
        static_cast<llvm::MSDemangleFlags>(flags));
}

char *llgoLLVMRustDemangle(const char *data, size_t size) {
    return llvm::rustDemangle(llgoLLVMStringView(data, size));
}

} // extern "C"
