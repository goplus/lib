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

#include <cstdlib>
#include <iostream>
#include <string>

extern "C" {

char *llgoLLVMItaniumDemangle(const char *data, size_t size, int parseParams);
char *llgoLLVMMicrosoftDemangle(const char *data, size_t size, size_t *nRead,
                                int *status, int flags);
char *llgoLLVMRustDemangle(const char *data, size_t size);

} // extern "C"

static bool expectEqual(const char *name, char *value, const char *want) {
  if (value == nullptr) {
    std::cerr << name << " returned null\n";
    return false;
  }
  const std::string got(value);
  std::free(value);
  if (got != want) {
    std::cerr << name << " returned " << got << ", want " << want << '\n';
    return false;
  }
  return true;
}

int main() {
  constexpr char itaniumName[] = "_Z3foov";
  bool ok = expectEqual(
      "Itanium",
      llgoLLVMItaniumDemangle(itaniumName, sizeof(itaniumName) - 1, 1),
      "foo()");

  constexpr char microsoftName[] = "?foo@@YAHH@Z";
  size_t nRead = 0;
  int status = -1;
  char *microsoft = llgoLLVMMicrosoftDemangle(
      microsoftName, sizeof(microsoftName) - 1, &nRead, &status, 0);
  if (microsoft == nullptr) {
    std::cerr << "Microsoft returned null\n";
    ok = false;
  } else {
    const std::string got(microsoft);
    std::free(microsoft);
    if (status != 0 || nRead == 0 || got.find("foo") == std::string::npos) {
      std::cerr << "Microsoft returned " << got << ", nRead " << nRead
                << ", status " << status << '\n';
      ok = false;
    }
  }

  constexpr char rustName[] = "_RNvC6_123foo3bar";
  ok = expectEqual("Rust", llgoLLVMRustDemangle(rustName, sizeof(rustName) - 1),
                   "123foo::bar") &&
       ok;
  return ok ? EXIT_SUCCESS : EXIT_FAILURE;
}
