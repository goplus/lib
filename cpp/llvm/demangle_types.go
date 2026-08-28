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

package llvm

import "github.com/goplus/lib/c"

/*
	enum MSDemangleFlags {
	  MSDF_None = 0,
	  MSDF_DumpBackrefs = 1 << 0,
	  MSDF_NoAccessSpecifier = 1 << 1,
	  MSDF_NoCallingConvention = 1 << 2,
	  MSDF_NoReturnType = 1 << 3,
	  MSDF_NoMemberType = 1 << 4,
	  MSDF_NoVariableType = 1 << 5,
	};
*/
type MSDemangleFlags c.Int
