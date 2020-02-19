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
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func f1(c *BlockContext) Middleware {
	return func() error {
		println(" inside f1")
		return c.Next()
	}
}

func f3(c *BlockContext) Middleware {
	return func() error {
		println(" inside f3")
		return nil
	}
}

func f4(c *BlockContext) Middleware {
	return func() error {
		println(" inside f4")
		return c.Next()
	}
}

func f5(c *BlockContext) Middleware {
	return func() error {
		println(" inside f5")
		return c.Next()
	}
}

func f6(c *BlockContext) Middleware {
	return func() error {
		println(" inside f6")
		return c.Next()
	}
}

func f2() error {
	println("inside f2")
	return nil
}

func TestProcWithMiddleware(t *testing.T) {
	bc := NewBlockContext(nil, nil)
	assert.Equal(t, true, len(bc.middlewares) == 0)
	bc.Use(f3(bc))
	bc.Use(f2, f4(bc))
	bc.Use(f1(bc))
	assert.Equal(t, true, len(bc.middlewares) == 4)
	//m.ProcessBlock()
	bc.Process(f5(bc), f6(bc))
}

func CheckBlock(c *BlockContext) Middleware {
	return func() error {
		println(" inside checkBlock")
		return errors.New("test error")
	}
}

func TestFailedFirstOne(t *testing.T) {
	bc := NewBlockContext(nil, nil)
	//assert.Equal(t, true, len(bc.middlewares) == 0)
	assert.Len(t, bc.middlewares, 0)
	bc.Use(CheckBlock(bc))
	bc.Use(ValidateBlockNumber(bc))
	bc.Use(UpdateStateRoot(bc))
	bc.Use(UpdateBlockVerifier(bc))
	bc.Use(InsertBlock(bc))
	//assert.Equal(t, true, len(bc.middlewares) == 5)
	assert.Len(t, bc.middlewares, 5)
	err := bc.Process()
	assert.Error(t, err)
}

func TestMiddlewareChain_Last(t *testing.T) {
	var c MiddlewareChain
	assert.Nil(t, c.Last())

	c = []Middleware{func() error {
		return nil
	}}
	assert.NotNil(t, c.Last())

	mc := &MiddlewareContext{
		middlewares: c,
	}
	assert.NotNil(t, mc.Middleware())
}
