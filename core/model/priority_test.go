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

package model

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCalPriority(t *testing.T) {
	reputation := uint64(5014)
	luck := common.HexToHash("3848465468464823742384877578453654634536842542465141654353536656")
	priority, err := CalPriority(luck, reputation)
	assert.NoError(t, err)
	assert.NotNil(t, priority)
}

func TestCalReputation(t *testing.T) {
	x := int64(5000)
	stake := big.NewInt(x)
	nonce := uint64(0)
	performance := uint64(0)
	reputation, err := CalReputation(nonce, stake, performance)
	assert.NoError(t, err)
	assert.NotNil(t, reputation)

	reputation, err = CalReputation(nonce, big.NewInt(0), performance)
	assert.Equal(t, "stake not sufficient", err.Error())
	assert.Equal(t, uint64(0), reputation)
}

type Value struct {
	name  string
	value float64
}

/*1000 txs,the effect of nonce decreases rapidly*/
func TestElemNonce(t *testing.T) {
	//epsilon1 := 1.0
	//epsilon2 := 0.1
	//epsilon3 := 50.0
	epsilon1 := 100.0
	epsilon2 := 0.008
	epsilon3 := 500.0

	nonceArray := []Value{
		{name: "0", value: 0.0},
		{name: "10", value: 10.0},
		{name: "50", value: 50.0},
		{name: "80", value: 80.0},
		{name: "100", value: 100.0},
		{name: "200", value: 200.0},
		{name: "300", value: 300.0},
		{name: "500", value: 500.0},
		{name: "600", value: 600.0},
		{name: "700", value: 700.0},
		{name: "800", value: 800.0},
		{name: "900", value: 900.0},
		{name: "1000", value: 1000.0},
		{name: "2000", value: 2000.0},
		//{name:"3000" ,value: 3000.0},
		//{name:"4000" ,value: 4000.0},
		//{name:"5000" ,value: 5000.0},
		//{name:"6000" ,value: 6000.0},
		//{name:"7000" ,value: 7000.0},
		//{name:"8000" ,value: 8000.0},
		//{name:"9000" ,value: 9000.0},
		{name: "10000", value: 10000.0},
		//{name:"20000" ,value: 20000.0},
	}

	for key := range nonceArray {
		fmt.Print(nonceArray[key].name, "   ")
		fmt.Println(Elem(nonceArray[key].value, epsilon1, epsilon2, epsilon3))
	}

}

/*9000 stakes,the effect of stake decreases rapidly*/
func TestElemStake(t *testing.T) {
	//epsilon1 := 1.0
	//epsilon2 := 0.1
	//epsilon3 := 50.0
	epsilon1 := 10000.0
	epsilon2 := 0.001
	epsilon3 := 1000.0

	stakeArray := []Value{
		{name: "0", value: 0.0},
		{name: "10", value: 10.0},
		{name: "50", value: 50.0},
		{name: "80", value: 80.0},
		{name: "100", value: 100.0},
		{name: "200", value: 200.0},
		{name: "300", value: 300.0},
		{name: "500", value: 500.0},
		{name: "600", value: 600.0},
		{name: "700", value: 700.0},
		{name: "800", value: 800.0},
		{name: "900", value: 900.0},
		{name: "1000", value: 1000.0},
		{name: "2000", value: 2000.0},
		{name: "3000", value: 3000.0},
		{name: "4000", value: 4000.0},
		{name: "5000", value: 5000.0},
		{name: "6000", value: 6000.0},
		{name: "7000", value: 7000.0},
		{name: "8000", value: 8000.0},
		{name: "9000", value: 9000.0},
		{name: "10000", value: 10000.0},
		{name: "20000", value: 20000.0},
	}

	for key := range stakeArray {
		fmt.Print(stakeArray[key].name, "   ")
		fmt.Println(Elem(stakeArray[key].value, epsilon1, epsilon2, epsilon3))
	}

}

