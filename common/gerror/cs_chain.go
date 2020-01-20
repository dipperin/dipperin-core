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

package gerror

import "errors"

var (
	/*Cached chain state errors*/
	ErrTargetOutOfRange    = errors.New("the target special block is out of range")
	ErrPreTargetBlockIsNil = errors.New("pre target block is nil")

	/*Cached chain service errors*/
	ErrAlreadyHaveThisBlock = errors.New("already have this block")
	ErrNoGenesis            = errors.New("genesis not found in chain")
	ErrLastNumIsNil         = errors.New("last number is nil")

	/*Chain state errors*/
	ErrBlockNotFound = errors.New("block not found")
)
