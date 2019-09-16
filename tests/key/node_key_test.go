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

package key

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPubKey(t *testing.T) {
	//pk, _ := crypto.HexToECDSA("6008df48c4c385961e8e0a56681f6ac2600075905bfed416ee27b3f531a81888")
	pk, _ := crypto.HexToECDSA("6008df48c4c385961e8e0a56681f6ac2600075905bfed416ee27b3f531a81888")
	n := enode.NewV4(&pk.PublicKey, net.ParseIP("127.0.0.1"), 10000, 10000)
	fmt.Println(n.String())

	sig1, err := crypto.Sign(common.HexToHash("0x231231").Bytes(), pk)
	assert.NoError(t, err)
	sig2, err := crypto.Sign(common.HexToHash("0x231231").Bytes(), pk)
	assert.NoError(t, err)
	assert.Equal(t, hexutil.Encode(sig1), hexutil.Encode(sig2))
}

func TestReplaceIPAddr(t *testing.T) {
	str := "enode://8b610c5400bfdb355c7a204beb65cb261fe5e89cc2c1837dc3cf752d16df65cfd95c6ee79be3720ef9dc6ba0b6876c63530a6352cb18298afb2b282b111ec7cf@[::]:40006"
	remote := "192.168.122.102:40006"

	i1 := strings.Index(str, "@")
	i2 := strings.LastIndex(str, ":")
	i3 := strings.LastIndex(remote, ":")

	str = str[:i1+1] + remote[:i3] + str[i2:]
	assert.Equal(t, "enode://8b610c5400bfdb355c7a204beb65cb261fe5e89cc2c1837dc3cf752d16df65cfd95c6ee79be3720ef9dc6ba0b6876c63530a6352cb18298afb2b282b111ec7cf@192.168.122.102:40006", str)
}

func TestDivBig(t *testing.T) {
	pub, _ := big.NewInt(0).SetString("e731657120688fd98488e97ea855333be46ae5ffe4a56a8a7285beae7243debf454c4419ddee0444806de7adba026daed02a022f996d215c2e5c17f750784de6", 16)
	priv, _ := big.NewInt(0).SetString("9828f4c2e2844c7530ea5f41729186d6034cd08d0e72ec5614588d570f0305f9", 16)

	fmt.Println(pub)
	fmt.Println(pub.Div(pub, priv))
}

func TestFilepath(t *testing.T) {
	fmt.Println(filepath.Dir("sdf/zz/x/y"))
	fmt.Println(filepath.Base("sdf/zz/x/y"))
}
