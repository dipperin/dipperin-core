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
	"testing"
	"math/big"

	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"fmt"
	"reflect"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"encoding/json"
)

type testAccountDB struct {
	balance map[common.Address]*big.Int
}

func (testDB *testAccountDB) GetBalance(addr common.Address) (*big.Int, error) {
	return testDB.balance[addr], nil
}
func (testDB *testAccountDB) AddBalance(addr common.Address, amount *big.Int) error {
	if _, ok := testDB.balance[addr]; !ok {
		testDB.balance[addr] = big.NewInt(0)
	}
	testDB.balance[addr].Add(testDB.balance[addr], amount)
	return nil
}

func (testDB *testAccountDB) SubBalance(addr common.Address, amount *big.Int) error {
	cmpResult := testDB.balance[addr].Cmp(amount)
	if cmpResult == -1 {
		return errors.New("the balance is not enough")
	}

	testDB.balance[addr].Sub(testDB.balance[addr], amount)

	return nil
}

var testDB = &testAccountDB{
	balance: make(map[common.Address]*big.Int, 0),
}

type testJson struct {
	BuiltInERC20Token
	Element1 *big.Int `json:"element_1"`
	element4 economy_model.Foundation
	Element2 []int    `json:"element_2"`
	Element3 []int64  `json:"element_3"`
}

func TestMap(t *testing.T) {
	testMap := make(map[string]int,0)
	testMap["a"] = 1
	testMap["a"] = 2

	log.Info("the testMap is:","testMap",testMap)
}

func TestReflect(t *testing.T) {
	value := reflect.ValueOf(&testJson{})

	fmt.Printf("the testJson method num is:%v\r\n", value.NumMethod())
}

func TestJsonEncode(t *testing.T) {
	test := testJson{
		BuiltInERC20Token: *newTestToken(),
		Element1:          big.NewInt(10),
		element4:          economy_model.MakeDipperinFoundation(economy_model.DIPProportion),
		Element2:          []int{1, 2, 3},
		Element3:          []int64{6, 7, 8},
	}

	encodeResult := util.StringifyJson(test)

	fmt.Printf("the encodeResult is:%v\r\n", encodeResult)
}

func TestEarlyTokenUnmarshalJSON(t *testing.T){
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	log.Info("the Token owner is:","owner",contract.Owner.Hex())
	value,ok := contract.Balances[contract.Owner.Hex()]
	assert.EqualValues(t,true,ok)

	log.Info("the value is:","value",value)
}

