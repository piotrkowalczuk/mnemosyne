package mnemosynetest

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

import "golang.org/x/net/context"

type Mnemosyne struct {
	mock.Mock
}

// FromContext provides a mock function with given fields: _a0
func (_m *Mnemosyne) FromContext(_a0 context.Context) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(context.Context) *mnemosynerpc.Session); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *Mnemosyne) Get(_a0 context.Context, _a1 mnemosynerpc.AccessToken) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(context.Context, mnemosynerpc.AccessToken) *mnemosynerpc.Session); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosynerpc.AccessToken) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: _a0, _a1
func (_m *Mnemosyne) Exists(_a0 context.Context, _a1 mnemosynerpc.AccessToken) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, mnemosynerpc.AccessToken) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosynerpc.AccessToken) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Mnemosyne) Start(_a0 context.Context, _a1 string, _a2 string, _a3 map[string]string) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]string) *mnemosynerpc.Session); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, map[string]string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: _a0, _a1
func (_m *Mnemosyne) Abandon(_a0 context.Context, _a1 mnemosynerpc.AccessToken) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, mnemosynerpc.AccessToken) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetValue provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Mnemosyne) SetValue(_a0 context.Context, _a1 mnemosynerpc.AccessToken, _a2 string, _a3 string) (map[string]string, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, mnemosynerpc.AccessToken, string, string) map[string]string); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosynerpc.AccessToken, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type Storage struct {
	mock.Mock
}

// Setup provides a mock function with given fields:
func (_m *Storage) Setup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TearDown provides a mock function with given fields:
func (_m *Storage) TearDown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0, _a1
func (_m *Storage) Start(_a0 string, _a1 map[string]string) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(string, map[string]string) *mnemosynerpc.Session); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: _a0
func (_m *Storage) Abandon(_a0 *mnemosynerpc.AccessToken) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*mnemosynerpc.AccessToken) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosynerpc.AccessToken) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *Storage) Get(_a0 *mnemosynerpc.AccessToken) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(*mnemosynerpc.AccessToken) *mnemosynerpc.Session); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosynerpc.AccessToken) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Storage) List(_a0 int64, _a1 int64, _a2 *time.Time, _a3 *time.Time) ([]*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []*mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(int64, int64, *time.Time, *time.Time) []*mnemosynerpc.Session); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64, *time.Time, *time.Time) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: _a0
func (_m *Storage) Exists(_a0 *mnemosynerpc.AccessToken) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*mnemosynerpc.AccessToken) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosynerpc.AccessToken) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1, _a2
func (_m *Storage) Delete(_a0 *mnemosynerpc.AccessToken, _a1 *time.Time, _a2 *time.Time) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*mnemosynerpc.AccessToken, *time.Time, *time.Time) int64); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosynerpc.AccessToken, *time.Time, *time.Time) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: _a0, _a1, _a2
func (_m *Storage) SetValue(_a0 *mnemosynerpc.AccessToken, _a1 string, _a2 string) (map[string]string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(*mnemosynerpc.AccessToken, string, string) map[string]string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosynerpc.AccessToken, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type suite struct {
	mock.Mock
}

// setup provides a mock function with given fields: _a0
func (_m *suite) setup(_a0 testing.T) {
	_m.Called(_a0)
}

// teardown provides a mock function with given fields: _a0
func (_m *suite) teardown(_a0 testing.T) {
	_m.Called(_a0)
}

type RPCClient struct {
	mock.Mock
}

// Context provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Context(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*mnemosynerpc.ContextResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.ContextResponse
	if rf, ok := ret.Get(0).(func(context.Context, *empty.Empty, ...grpc.CallOption) *mnemosynerpc.ContextResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ContextResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *empty.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Get(ctx context.Context, in *mnemosynerpc.GetRequest, opts ...grpc.CallOption) (*mnemosynerpc.GetResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.GetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.GetRequest, ...grpc.CallOption) *mnemosynerpc.GetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.GetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.GetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) List(ctx context.Context, in *mnemosynerpc.ListRequest, opts ...grpc.CallOption) (*mnemosynerpc.ListResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.ListResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.ListRequest, ...grpc.CallOption) *mnemosynerpc.ListResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ListResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.ListRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Exists(ctx context.Context, in *mnemosynerpc.ExistsRequest, opts ...grpc.CallOption) (*mnemosynerpc.ExistsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.ExistsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.ExistsRequest, ...grpc.CallOption) *mnemosynerpc.ExistsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ExistsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.ExistsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Start(ctx context.Context, in *mnemosynerpc.StartRequest, opts ...grpc.CallOption) (*mnemosynerpc.StartResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.StartResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.StartRequest, ...grpc.CallOption) *mnemosynerpc.StartResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.StartResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.StartRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Abandon(ctx context.Context, in *mnemosynerpc.AbandonRequest, opts ...grpc.CallOption) (*mnemosynerpc.AbandonResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.AbandonResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.AbandonRequest, ...grpc.CallOption) *mnemosynerpc.AbandonResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.AbandonResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.AbandonRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) SetValue(ctx context.Context, in *mnemosynerpc.SetValueRequest, opts ...grpc.CallOption) (*mnemosynerpc.SetValueResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.SetValueResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.SetValueRequest, ...grpc.CallOption) *mnemosynerpc.SetValueResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.SetValueResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.SetValueRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Delete(ctx context.Context, in *mnemosynerpc.DeleteRequest, opts ...grpc.CallOption) (*mnemosynerpc.DeleteResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.DeleteResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.DeleteRequest, ...grpc.CallOption) *mnemosynerpc.DeleteResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.DeleteResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.DeleteRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type RPCServer struct {
	mock.Mock
}

// Context provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Context(_a0 context.Context, _a1 *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.ContextResponse
	if rf, ok := ret.Get(0).(func(context.Context, *empty.Empty) *mnemosynerpc.ContextResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ContextResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *empty.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Get(_a0 context.Context, _a1 *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.GetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.GetRequest) *mnemosynerpc.GetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.GetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.GetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) List(_a0 context.Context, _a1 *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.ListResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.ListRequest) *mnemosynerpc.ListResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ListResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.ListRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Exists(_a0 context.Context, _a1 *mnemosynerpc.ExistsRequest) (*mnemosynerpc.ExistsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.ExistsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.ExistsRequest) *mnemosynerpc.ExistsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ExistsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.ExistsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Start(_a0 context.Context, _a1 *mnemosynerpc.StartRequest) (*mnemosynerpc.StartResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.StartResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.StartRequest) *mnemosynerpc.StartResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.StartResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.StartRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Abandon(_a0 context.Context, _a1 *mnemosynerpc.AbandonRequest) (*mnemosynerpc.AbandonResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.AbandonResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.AbandonRequest) *mnemosynerpc.AbandonResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.AbandonResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.AbandonRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) SetValue(_a0 context.Context, _a1 *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.SetValueResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.SetValueRequest) *mnemosynerpc.SetValueResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.SetValueResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.SetValueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Delete(_a0 context.Context, _a1 *mnemosynerpc.DeleteRequest) (*mnemosynerpc.DeleteResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosynerpc.DeleteResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.DeleteRequest) *mnemosynerpc.DeleteResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.DeleteResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosynerpc.DeleteRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type RandomBytesGenerator struct {
	mock.Mock
}

// GenerateRandomBytes provides a mock function with given fields: _a0
func (_m *RandomBytesGenerator) GenerateRandomBytes(_a0 int) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(int) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
