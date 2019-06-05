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


package mineworker

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/health-info-log"
)

func NewDefaultWorkExecutor(work *minemsg.DefaultWork, submitter workSubmitter) *defaultWorkExecutor {
	ex := &defaultWorkExecutor{
		curWork:   work,
		submitter: submitter,
		//nonceSuffix:big.NewInt(0),
	}
	return ex

}

type workSubmitter interface {
	SubmitWork(work minemsg.Work)
}

type defaultWorkExecutor struct {
	curWork   *minemsg.DefaultWork
	submitter workSubmitter

	// The first 8 bytes are allocated by the fragment server and the miner server, and cannot be changed.
	// The latter field is freely played by the miners.

	nonceSuffix [common.NonceLength - 8]byte
	//nonceSuffix *big.Int
}

//
func (executor *defaultWorkExecutor) ChangeNonce() bool {
	//todo use byte operation to optimal this nonce change step
	// Nonce increments and verify the validity
	for index := len(executor.nonceSuffix) - 1; index >= 0; {
		if executor.nonceSuffix[index] < 255 {
			executor.nonceSuffix[index]++
			break
		} else {
			executor.nonceSuffix[index] = 0
			index--
		}
	}
	//log.Info("change nonce", "suffix", executor.nonceSuffix)
	//put result to header nonce
	copy(executor.curWork.BlockHeader.Nonce[8:], executor.nonceSuffix[:])

	//
	//executor.nonceSuffix.Add(executor.nonceSuffix,big.NewInt(1))
	//copy(executor.curWork.BlockHeader.Nonce[8:], executor.nonceSuffix.Bytes())

	//log.Debug("defaultWorkExecutor ChangeNonce", "executor.curWork.BlockHeader.Hash()", executor.curWork.BlockHeader.Hash().Hex())
	// check nonce valid
	//if executor.curWork.BlockHeader.Hash().ValidHashForDifficulty(executor.curWork.BlockHeader.Diff) {
	//new hash calculation replacing upstairs
	bHash, err := executor.curWork.CalHash()
	if err == nil {
		if bHash.ValidHashForDifficulty(executor.curWork.BlockHeader.Diff) {
			// some thing interesting here
			executor.curWork.ResultNonce = executor.curWork.BlockHeader.Nonce
			log.Info("ChangeNonce successful")
			//fmt.Println(executor.curWork.BlockHeader.String())
			health_info_log.Info("found nonce", "height", executor.curWork.BlockHeader.Number)
			return true
		}
	} else {
		log.Info("search nonce", "error", err)
	}
	return false
}

func (executor *defaultWorkExecutor) Submit() {
	executor.submitter.SubmitWork(executor.curWork)
}