func TestUnmarshalData(t *testing.T){
	return
	data := `0x7b226572635f3230223a7b2262616c616e636573223a7b2230783030303065613442393738416435324435636245464164366634653364623136383865443642363539313337223a2230783237353134222c2230783030303065343437423842373835314433464244354336413033363235443238386366453942623565463045223a2230783237393466222c2230783030303036664337453942333964364330304137363741416441336530354145413762613864373145443644223a2230783236636437222c2230783030303033433432453364313644443539446235364537343164323930373735343433343439363838423230223a2230783237363331222c2230783030303034393234343745343038314437333532314236383334436536324545463446334635453862426237223a2230783237306439222c2230783030303046353562453637316538666632313834423043373138316165304534434439323432394330333443223a2230783236396239222c2230783030303036353332323535363630443965323238443939376463443832374465433638356239613137636131223a2230783237316264222c2230783030303037654465344435443830384441386132363732383462333845303041426363623432383839644632223a2230783236336236222c2230783030303032423941663833393043336361314461353730353464433839343763396639326131434130433939223a2230783236656438222c2230783030303034313739443537653435436233623534443646414546363965373436626632343045323837393738223a2230783236393934222c2230783030303044384642613764443934363534453364666266613738423736626234306442363035623445366341223a2230783236663131222c2230783030303032344644464265424431426464346144663542616337663034323031436631426265454463363837223a223078336332326563222c2230783030303062353738393845623830363439623246393939336438413339343145643139353936313336384539223a2230783236666635222c2230783030303044416132386543353263323834636138344161634244363930333932363964376133363234376332223a2230783237326461222c2230783030303035454343463041416136453846343531303738343438613138323937306538306362446432353362223a2230783236396239222c2230783030303030433642383744443033643830643232393033316461644332636431304664384138433234313333223a2230783236663833222c2230783030303063613836304237334136663045623365304143326233343239443137393136616565464236343938223a2230786466303563222c2230783030303046634439303339346539453930323232304266643741326430343433316233434642613265326444223a2230783237336637222c2230783030303030616464303461633464353237446538363663444534633933614531363632323134363137423132223a2230783237353836222c2230783030303035393636464645436339314632364246396236424232623745303941306430303436356532393430223a2230783237303637222c2230783030303034336235663165393830393242374437623932364566463536636342384439323841646461383842223a22307831316666303163222c2230783030303034324633643263323230333738434638396630463661333338356339323735333037363231414234223a22307831663333313136343663222c2230783030303039313863373733383830423436323932394143453446393735436366454439426532643845666339223a224561726c79546f6b656e222c2230783030303046393834373432423333304543393837433344463739433731634531653732393439384363363133223a2230783236363632222c2230783030303064353532433765633737333536363835373161386564393262353731323246323166436561353939223a2230783236623438222c2230783030303043324335364336363162363446364564353930353632383141303834443743444334353132413830223a2230783237363063227d2c226f776e6572223a2230783030303034326633643263323230333738636638396630663661333338356339323735333037363231616234222c22746f6b656e5f646563696d616c73223a332c22746f6b656e5f73796d626f6c223a224561726c79526577617264222c22746f6b656e5f746f74616c5f737570706c79223a22307831663334623066623030227d2c226368616e67655f746f5f63736b5f746f6b656e223a22307830222c2265786368616e67655f72617465223a5b3332365d2c226e6565645f63736b223a22307866383564646539356234383030227d`
	bytes,err := hexutil.Decode(data)
	assert.NoError(t,err)

	log.Info("the bytes string is:","bytes",string(bytes))
	var contract EarlyRewardContract
	err = json.Unmarshal(bytes, &contract)
	assert.NoError(t,err)

}

func TestMakeEarlyRewardContract(t *testing.T) {
	foundation := economy_model.MakeDipperinFoundation(economy_model.DIPProportion)

	owner := economy_model.EarlyTokenAddresses[0]
	decimalBase := big.NewInt(99999999999)
	initAmount := big.NewInt(99999999999999)
	testContract,err := MakeEarlyRewardContract(foundation, initAmount, economy_model.InitExchangeRate, tokenName, DecimalUnits, tokenSymbol, owner)
	assert.Error(t,err,errors.New("the DIP isn't enough"))

	decimalBase = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(DecimalUnits)), nil)
	initAmount = big.NewInt(0).Mul(economy_model.EarlyTokenAmount, decimalBase)

	testContract,err = MakeEarlyRewardContract(foundation, initAmount, economy_model.InitExchangeRate, tokenName, DecimalUnits, tokenSymbol, owner)
	assert.NoError(t,err)

	log.Info("the remainder DIP is:","remainder",big.NewInt(0).Sub(economy_model.EarlyTokenDIP,testContract.NeedDIP))

	log.Info("the initial exchangeRate is:", "initialExchangeRate", economy_model.InitExchangeRate)
	earlyContractStr := util.StringifyJson(testContract)

	fmt.Printf("the earlyContractStr is:%v\r\n", earlyContractStr)
}

