// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package middleware

const (
	bodyCacheLimit      = 256
	blockCacheLimit     = 256
	maxFutureBlocks     = 256
	headerCacheLimit    = 512
	numberCacheLimit    = 2048
	verifierCacheLimit  = 12
	maxTimeFutureBlocks = 30
	txSizeMax           = 100000000
	txSizeMin           = 0
	txAmountMax         = 100000000
	txAmountMin         = 0

	// test env have not enough mem
	bodyCacheLimitTestEnv      = 30
	blockCacheLimitTestEnv     = 30
	maxFutureBlocksTestEnv     = 30
	headerCacheLimitTestEnv    = 60
	numberCacheLimitTestEnv    = 512
	maxTimeFutureBlocksTestEnv = 15
)
