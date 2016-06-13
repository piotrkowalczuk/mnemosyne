package mnemosynetest

import "github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
import "github.com/stretchr/testify/mock"

import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"
import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

type SessionManagerClient struct {
	mock.Mock
}

// Context provides a mock function with given fields: ctx, in, opts
func (_m *SessionManagerClient) Context(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*mnemosynerpc.ContextResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *mnemosynerpc.ContextResponse
	if rf, ok := ret.Get(0).(func(context.Context, *google_protobuf1.Empty, ...grpc.CallOption) *mnemosynerpc.ContextResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.ContextResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *google_protobuf1.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *SessionManagerClient) Get(ctx context.Context, in *mnemosynerpc.GetRequest, opts ...grpc.CallOption) (*mnemosynerpc.GetResponse, error) {
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
func (_m *SessionManagerClient) List(ctx context.Context, in *mnemosynerpc.ListRequest, opts ...grpc.CallOption) (*mnemosynerpc.ListResponse, error) {
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
func (_m *SessionManagerClient) Exists(ctx context.Context, in *mnemosynerpc.ExistsRequest, opts ...grpc.CallOption) (*mnemosynerpc.ExistsResponse, error) {
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
func (_m *SessionManagerClient) Start(ctx context.Context, in *mnemosynerpc.StartRequest, opts ...grpc.CallOption) (*mnemosynerpc.StartResponse, error) {
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
func (_m *SessionManagerClient) Abandon(ctx context.Context, in *mnemosynerpc.AbandonRequest, opts ...grpc.CallOption) (*mnemosynerpc.AbandonResponse, error) {
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
func (_m *SessionManagerClient) SetValue(ctx context.Context, in *mnemosynerpc.SetValueRequest, opts ...grpc.CallOption) (*mnemosynerpc.SetValueResponse, error) {
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
func (_m *SessionManagerClient) Delete(ctx context.Context, in *mnemosynerpc.DeleteRequest, opts ...grpc.CallOption) (*mnemosynerpc.DeleteResponse, error) {
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
