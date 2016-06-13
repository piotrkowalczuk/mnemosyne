package mnemosynetest

import "github.com/stretchr/testify/mock"

import "github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
import "golang.org/x/net/context"

type Mnemosyne struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Mnemosyne) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FromContext provides a mock function with given fields: ctx
func (_m *Mnemosyne) FromContext(ctx context.Context) (*mnemosynerpc.Session, error) {
	ret := _m.Called(ctx)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(context.Context) *mnemosynerpc.Session); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, token
func (_m *Mnemosyne) Get(ctx context.Context, token string) (*mnemosynerpc.Session, error) {
	ret := _m.Called(ctx, token)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(context.Context, string) *mnemosynerpc.Session); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: ctx, subjectID, subjectClient, bag
func (_m *Mnemosyne) Start(ctx context.Context, subjectID string, subjectClient string, bag map[string]string) (*mnemosynerpc.Session, error) {
	ret := _m.Called(ctx, subjectID, subjectClient, bag)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]string) *mnemosynerpc.Session); ok {
		r0 = rf(ctx, subjectID, subjectClient, bag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, map[string]string) error); ok {
		r1 = rf(ctx, subjectID, subjectClient, bag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: ctx, token
func (_m *Mnemosyne) Exists(ctx context.Context, token string) (bool, error) {
	ret := _m.Called(ctx, token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: ctx, token
func (_m *Mnemosyne) Abandon(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetValue provides a mock function with given fields: ctx, token, key, value
func (_m *Mnemosyne) SetValue(ctx context.Context, token string, key string, value string) (map[string]string, error) {
	ret := _m.Called(ctx, token, key, value)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) map[string]string); ok {
		r0 = rf(ctx, token, key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, token, key, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