/*elections reached about 70, the effect of performance decreases rapidly*/
func TestElemPerformance(t *testing.T) {

	epsilon1 := 1.0
	epsilon2 := 5.0
	epsilon3 := 0.1

	performanceArray := []Value{
		{name: "0", value: 0.0},
		{name: "0", value: 0.0},
		{name: "10", value: 10.0},
		{name: "50", value: 50.0},
		{name: "80", value: 80.0},
		{name: "100", value: 100.0},
		{name: "200", value: 200.0},
		{name: "300", value: 300.0},
		{name: "500", value: 500.0},
		{name: "600", value: 600.0},
		{name: "700", value: 700.0},
		{name: "800", value: 800.0},
		{name: "900", value: 900.0},
		{name: "1000", value: 1000.0},
		{name: "2000", value: 2000.0},
		{name: "3000", value: 3000.0},
		{name: "4000", value: 4000.0},
		{name: "5000", value: 5000.0},
		{name: "6000", value: 6000.0},
		{name: "7000", value: 7000.0},
		{name: "8000", value: 8000.0},
		{name: "9000", value: 9000.0},
		{name: "10000", value: 10000.0},
	}

	for key := range performanceArray {
		fmt.Print(performanceArray[key].name, "   ")
		fmt.Println(Elem(performanceArray[key].value/10000.0, epsilon1, epsilon2, epsilon3))
	}

}

func TestTestCalculator_GetElectPriority(t *testing.T) {
	var ttt TestCalculator
	hash := common.HexToHash("3848465468464823742384878845453654634536842542465141654353539999")
	x := int64(500000)
	stake := big.NewInt(x)
	nonce := uint64(70)
	performance := uint64(50)

	fmt.Println(ttt.GetElectPriority(hash, nonce, stake, performance))

}

const (
	BaseMulti       = 3
	NonceMax        = 10000
	PerformanceMax  = 100
	StakeMax        = 500
	BasePerformance = 30
	BaseNonce       = 100
	MaxPerformance  = 100
	BasePenalty     = 10
	BaseReword      = 1
)

type TestAccount struct {
	accountAddress common.Address `json:"-"`
	Name           string
	nonce          uint64
	stake          *big.Int
	performance    uint64
	count          uint64
	priority       uint64
}

/*
按每8秒出一个算：
   	一轮110个块  花费14.6分钟
    一天出块10800块  将近100轮
	按performance=30,stake=100,nonce=100计算
		100个竞选验证者  运行1000轮 大约10天  大约有66个验证者performance值大于90
		1000个竞选验证者  运行1000轮 大约10天  保持stake，nonce不变，大约有75个验证者performance值大约90

	结论：performance可能增长过快，luck值影响有点小

Imagine one block is generated every 8 seconds:
	110 blocks a round will take approximately 14.6 minutes
	in one day 10800 new blocks or 100 rounds will be generated
	suppose performance=30,stake=100,nonce=100
    	100 candidates, 1000 rounds, 10days, there will be 66 verifiers whose performances are greater than 90
		1000 candidates, 1000 rounds, 10days, with the same stake and nonce, there will be 75 verifiers whose performances are greater than 90
	conclusion: performance will grow very fast, check the influence of the luck value
*/
func TestCalculator_GetElectPriority_Multi(t *testing.T) {
	//var tc TestCalculator
	// accounts num
	accountNum := uint(1000)
	// round num of elect
	round := uint(1000)
	// wanted verifier num
	topVerifierNum := uint(22)
	resultMap := make(map[uint][]*TestAccount)
	//accounts := CreateMultiAccountAddress(accountNum)
	// generate rand nonce,stake
	//accounts,stakeCount := CreateMultiAccount(accountNum, GenerateRand,NonceMax,generateNormalPerformance,BasePerformance,generateRandStake,StakeMax)
	// generate rand stake
	//accounts,stakeCount := CreateMultiAccount("account", accountNum, generateNormalNonce,BaseNonce,generateNormalPerformance,BasePerformance,generateRandStake,StakeMax)
	// generate normal nonce,performance,stake
	accounts, stakeCount := CreateMultiAccount("account", accountNum, generateNormalNonce, BaseNonce, generateNormalPerformance, BasePerformance, generateNormalStake, StakeMax)
	log.DLogger.Info("TestCalculator_GetElectPriority_Multi stakeCount info", zap.Int64("stakeCount", stakeCount.Int64()))
	baseHashStr := "3848465468464823742384878845453654634536842542465141654353539999"
	for j := uint(0); j < round; j++ {
		hash := create256Seed(baseHashStr + strconv.Itoa(int(j)))
		var jAccounts []*TestAccount
		//log.DLogger.Info("TestCalculator_GetElectPriority_Multi round info", "round", j)
		for i := 0; i < len(accounts); i++ {
			luck := getLuck(accounts[i].accountAddress, hash)
			if accounts[i].performance > MaxPerformance {
				accounts[i].performance = MaxPerformance
			}
			//priority, err := tc.GetElectPriority(luck, accounts[i].nonce, accounts[i].stake, accounts[i].performance)
			reputation, err := CalReputation(accounts[i].nonce, accounts[i].stake, accounts[i].performance)
			if err != nil {
				log.DLogger.Error("TestCalculator_GetElectPriority_Multi calReputation", zap.Uint64("reputation", reputation))
				panic("calReputation err")
			}
			priority, err := CalPriority(luck, reputation)
			if err != nil {
				log.DLogger.Error("TestCalculator_GetElectPriority_Multi  priority err", zap.Error(err))
				panic("priority is wrong")
			}
			//log.DLogger.Info("priofity info", "priority", priority)
			accounts[i].priority = priority
			jAccounts = GetSortAccount(jAccounts, accounts[i])
		}
		//log.DLogger.Info("jAccounts length", "jAccounts", len(jAccounts))
		//PrintAccountInfo(jAccounts)
		//random verifier-node's performance
		// top topVerifierNum nodes,add once elections
		jAccounts = jAccounts[:topVerifierNum]
		random := GenerateRand(topVerifierNum) + uint64(len(jAccounts)*2/3)
		for k := 0; k < len(jAccounts); k++ {
			jAccounts[k].count += 1
			//log.DLogger.Info("TestCalculator_GetElectPriority_Multi jAccounts info topVerifier", "jAccounts"+strconv.Itoa(k), jAccounts[k])
			if k < len(jAccounts)*2/3 {
				jAccounts[k].performance += BaseReword
			} else if k == int(random) {
				jAccounts[k].performance -= BasePenalty
			}
		}
		resultMap[j] = jAccounts
		//baseHashStr = GenerateNewStr(baseHashStr+strconv.Itoa(int(j)))
	}
	//PrintAccountElectResult(resultMap, round, stakeCount, topVerifierNum)
}

