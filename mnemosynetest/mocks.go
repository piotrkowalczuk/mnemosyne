package mnemosynetest

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
)
import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"
import "google.golang.org/grpc"

type Mnemosyne struct {
	mock.Mock
}

// FromContext provides a mock function with given fields: _a0
func (_m *Mnemosyne) FromContext(_a0 context.Context) (*mnemosyne.Session, error) {
	ret := _m.Called(_a0)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(context.Context) *mnemosyne.Session); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
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
func (_m *Mnemosyne) Get(_a0 context.Context, _a1 mnemosyne.Token) (*mnemosyne.Session, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.Token) *mnemosyne.Session); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosyne.Token) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: _a0, _a1
func (_m *Mnemosyne) Exists(_a0 context.Context, _a1 mnemosyne.Token) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.Token) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosyne.Token) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0, _a1, _a2
func (_m *Mnemosyne) Start(_a0 context.Context, _a1 string, _a2 map[string]string) (*mnemosyne.Session, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string) *mnemosyne.Session); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: _a0, _a1
func (_m *Mnemosyne) Abandon(_a0 context.Context, _a1 mnemosyne.Token) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.Token) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetValue provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Mnemosyne) SetValue(_a0 context.Context, _a1 mnemosyne.Token, _a2 string, _a3 string) (map[string]string, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.Token, string, string) map[string]string); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosyne.Token, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type RPCClient struct {
	mock.Mock
}

// Context provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Context(ctx context.Context, in *mnemosyne.Empty, opts ...grpc.CallOption) (*mnemosyne.Session, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.Empty, ...grpc.CallOption) *mnemosyne.Session); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Get(ctx context.Context, in *mnemosyne.GetRequest, opts ...grpc.CallOption) (*mnemosyne.GetResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.GetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.GetRequest, ...grpc.CallOption) *mnemosyne.GetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.GetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.GetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) List(ctx context.Context, in *mnemosyne.ListRequest, opts ...grpc.CallOption) (*mnemosyne.ListResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.ListResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.ListRequest, ...grpc.CallOption) *mnemosyne.ListResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.ListResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.ListRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Exists(ctx context.Context, in *mnemosyne.ExistsRequest, opts ...grpc.CallOption) (*mnemosyne.ExistsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.ExistsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.ExistsRequest, ...grpc.CallOption) *mnemosyne.ExistsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.ExistsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.ExistsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Start(ctx context.Context, in *mnemosyne.StartRequest, opts ...grpc.CallOption) (*mnemosyne.StartResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.StartResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.StartRequest, ...grpc.CallOption) *mnemosyne.StartResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.StartResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.StartRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Abandon(ctx context.Context, in *mnemosyne.AbandonRequest, opts ...grpc.CallOption) (*mnemosyne.AbandonResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.AbandonResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.AbandonRequest, ...grpc.CallOption) *mnemosyne.AbandonResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.AbandonResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.AbandonRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) SetValue(ctx context.Context, in *mnemosyne.SetValueRequest, opts ...grpc.CallOption) (*mnemosyne.SetValueResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.SetValueResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.SetValueRequest, ...grpc.CallOption) *mnemosyne.SetValueResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.SetValueResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.SetValueRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *RPCClient) Delete(ctx context.Context, in *mnemosyne.DeleteRequest, opts ...grpc.CallOption) (*mnemosyne.DeleteResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosyne.DeleteResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.DeleteRequest, ...grpc.CallOption) *mnemosyne.DeleteResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.DeleteResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.DeleteRequest, ...grpc.CallOption) error); ok {
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
func (_m *RPCServer) Context(_a0 context.Context, _a1 *mnemosyne.Empty) (*mnemosyne.Session, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.Empty) *mnemosyne.Session); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Get(_a0 context.Context, _a1 *mnemosyne.GetRequest) (*mnemosyne.GetResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.GetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.GetRequest) *mnemosyne.GetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.GetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.GetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) List(_a0 context.Context, _a1 *mnemosyne.ListRequest) (*mnemosyne.ListResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.ListResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.ListRequest) *mnemosyne.ListResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.ListResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.ListRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Exists(_a0 context.Context, _a1 *mnemosyne.ExistsRequest) (*mnemosyne.ExistsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.ExistsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.ExistsRequest) *mnemosyne.ExistsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.ExistsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.ExistsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Start(_a0 context.Context, _a1 *mnemosyne.StartRequest) (*mnemosyne.StartResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.StartResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.StartRequest) *mnemosyne.StartResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.StartResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.StartRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Abandon(_a0 context.Context, _a1 *mnemosyne.AbandonRequest) (*mnemosyne.AbandonResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.AbandonResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.AbandonRequest) *mnemosyne.AbandonResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.AbandonResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.AbandonRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) SetValue(_a0 context.Context, _a1 *mnemosyne.SetValueRequest) (*mnemosyne.SetValueResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.SetValueResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.SetValueRequest) *mnemosyne.SetValueResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.SetValueResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.SetValueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *RPCServer) Delete(_a0 context.Context, _a1 *mnemosyne.DeleteRequest) (*mnemosyne.DeleteResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.DeleteResponse
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosyne.DeleteRequest) *mnemosyne.DeleteResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.DeleteResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *mnemosyne.DeleteRequest) error); ok {
		r1 = rf(_a0, _a1)
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
func (_m *Storage) Start(_a0 string, _a1 map[string]string) (*mnemosyne.Session, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(string, map[string]string) *mnemosyne.Session); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
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
func (_m *Storage) Abandon(_a0 *mnemosyne.Token) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*mnemosyne.Token) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosyne.Token) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *Storage) Get(_a0 *mnemosyne.Token) (*mnemosyne.Session, error) {
	ret := _m.Called(_a0)

	var r0 *mnemosyne.Session
	if rf, ok := ret.Get(0).(func(*mnemosyne.Token) *mnemosyne.Session); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosyne.Token) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Storage) List(_a0 int64, _a1 int64, _a2 *time.Time, _a3 *time.Time) ([]*mnemosyne.Session, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []*mnemosyne.Session
	if rf, ok := ret.Get(0).(func(int64, int64, *time.Time, *time.Time) []*mnemosyne.Session); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*mnemosyne.Session)
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
func (_m *Storage) Exists(_a0 *mnemosyne.Token) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*mnemosyne.Token) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosyne.Token) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1, _a2
func (_m *Storage) Delete(_a0 *mnemosyne.Token, _a1 *time.Time, _a2 *time.Time) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*mnemosyne.Token, *time.Time, *time.Time) int64); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosyne.Token, *time.Time, *time.Time) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: _a0, _a1, _a2
func (_m *Storage) SetValue(_a0 *mnemosyne.Token, _a1 string, _a2 string) (map[string]string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(*mnemosyne.Token, string, string) map[string]string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*mnemosyne.Token, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
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
