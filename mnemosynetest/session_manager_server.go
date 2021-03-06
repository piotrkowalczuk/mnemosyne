// Code generated by mockery v1.0.0. DO NOT EDIT.

package mnemosynetest

import context "context"
import empty "github.com/golang/protobuf/ptypes/empty"
import mnemosynerpc "github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
import mock "github.com/stretchr/testify/mock"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

// SessionManagerServer is an autogenerated mock type for the SessionManagerServer type
type SessionManagerServer struct {
	mock.Mock
}

// Abandon provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) Abandon(_a0 context.Context, _a1 *mnemosynerpc.AbandonRequest) (*wrappers.BoolValue, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *wrappers.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.AbandonRequest) *wrappers.BoolValue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wrappers.BoolValue)
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

// Context provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) Context(_a0 context.Context, _a1 *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
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

// Delete provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) Delete(_a0 context.Context, _a1 *mnemosynerpc.DeleteRequest) (*wrappers.Int64Value, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *wrappers.Int64Value
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.DeleteRequest) *wrappers.Int64Value); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wrappers.Int64Value)
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

// Exists provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) Exists(_a0 context.Context, _a1 *mnemosynerpc.ExistsRequest) (*wrappers.BoolValue, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *wrappers.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *mnemosynerpc.ExistsRequest) *wrappers.BoolValue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wrappers.BoolValue)
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

// Get provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) Get(_a0 context.Context, _a1 *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
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
func (_m *SessionManagerServer) List(_a0 context.Context, _a1 *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
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

// SetValue provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) SetValue(_a0 context.Context, _a1 *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
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

// Start provides a mock function with given fields: _a0, _a1
func (_m *SessionManagerServer) Start(_a0 context.Context, _a1 *mnemosynerpc.StartRequest) (*mnemosynerpc.StartResponse, error) {
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
