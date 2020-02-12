package mineworker

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeLocalWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	masterServer := NewMockMasterServer(ctrl)
	worker := MakeLocalWorker(common.HexToAddress("123"), 1, masterServer)
	assert.NotNil(t, worker)
}

func TestMakeRemoteWorker(t *testing.T) {
	worker, remoteConnector := MakeRemoteWorker(common.HexToAddress("123"), 1)
	assert.NotNil(t, worker)
	assert.NotNil(t, remoteConnector)
}