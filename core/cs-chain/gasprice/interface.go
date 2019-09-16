package gasprice

import "github.com/dipperin/dipperin-core/core/model"

type Chain interface {
	GetBlockByNumber(number uint64) model.AbstractBlock
	CurrentBlock() model.AbstractBlock
}
