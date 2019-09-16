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

package commands

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestGetUnit(t *testing.T) {
	value, unit := GetUnit("1.23456DIP")
	assert.Equal(t, "1.23456", value)
	assert.Equal(t, "DIP", unit)

	value, unit = GetUnit("123456WU")
	assert.Equal(t, "123456", value)
	assert.Equal(t, "WU", unit)

	value, unit = GetUnit("0.9999")
	assert.Equal(t, "0.9999", value)
	assert.Equal(t, "WU", unit)
}

func TestMoneyValueToCSCoin(t *testing.T) {
	value, err := MoneyValueToCSCoin("0.001")
	assert.Error(t, err)
	assert.Nil(t, value)

	value, err = MoneyValueToCSCoin("S.001DIP")
	assert.Error(t, err)
	assert.Nil(t, value)

	value, err = MoneyValueToCSCoin("0.0000000000000000001DIP")
	assert.Error(t, err)
	assert.Nil(t, value)

	value, err = MoneyValueToCSCoin("10000000000000000000WU")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(10), big.NewInt(consts.DIP)), value)

	value, err = MoneyValueToCSCoin("1.23456DIP")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1.23456*consts.DIP), value)

	value, err = MoneyValueToCSCoin("1.23456UDIP")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1.23456*consts.UDIP), value)

	value, err = MoneyValueToCSCoin("1.23456MDIP")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1.23456*consts.MDIP), value)

	value, err = MoneyValueToCSCoin("1.23456GWU")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1.23456*consts.GDIPUNIT), value)

	value, err = MoneyValueToCSCoin("1.23456MWU")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1.23456*consts.MDIPUNIT), value)

	value, err = MoneyValueToCSCoin("1.234KWU")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1.234*consts.KDIPUNIT), value)

	value, err = MoneyValueToCSCoin("1234WU")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1234*consts.DIPUNIT), value)
}

func TestCSCoinToMoneyValue(t *testing.T) {
	csCoinValue1 := (*hexutil.Big)(big.NewInt(0).Mul(big.NewInt(10000000), big.NewInt(consts.GDIPUNIT)))

	csCoinValue2 := (*hexutil.Big)(big.NewInt(0).Mul(big.NewInt(34545000), big.NewInt(consts.GDIPUNIT)))

	csCoinValue3 := (*hexutil.Big)(big.NewInt(0).Mul(big.NewInt(600000000), big.NewInt(consts.GDIPUNIT)))

	csCoinValue4 := (*hexutil.Big)(big.NewInt(0).Mul(big.NewInt(897878600000000), big.NewInt(consts.GDIPUNIT)))

	csCoinValue5 := (*hexutil.Big)(big.NewInt(0).Mul(big.NewInt(3069000000000), big.NewInt(consts.GDIPUNIT)))

	moneyValue1, err := CSCoinToMoneyValue(csCoinValue1)
	assert.NoError(t, err)
	assert.Equal(t, "0.01DIP", moneyValue1)

	moneyValue2, err := CSCoinToMoneyValue(csCoinValue2)
	assert.NoError(t, err)
	assert.Equal(t, "0.034545DIP", moneyValue2)

	moneyValue3, err := CSCoinToMoneyValue(csCoinValue3)
	assert.NoError(t, err)
	assert.Equal(t, "0.6DIP", moneyValue3)

	moneyValue4, err := CSCoinToMoneyValue(csCoinValue4)
	assert.NoError(t, err)
	assert.Equal(t, "897878.6DIP", moneyValue4)

	moneyValue5, err := CSCoinToMoneyValue(csCoinValue5)
	assert.NoError(t, err)
	assert.Equal(t, "3069DIP", moneyValue5)
}

func TestMoneyWithUnit(t *testing.T) {
	assert.Equal(t, "100WU", MoneyWithUnit("100"))
	assert.Equal(t, "100dip", MoneyWithUnit("100dip"))
	assert.Equal(t, "100aaa", MoneyWithUnit("100aaa"))
	assert.Equal(t, "100.00WU", MoneyWithUnit("100.00"))
}

func TestParseWalletPathAndName(t *testing.T) {
	gotPath, gotName := ParseWalletPathAndName("A.B")
	assert.Equal(t, "A.B", gotPath)
	assert.Equal(t, "A.B", gotName)

	gotPath, gotName = ParseWalletPathAndName("A/B")
	assert.Equal(t, "A/B", gotPath)
	assert.Equal(t, "B", gotName)
}

func TestCheckAndChangeHexToAddress(t *testing.T) {
	addr, err := CheckAndChangeHexToAddress("1234")
	assert.Equal(t, g_error.ErrInvalidAddressLen, err)
	assert.Equal(t, common.Address{}, addr)

	addr, err = CheckAndChangeHexToAddress("0000005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	assert.Equal(t, g_error.ErrInvalidAddressPrefix, err)
	assert.Equal(t, common.Address{}, addr)

	addr, err = CheckAndChangeHexToAddress("0x12345586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	assert.Equal(t, g_error.UnknownTxTypeErr, err)
	assert.Equal(t, common.Address{}, addr)

	addr, err = CheckAndChangeHexToAddress(from)
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress(from), addr)
}

func TestDecimalToInter(t *testing.T) {
	moneyValue1 := "0.001"
	moneyValue2 := "7.89"

	moneyValue3 := "a.2234"

	moneyValue4 := "300"

	moneyValue5 := "0.0001"
	unitBit := 3

	unit := 1
	for i := 0; i < unitBit; i++ {
		unit *= 10
	}
	assert.EqualValues(t, 1000, unit)

	value, err := DecimalToInter(moneyValue1, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1), value)

	value, err = DecimalToInter(moneyValue2, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(7890), value)

	value, err = DecimalToInter(moneyValue3, unitBit)
	assert.Error(t, err)

	value, err = DecimalToInter(moneyValue4, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(300000), value)

	value, err = DecimalToInter(moneyValue5, unitBit)
	assert.Error(t, err)
}

func TestInterToDecimal(t *testing.T) {
	csCoinValue1 := (*hexutil.Big)(big.NewInt(100))
	csCoinValue2 := (*hexutil.Big)(big.NewInt(34545))
	csCoinValue3 := (*hexutil.Big)(big.NewInt(6000))
	csCoinValue4 := (*hexutil.Big)(big.NewInt(8978786))
	csCoinValue5 := (*hexutil.Big)(big.NewInt(30690000))
	csCoinValue6 := (*hexutil.Big)(big.NewInt(23400))
	unitBit := 3

	moneyValue1, err := InterToDecimal(csCoinValue1, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "0.1", moneyValue1)

	moneyValue2, err := InterToDecimal(csCoinValue2, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "34.545", moneyValue2)

	moneyValue3, err := InterToDecimal(csCoinValue3, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "6", moneyValue3)

	moneyValue4, err := InterToDecimal(csCoinValue4, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "8978.786", moneyValue4)

	moneyValue5, err := InterToDecimal(csCoinValue5, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "30690", moneyValue5)

	moneyValue6, err := InterToDecimal(csCoinValue6, 0)
	assert.NoError(t, err)
	assert.Equal(t, "23400", moneyValue6)
}

func TestCheckRegistration(t *testing.T) {
	assert.Equal(t, false, CheckRegistration())
}