func TestEarlyRewardContract_SetExchangeRate(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	initialDIP := contract.NeedDIP
	log.Info("the initialDIP is:","initialDIP",initialDIP)
	log.Info("the early Token account balance is:","balance",big.NewInt(0).Sub(economy_model.EarlyTokenDIP,initialDIP))
	log.Info("the need DIP is:", "NeedDIP", contract.NeedDIP)
	log.Info("the total supply is:", "supply", contract.TokenTotalSupply)
	log.Info("the initial exchange rate is:", "exchangeRate", contract.ExchangeRate[0])

	exchangeRate := int64(300)

	errT := contract.SetExchangeRate(common.HexToAddress("1234"), exchangeRate)
	assert.Error(t, errT, errors.New("the address isn't foundation maintenance address"))

	maintenanceAddress := economy_model.MaintenanceAddresses[0]
	contract.AccountDB = testDB

	errF := contract.SetExchangeRate(maintenanceAddress, contract.ExchangeRate[len(contract.ExchangeRate)-1])
	assert.Nil(t, errF)

	err := contract.SetExchangeRate(maintenanceAddress, exchangeRate)
	assert.NoError(t, err)

	log.Info("the tokenTotal is:", "EarlyTokenAmount", economy_model.EarlyTokenAmount)
	needValue := big.NewInt(0).Mul(economy_model.EarlyTokenAmount, big.NewInt(exchangeRate))
	needValue.Div(needValue, big.NewInt(economy_model.EarlyTokenExchangeBase))
	needValue.Mul(needValue, big.NewInt(consts.DIP))

	assert.EqualValues(t, needValue, contract.NeedDIP)

	returnValue := big.NewInt(0).Sub(initialDIP, needValue)

	contractCreatorAddressBalance, err := contract.AccountDB.GetBalance(contract.Owner)

	log.Info("the  returnValue is:","returnValue",returnValue)
	assert.NoError(t, err)
	assert.EqualValues(t, returnValue, contractCreatorAddressBalance)
	return
}

func TestEarlyRewardContract_TransferEDIPToDIP(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	errF := contract.TransferEDIPToDIP(economy_model.MaintenanceAddresses[0], (*hexutil.Big)(big.NewInt(0x21fc)))
	assert.Error(t, errF, errors.New("the address isn't NotFoundationAddress"))

	contract.AccountDB = testDB

	rewardAddress := common.HexToAddress("0x0000f01dA91C64eF6202c735e9362010196a556C7fc7")
	DIPReward := big.NewInt(1740000000)
	blockNumber := uint64(30)

	err := contract.RewardMineMaster(DIPReward, blockNumber, rewardAddress)
	assert.NoError(t, err)

	tokenValue := contract.BalanceOf(rewardAddress)
	assert.EqualValues(t, big.NewInt(8700), tokenValue)

	err = contract.TransferEDIPToDIP(rewardAddress, tokenValue)
	assert.NoError(t, err)

	tokenValue = contract.BalanceOf(rewardAddress)
	assert.EqualValues(t, tokenValue, big.NewInt(0))

	exchangeRate := contract.GetExchangeRate()
	decimal := contract.Decimals()
	base := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	DIP := big.NewInt(0).Mul(big.NewInt(8700), big.NewInt(exchangeRate))
	DIP.Mul(DIP, big.NewInt(consts.DIP))
	DIP.Div(DIP, base)
	DIP.Div(DIP, big.NewInt(economy_model.EarlyTokenExchangeBase))

	balance, err := contract.AccountDB.GetBalance(rewardAddress)
	assert.NoError(t, err)
	assert.EqualValues(t, DIP, balance)
}

func TestEarlyRewardContract_Destroy(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	assert.Error(t, contract.Destroy(common.HexToAddress("1234")), errors.New("can't destroy the contract"))
}

func TestEarlyRewardContract_Transfer(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	errF := contract.TransferEDIPToDIP(economy_model.MaintenanceAddresses[0], (*hexutil.Big)(big.NewInt(0x21fc)))
	assert.Error(t, errF, errors.New("the address should be normalAddress"))
}

func TestEarlyRewardContract_TransferFrom(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	ret := contract.TransferFrom(economy_model.MaintenanceAddresses[0], economy_model.MaintenanceAddresses[0], (*hexutil.Big)(big.NewInt(0x21fc)))
	assert.Equal(t, ret, false)

	ret = contract.TransferFrom(common.HexToAddress("1234"), economy_model.MaintenanceAddresses[0], (*hexutil.Big)(big.NewInt(0x21fc)))
	assert.Equal(t, ret, false)

}

