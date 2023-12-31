// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockMessage is an autogenerated mock type for the Message type
type MockMessage struct {
	mock.Mock
}

type MockMessage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessage) EXPECT() *MockMessage_Expecter {
	return &MockMessage_Expecter{mock: &_m.Mock}
}

// Ack provides a mock function with given fields:
func (_m *MockMessage) Ack() {
	_m.Called()
}

// MockMessage_Ack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ack'
type MockMessage_Ack_Call struct {
	*mock.Call
}

// Ack is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Ack() *MockMessage_Ack_Call {
	return &MockMessage_Ack_Call{Call: _e.mock.On("Ack")}
}

func (_c *MockMessage_Ack_Call) Run(run func()) *MockMessage_Ack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Ack_Call) Return() *MockMessage_Ack_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockMessage_Ack_Call) RunAndReturn(run func()) *MockMessage_Ack_Call {
	_c.Call.Return(run)
	return _c
}

// Duplicate provides a mock function with given fields:
func (_m *MockMessage) Duplicate() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Duplicate")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockMessage_Duplicate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Duplicate'
type MockMessage_Duplicate_Call struct {
	*mock.Call
}

// Duplicate is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Duplicate() *MockMessage_Duplicate_Call {
	return &MockMessage_Duplicate_Call{Call: _e.mock.On("Duplicate")}
}

func (_c *MockMessage_Duplicate_Call) Run(run func()) *MockMessage_Duplicate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Duplicate_Call) Return(_a0 bool) *MockMessage_Duplicate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Duplicate_Call) RunAndReturn(run func() bool) *MockMessage_Duplicate_Call {
	_c.Call.Return(run)
	return _c
}

// MessageID provides a mock function with given fields:
func (_m *MockMessage) MessageID() uint16 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MessageID")
	}

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// MockMessage_MessageID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MessageID'
type MockMessage_MessageID_Call struct {
	*mock.Call
}

// MessageID is a helper method to define mock.On call
func (_e *MockMessage_Expecter) MessageID() *MockMessage_MessageID_Call {
	return &MockMessage_MessageID_Call{Call: _e.mock.On("MessageID")}
}

func (_c *MockMessage_MessageID_Call) Run(run func()) *MockMessage_MessageID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_MessageID_Call) Return(_a0 uint16) *MockMessage_MessageID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_MessageID_Call) RunAndReturn(run func() uint16) *MockMessage_MessageID_Call {
	_c.Call.Return(run)
	return _c
}

// Payload provides a mock function with given fields:
func (_m *MockMessage) Payload() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Payload")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockMessage_Payload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Payload'
type MockMessage_Payload_Call struct {
	*mock.Call
}

// Payload is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Payload() *MockMessage_Payload_Call {
	return &MockMessage_Payload_Call{Call: _e.mock.On("Payload")}
}

func (_c *MockMessage_Payload_Call) Run(run func()) *MockMessage_Payload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Payload_Call) Return(_a0 []byte) *MockMessage_Payload_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Payload_Call) RunAndReturn(run func() []byte) *MockMessage_Payload_Call {
	_c.Call.Return(run)
	return _c
}

// Qos provides a mock function with given fields:
func (_m *MockMessage) Qos() byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Qos")
	}

	var r0 byte
	if rf, ok := ret.Get(0).(func() byte); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(byte)
	}

	return r0
}

// MockMessage_Qos_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Qos'
type MockMessage_Qos_Call struct {
	*mock.Call
}

// Qos is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Qos() *MockMessage_Qos_Call {
	return &MockMessage_Qos_Call{Call: _e.mock.On("Qos")}
}

func (_c *MockMessage_Qos_Call) Run(run func()) *MockMessage_Qos_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Qos_Call) Return(_a0 byte) *MockMessage_Qos_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Qos_Call) RunAndReturn(run func() byte) *MockMessage_Qos_Call {
	_c.Call.Return(run)
	return _c
}

// Retained provides a mock function with given fields:
func (_m *MockMessage) Retained() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Retained")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockMessage_Retained_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Retained'
type MockMessage_Retained_Call struct {
	*mock.Call
}

// Retained is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Retained() *MockMessage_Retained_Call {
	return &MockMessage_Retained_Call{Call: _e.mock.On("Retained")}
}

func (_c *MockMessage_Retained_Call) Run(run func()) *MockMessage_Retained_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Retained_Call) Return(_a0 bool) *MockMessage_Retained_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Retained_Call) RunAndReturn(run func() bool) *MockMessage_Retained_Call {
	_c.Call.Return(run)
	return _c
}

// Topic provides a mock function with given fields:
func (_m *MockMessage) Topic() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Topic")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockMessage_Topic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Topic'
type MockMessage_Topic_Call struct {
	*mock.Call
}

// Topic is a helper method to define mock.On call
func (_e *MockMessage_Expecter) Topic() *MockMessage_Topic_Call {
	return &MockMessage_Topic_Call{Call: _e.mock.On("Topic")}
}

func (_c *MockMessage_Topic_Call) Run(run func()) *MockMessage_Topic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_Topic_Call) Return(_a0 string) *MockMessage_Topic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_Topic_Call) RunAndReturn(run func() string) *MockMessage_Topic_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMessage creates a new instance of MockMessage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessage {
	mock := &MockMessage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
