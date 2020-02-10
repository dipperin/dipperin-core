package mineworker

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newRemoteConnector(t *testing.T) {
	rc := newRemoteConnector()
	assert.NotNil(t, rc)
}

func TestRemoteConnector_SetMineMasterPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// mock
	peerMock := NewMockPmAbstractPeer(ctrl)
	peerMock.EXPECT().SendMsg(gomock.Any(),gomock.Any()).Return(nil).AnyTimes()
	workerMock := NewMockWorker(ctrl)
	workerMock.EXPECT().CurrentCoinbaseAddress().Return(common.HexToAddress("123")).AnyTimes()
	workerMock.EXPECT().Start().Return().AnyTimes()
	// init
	rc := newRemoteConnector()
	assert.Nil(t, rc.peer)
	rc.worker = workerMock
	// test
	rc.SetMineMasterPeer(peerMock)
	assert.NotNil(t, rc.peer)
}

func TestRemoteConnector_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// test case
	situations := []struct {
		name      string
		given     func() *RemoteConnector
		expectErr bool
	}{
		{
			"occur error",
			func() *RemoteConnector {
				// mock
				peerMock1 := NewMockPmAbstractPeer(ctrl)
				peerMock1.EXPECT().SendMsg(gomock.Any(),gomock.Any()).Return(errors.New("occur error")).AnyTimes()
				workerMock1 := NewMockWorker(ctrl)
				workerMock1.EXPECT().CurrentCoinbaseAddress().Return(common.HexToAddress("123")).AnyTimes()
				// init
				rc := newRemoteConnector()
				rc.peer = peerMock1
				rc.worker = workerMock1
				// return
				return rc
			},
			true,
		},
		{
			"normal case",
			func() *RemoteConnector {
				// mock
				peerMock2 := NewMockPmAbstractPeer(ctrl)
				peerMock2.EXPECT().SendMsg(gomock.Any(),gomock.Any()).Return(nil).AnyTimes()
				workerMock2 := NewMockWorker(ctrl)
				workerMock2.EXPECT().CurrentCoinbaseAddress().Return(common.HexToAddress("123")).AnyTimes()
				// init
				rc := newRemoteConnector()
				rc.peer = peerMock2
				rc.worker = workerMock2
				// return
				return rc
			},
			false,
		},
	}
	// test
	for _, situation := range situations {
		executor := situation.given()
		err := executor.Register()
		if situation.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestRemoteConnector_UnRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// mock
	peerMock := NewMockPmAbstractPeer(ctrl)
	peerMock.EXPECT().SendMsg(gomock.Any(),gomock.Any()).Return(nil).AnyTimes()
	// init
	rc := newRemoteConnector()
	rc.peer = peerMock
	// test
	rc.UnRegister()
}

func TestRemoteConnector_SendMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// test case
	situations := []struct {
		name      string
		given     func() *RemoteConnector
		expectErr bool
	}{
		{
			"occur error",
			func() *RemoteConnector {
				// mock
				peerMock1 := NewMockPmAbstractPeer(ctrl)
				peerMock1.EXPECT().SendMsg(gomock.Any(),gomock.Any()).Return(errors.New("occur error")).AnyTimes()
				workerMock1 := NewMockWorker(ctrl)
				workerMock1.EXPECT().CurrentCoinbaseAddress().Return(common.HexToAddress("123")).AnyTimes()
				// init
				rc := newRemoteConnector()
				rc.peer = peerMock1
				rc.worker = workerMock1
				// return
				return rc
			},
			true,
		},
		{
			"normal case",
			func() *RemoteConnector {
				// mock
				peerMock2 := NewMockPmAbstractPeer(ctrl)
				peerMock2.EXPECT().SendMsg(gomock.Any(),gomock.Any()).Return(nil).AnyTimes()
				workerMock2 := NewMockWorker(ctrl)
				workerMock2.EXPECT().CurrentCoinbaseAddress().Return(common.HexToAddress("123")).AnyTimes()
				// init
				rc := newRemoteConnector()
				rc.peer = peerMock2
				rc.worker = workerMock2
				// return
				return rc
			},
			false,
		},
	}
	// test
	for _, situation := range situations {
		executor := situation.given()
		err := executor.SendMsg(1, 1)
		if situation.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
