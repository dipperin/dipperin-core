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

package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
)

// codec struct
type paramsCodec struct {
	decMu  sync.Mutex                // guards the decoder
	decode func(v interface{}) error // decoder to allow multiple transports
}

// create codec
func NewParamsCodec(r io.Reader) *paramsCodec {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	return &paramsCodec{
		decode: dec.Decode,
	}
}

// ParseRequestArguments tries to parse the given params (json.RawMessage) with the given
// types. It returns the parsed values or an error when the parsing failed.
func (c *paramsCodec) ParseRequestArguments(argTypes []reflect.Type) ([]reflect.Value, error) {
	var incomingMsg json.RawMessage

	if err := c.decode(&incomingMsg); err != nil {
		return nil, err
	}
	return parsePositionalArguments(incomingMsg, argTypes)
}

// parsePositionalArguments tries to parse the given args to an array of values with the
// given types. It returns the parsed values or an error when the args could not be
// parsed. Missing optional arguments are returned as reflect.Zero values.
func parsePositionalArguments(rawArgs json.RawMessage, types []reflect.Type) ([]reflect.Value, error) {
	// Read beginning of the args array.
	dec := json.NewDecoder(bytes.NewReader(rawArgs))
	if tok, _ := dec.Token(); tok != json.Delim('[') {
		return nil, errors.New("non-array args")
	}
	// Read args.
	args := make([]reflect.Value, 0, len(types))
	for i := 0; dec.More(); i++ {
		if i >= len(types) {
			return nil, errors.New(fmt.Sprintf("too many arguments, want at most %d", len(types)))
		}
		argval := reflect.New(types[i])
		if err := dec.Decode(argval.Interface()); err != nil {
			return nil, errors.New(fmt.Sprintf("invalid argument %d: %v", i, err))
		}
		if argval.IsNil() && types[i].Kind() != reflect.Ptr {
			return nil, errors.New(fmt.Sprintf("missing value for required argument %d", i))
		}
		args = append(args, argval.Elem())
	}
	// Read end of args array.
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	// Set any missing args to nil.
	for i := len(args); i < len(types); i++ {
		if types[i].Kind() != reflect.Ptr {
			return nil, errors.New(fmt.Sprintf("missing value for required argument %d", i))
		}
		args = append(args, reflect.Zero(types[i]))
	}
	return args, nil
}
