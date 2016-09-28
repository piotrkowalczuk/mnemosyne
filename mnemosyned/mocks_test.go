package mnemosyned

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/stretchr/testify/mock"
)

type mockRandomBytesGenerator struct {
	mock.Mock
}

// generateRandomBytes provides a mock function with given fields: _a0
func (_m *mockRandomBytesGenerator) generateRandomBytes(_a0 int) ([]byte, error) {
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

type mockStorage struct {
	mock.Mock
}

// Setup provides a mock function with given fields:
func (_m *mockStorage) Setup() error {
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
func (_m *mockStorage) TearDown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *mockStorage) Start(_a0 string, _a1 string, _a2 string, _a3 map[string]string) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(string, string, string, map[string]string) *mnemosynerpc.Session); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, map[string]string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Abandon provides a mock function with given fields: _a0
func (_m *mockStorage) Abandon(_a0 string) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *mockStorage) Get(_a0 string) (*mnemosynerpc.Session, error) {
	ret := _m.Called(_a0)

	var r0 *mnemosynerpc.Session
	if rf, ok := ret.Get(0).(func(string) *mnemosynerpc.Session); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosynerpc.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *mockStorage) List(_a0 int64, _a1 int64, _a2 *time.Time, _a3 *time.Time) ([]*mnemosynerpc.Session, error) {
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
func (_m *mockStorage) Exists(_a0 string) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockStorage) Delete(_a0 string, _a1 *time.Time, _a2 *time.Time) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, *time.Time, *time.Time) int64); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *time.Time, *time.Time) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockStorage) SetValue(_a0 string, _a1 string, _a2 string) (map[string]string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(string, string, string) map[string]string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