func generateNormalPerformance(performance uint) uint64 {
	return BasePerformance
}

func generateNormalNonce(nonce uint) uint64 {
	return BaseNonce
}

func generateRandStake(stake uint) *big.Int {
	return big.NewInt(0).Add(big.NewInt(int64(GenerateRand(StakeMax))), big.NewInt(100))
}

func generateNormalStake(stake uint) *big.Int {
	return big.NewInt(100)
}

// generate luck
func getLuck(addr common.Address, seed common.Hash) common.Hash {
	list := append(seed.Bytes(), addr.Bytes()...)
	return common.RlpHashKeccak256(list)
}

func PrintAccountInfo(tAccounts []*TestAccount) {
	for i := range tAccounts {
		log.DLogger.Info("PrintAccountInfo", zap.Any("account", tAccounts[i]))
	}
}

// Create multiAccountAddress,random nonce
func CreateMultiAccountAddress(count uint) (ta []*TestAccount) {
	ta = make([]*TestAccount, count)
	var testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	privLen := len(testPriv1)
	for i := uint(0); i < count; i++ {
		account := &TestAccount{}
		testPriv1 = GenerateNewStr(testPriv1 + strconv.Itoa(int(i)))
		testPriv1 = testPriv1[len(testPriv1)-privLen:]
		//log.DLogger.Info("CreateMultiAccountAddress  alicePriv info", "alicePriv", alicePriv)
		pk, err := crypto.HexToECDSA(testPriv1)
		if err != nil {
		}
		//signer := NewSigner(big.NewInt(int64((i+1) * BaseMulti)))
		addr := cs_crypto.GetNormalAddress(pk.PublicKey)
		account.accountAddress = addr
		account.Name = "account" + strconv.Itoa(int(i))
		account.nonce = GenerateRand(NonceMax)
		account.performance = BasePerformance
		account.stake = big.NewInt(0).Add(big.NewInt(int64(GenerateRand(StakeMax))), big.NewInt(100))
		ta[i] = account
		//log.DLogger.Info("CreateMultiAccountAddress  accounts info", "accounts", fmt.Sprintf("%v", ta[i]), "stake", ta[i].stake)
	}
	return
}

