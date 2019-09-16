package ver_halt_check

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"testing"
)

func TestEvalulate(t *testing.T) {
	index, proof := crypto.Evaluate(testVerBootAccounts[0].Pk, common.HexToHash("0x0000003").Bytes())
	index2, proof2 := crypto.Evaluate(testVerBootAccounts[0].Pk, common.HexToHash("0x0000003").Bytes())
	fmt.Println(index, index2, proof, proof2)
	//assert.Equal(t, proof2, proof)
}
