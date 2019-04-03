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


package address_util

import (
	"encoding/hex"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/dipperin/dipperin-core/common"
)

func TestGenERC20Address(t *testing.T) {
	_, err := GenERC20Address()
	assert.NoError(t, err)
}

func TestGenContractAddress(t *testing.T) {
	_, err := GenContractAddress(1)
	assert.NoError(t, err)
}

func TestPubKeyToAddress(t *testing.T) {

	type args struct {
		p           string
		addressType int
	}
	tests := []struct {
		name string
		args args
		want common.Address
	}{
		{
			name: "normal",
			args: args{
				p: "04760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1b01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d",
				addressType: 1,
			},
			want: common.HexToAddress("0x00015891906FeF64a5AE924c7FC5ED48C0f64a55fCE1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkh, _ := hex.DecodeString(tt.args.p)
			pk, _ := crypto.UnmarshalPubkey(pkh)
			if got := PubKeyToAddress(*pk, tt.args.addressType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubKeyToAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
