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

import (
	"github.com/dipperin/dipperin-core/core/model"
)

func NewBlockContext(block model.AbstractBlock, chain ChainInterface) *BlockContext {
	return &BlockContext{
		MiddlewareContext: MiddlewareContext{ index: -1 },
		Block: block,
		Chain: chain,
	}
}

/*
visit chain, db, state_root through processor
 */
type BlockContext struct {
	MiddlewareContext

	// the block to be handled
	Block model.AbstractBlock
	// chain
	Chain ChainInterface
}

// basic middleware, can be comprised by other middleware
type MiddlewareContext struct {
	// index of middleware, initial value=-1
	index int8
	middlewares MiddlewareChain
}

/*
the core function of middleware,
called only once after the registration of each middleware
 */
func (mc *MiddlewareContext) Next() error {
	mc.index++
	// this loop in middleware can go to the end even if next is not called
	for mc.index < int8(len(mc.middlewares)) {
		if err := mc.middlewares[mc.index](); err != nil {
			return err
		}
		mc.index++
	}
	return nil
}

func (mc *MiddlewareContext) Middleware() Middleware{
	return mc.middlewares.Last()
}

func (mc *MiddlewareContext) Use(m ...Middleware) {
	mc.middlewares = append(mc.middlewares, m...)
}

func (mc *MiddlewareContext) Process(m ...Middleware) error {
	mc.middlewares = append(mc.middlewares, m...)
	return mc.Next()
}

type Middleware func() error
type MiddlewareChain []Middleware

func (c MiddlewareChain) Last() Middleware{
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