func CreateMultiAccount(baseName string, count uint, generateNonce func(nonceMax uint) uint64, nonceMax uint, generatePerformance func(performanceMax uint) uint64, performanceMax uint, generateStake func(stakeMax uint) *big.Int, stakeMax uint) (ta []*TestAccount, stakeCount *big.Int) {
	ta = make([]*TestAccount, count)
	var testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	stakeCount = big.NewInt(0)
	privLen := len(testPriv1)
	for i := uint(0); i < count; i++ {
		account := &TestAccount{}
		testPriv1 = GenerateNewStr(testPriv1 + strconv.Itoa(int(i)))
		testPriv1 = testPriv1[len(testPriv1)-privLen:]
		//log.DLogger.Info("CreateMultiAccountAddress  alicePriv info", "alicePriv", alicePriv)
		pk, err := crypto.HexToECDSA(testPriv1)
		if err != nil {
		}
		//signer := NewSigner(big.NewInt(int64((i+1) * BaseMulti)))
		addr := cs_crypto.GetNormalAddress(pk.PublicKey)
		account.accountAddress = addr
		account.Name = baseName + strconv.Itoa(int(i))
		account.nonce = generateNonce(nonceMax)
		account.performance = generatePerformance(performanceMax)
		//account.stake  = big.NewInt(0).Add(big.NewInt(int64(GenerateRand(StakeMax))), big.NewInt(100))
		account.stake = generateStake(stakeMax)
		stakeCount.Add(stakeCount, account.stake)
		ta[i] = account
		//log.DLogger.Info("CreateMultiAccountAddress  accounts info", "accounts", fmt.Sprintf("%v", ta[i]), "stake", ta[i].stake)
	}
	return
}

func GenerateNewStr(baseStr string) string {
	rand.Seed(int64(time.Now().Nanosecond()))
	randomStr := strconv.Itoa(int(rand.Int63n(2 ^ 30)))
	//fmt.Println(randomStr)
	baseStr = baseStr[0:10] + randomStr + baseStr[len(randomStr)+10:]
	return baseStr
}

// generate seed
func create256Seed(baseStr string) common.Hash {
	//t := time.Now()
	//h := sha3.New256()
	h := crypto.Keccak256([]byte(baseStr))
	//fmt.Println(h)
	return common.BytesToHash(h)
}

func TestCreatePasswd(t *testing.T) {
	baseHashStr := "3848465468464823742384878845453654634536842542465141654353539999"
	for i := 0; i < 10; i++ {
		fmt.Println(i, "       ", create256Seed(baseHashStr+strconv.Itoa(i)))
	}
}

func PrintAccountElectResult(accountsMap map[uint][]*TestAccount, round uint, stakeCount *big.Int, topVerifierNum uint) {
	resultMap := make(map[string]*TestAccount)
	for i := uint(0); i < round; i++ {
		accounts := accountsMap[i]
		for j := 0; j < len(accounts); j++ {
			//log.DLogger.Info("PrintAccountElectResult", "round", i, "account.accountAddress", accounts[j].Name, "count", accounts[j].count)
			resultMap[accounts[j].Name] = accounts[j]
		}
	}
	var performanceCount int
	for _, value := range resultMap {
		//log.DLogger.Info("PrintAccountElectResult", "key", key, "count", value.count, "count percentage", float64(value.count)/float64(round)/float64(topVerifierNum), "stake", value.stake, "stake percentage", float64(value.stake.Int64())/float64(stakeCount.Int64()), "nonce", value.nonce, "performance", value.performance)
		if value.performance > 90 {
			performanceCount++
		}
	}
	//log.DLogger.Info("PrintAccountElectResult", "performanceCount > 90  ", performanceCount)

}

// Sort all accounts  TODO  reason
func GetSortAccount(tAccounts []*TestAccount, account *TestAccount) []*TestAccount {
	if len(tAccounts) > 0 {
		isMinest := true
		for i := 0; i < len(tAccounts); i++ {
			if tAccounts[i].priority < account.priority {
				bakAccounts := make([]*TestAccount, len(tAccounts))
				copy(bakAccounts, tAccounts)
				bak := bakAccounts[i]
				front := append(append(bakAccounts[:i], account), bak)
				//log.DLogger.Info("GetSortAccount", "front",front,"bak",bak)
				//PrintAccountInfo(front)
				tAccounts = append(front, tAccounts[i+1:]...)
				//log.DLogger.Info("GetSortAccount", "i", i, "front", front)
				//PrintAccountInfo(tAccounts)
				isMinest = false
				break
			}
		}
		if isMinest {
			tAccounts = append(tAccounts, account)
		}
	} else {
		tAccounts = append(tAccounts, account)
	}
	//log.DLogger.Info("GetSortAccount", "tAccounts", tAccounts, "account", account)
	return tAccounts
}

func GenerateRand(max uint) uint64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return uint64(rand.Int63n(int64(max)))
}

func TestTestCalculator_GetReputation(t *testing.T) {
	tc := TestCalculator{}
	result, err := tc.GetReputation(1, big.NewInt(100), 85)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
