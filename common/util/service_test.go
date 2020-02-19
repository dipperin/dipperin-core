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

package util

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testService struct {
	BaseService
}

func (testService) OnReset() error {
	return nil
}

func TestBaseServiceWait(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	ts.Start()

	waitFinished := make(chan struct{})
	go func() {
		ts.Wait()
		waitFinished <- struct{}{}
	}()

	go ts.Stop()

	select {
	case <-waitFinished:
		// all good
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected Wait() to finish within 100 ms.")
	}
}

func TestBaseServiceReset(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	ts.Start()

	err := ts.Reset()
	require.Error(t, err, "expected cant reset service error")

	ts.Stop()

	err = ts.Reset()
	require.NoError(t, err)

	err = ts.Start()
	require.NoError(t, err)
}

func TestBaseService_MulStart(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)

	ts.Start()
	assert.True(t, ts.IsRunning())
	ts.Stop()
	assert.False(t, ts.IsRunning())
	ts.Reset()

	ts.Start()
	assert.True(t, ts.IsRunning())
	ts.Stop()
	assert.False(t, ts.IsRunning())
}

func TestBaseService_SetLogger(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	var logger log.Logger
	log.New(&logger)
	ts.SetLogger(logger)

	assert.Equal(t, ts.Logger, logger)
}

func TestBaseService_Start(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	err := ts.Start()
	assert.NoError(t, err)

	assert.Equal(t, ts.IsRunning(), true)

	err = ts.Start()

	assert.Error(t, err)

	ts.Stop()

	err = ts.Start()

	assert.Error(t, err)
}

func TestBaseService_OnStart(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	assert.NoError(t, ts.OnStart(), nil)
}

func TestBaseService_Stop(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	err := ts.Start()
	assert.NoError(t, err)

	assert.Equal(t, ts.IsRunning(), true)

	ts.Stop()

	assert.Equal(t, ts.IsRunning(), false)

	ts.Stop()

	ts2 := &testService{}
	ts2.BaseService = *NewBaseService(nil, "TestService", ts2)

	ts2.Stop()
}

func TestBaseService_Reset(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)

	err := ts.Reset()
	assert.Error(t, err)

	err = ts.Start()

	assert.NoError(t, err)

	ts.Stop()

	err = ts.Reset()

	assert.NoError(t, err)

}

type errService struct{}

func (e *errService) Start() error {
	return nil
}

func (e *errService) OnStart() error {
	return errors.New("err start")
}

func (e *errService) Stop() {
	panic("implement me")
}

func (e *errService) OnStop() {
	panic("implement me")
}

func (e *errService) Reset() error {
	panic("implement me")
}

func (e *errService) OnReset() error {
	panic("implement me")
}

func (e *errService) IsRunning() bool {
	panic("implement me")
}

func (e *errService) Quit() <-chan struct{} {
	panic("implement me")
}

func (e *errService) String() string {
	return "errService"
}

func (e *errService) SetLogger(logger log.Logger) {
	panic("implement me")
}

func TestStart2(t *testing.T) {
	bs := NewBaseService(nil, "tests", &testService{})
	bs.stopped = 1
	assert.Error(t, bs.Start())
	assert.Panics(t, func() {
		bs.OnReset()
	})

	bs = NewBaseService(nil, "tests", &errService{})
	assert.Error(t, bs.Start())
	bs.Quit()
}

func TestNewBaseService(t *testing.T) {
	type args struct {
		logger log.Logger
		name   string
		impl   Service
	}
	tests := []struct {
		name string
		args args
		want *BaseService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseService(tt.args.logger, tt.args.name, tt.args.impl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseService_OnStop(t *testing.T) {
	tests := []struct {
		name string
		bs   *BaseService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bs.OnStop()
		})
	}
}

func TestBaseService_OnReset(t *testing.T) {
	tests := []struct {
		name    string
		bs      *BaseService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bs.OnReset(); (err != nil) != tt.wantErr {
				t.Errorf("BaseService.OnReset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseService_IsRunning(t *testing.T) {
	tests := []struct {
		name string
		bs   *BaseService
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bs.IsRunning(); got != tt.want {
				t.Errorf("BaseService.IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseService_Wait(t *testing.T) {
	tests := []struct {
		name string
		bs   *BaseService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bs.Wait()
		})
	}
}

func TestBaseService_String(t *testing.T) {
	tests := []struct {
		name string
		bs   *BaseService
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bs.String(); got != tt.want {
				t.Errorf("BaseService.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseService_Quit(t *testing.T) {
	tests := []struct {
		name string
		bs   *BaseService
		want <-chan struct{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bs.Quit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseService.Quit() = %v, want %v", got, tt.want)
			}
		})
	}
}