func TestEarlyRewardContract_Approve(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	contract.CurSender = economy_model.MaintenanceAddresses[0]
	ret := contract.Approve(common.HexToAddress("5678"), (*hexutil.Big)(big.NewInt(0x21fc)))
	assert.Equal(t, ret, false)

	contract.CurSender = common.HexToAddress("1234")
	ret = contract.Approve(common.HexToAddress("5678"), (*hexutil.Big)(big.NewInt(0x21fc)))
	assert.Equal(t, ret, true)
}

func TestEarlyRewardContract_Allowance(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}
	ret := contract.Allowance(economy_model.MaintenanceAddresses[0], economy_model.MaintenanceAddresses[0])
	assert.Equal(t, ret, big.NewInt(0))

	ret = contract.Allowance(common.HexToAddress("0x1234"), economy_model.MaintenanceAddresses[0])
	assert.Equal(t, ret, big.NewInt(0))

	ret = contract.Allowance(common.HexToAddress("0x1234"), common.HexToAddress("0x5678"))
	assert.Equal(t, ret, big.NewInt(0))
}

func TestEarlyRewardContract_BalanceOf(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}
	decimal := contract.Decimals()
	assert.EqualValues(t, 3, decimal)
	base := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	ownerValue := contract.BalanceOf(contract.Owner).ToInt()
	expectValue := big.NewInt(0).Mul(economy_model.EarlyTokenAmount, base)
	assert.EqualValues(t, expectValue, ownerValue)

	log.Info("the expectValue is:", "expectValue", expectValue)
}

func TestEarlyRewardContract_RewardMineMaster(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	rewardAddress := common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	DIPReward := big.NewInt(1740000000)

	testReward := []*big.Int{
		big.NewInt(8700),
		big.NewInt(8004),
		big.NewInt(6612),
		big.NewInt(4524),
		big.NewInt(1740),
	}

	expectValue := big.NewInt(0)
	//year 1~5
	for i := 1; i <= 1; i++ {
		blockNumber := uint64(i) * economy_model.HeightAfterOneYear
		ownerValue := big.NewInt(0).Set(contract.BalanceOf(contract.Owner).ToInt())

		err := contract.RewardMineMaster(DIPReward, blockNumber, rewardAddress)
		assert.NoError(t, err)
		tokenValue := contract.BalanceOf(rewardAddress).ToInt()
		expectValue.Add(expectValue, testReward[i-1])
		ownerValueAfterReward := big.NewInt(0).Set(contract.BalanceOf(contract.Owner).ToInt())

		assert.EqualValues(t, testReward[i-1], big.NewInt(0).Sub(ownerValue, ownerValueAfterReward))
		assert.EqualValues(t, expectValue, tokenValue)
	}
}

func TestEarlyRewardContract_RewardVerifier(t *testing.T) {
	var contract EarlyRewardContract
	if err := util.ParseJson(EarlyRewardContractStr, &contract); err != nil {
		panic(err.Error())
	}

	DIPReward := map[economy_model.VerifierType]*big.Int{
		economy_model.MasterVerifier:    big.NewInt(3768115),
		economy_model.CommitVerifier:    big.NewInt(7536231),
		economy_model.NotCommitVerifier: big.NewInt(1884057),
	}

	addresses := map[economy_model.VerifierType][]common.Address{
		economy_model.MasterVerifier:    {common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")},
		economy_model.CommitVerifier:    {common.HexToAddress("0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978")},
		economy_model.NotCommitVerifier: {common.HexToAddress("0x00006fC7E9B39d6C00A767AAdA3e05AEA7ba8d71ED6D")},
	}

	blockNumber := uint64(30)
	err := contract.RewardVerifier(DIPReward, blockNumber, addresses)

	masterBalance := contract.BalanceOf(addresses[economy_model.MasterVerifier][0]).ToInt()
	assert.EqualValues(t, big.NewInt(18), masterBalance)

	commitBalance := contract.BalanceOf(addresses[economy_model.CommitVerifier][0]).ToInt()
	assert.EqualValues(t, big.NewInt(37), commitBalance)

	notCommitBalance := contract.BalanceOf(addresses[economy_model.NotCommitVerifier][0]).ToInt()
	assert.EqualValues(t, big.NewInt(9), notCommitBalance)

	assert.NoError(t, err)
}
