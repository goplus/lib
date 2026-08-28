//go:build !windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package os

const (
	F_GETFL = 3
	F_SETFL = 4

	O_RDONLY   = 0x0000
	O_WRONLY   = 0x0001
	O_RDWR     = 0x0002
	O_ACCMODE  = 0x0003
	O_NONBLOCK = 0x00000004
	O_CREAT    = 0x00000200
	O_TRUNC    = 0x00000400

	EAGAIN = 35
)
