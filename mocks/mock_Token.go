// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockToken is an autogenerated mock type for the Token type
type MockToken struct {
	mock.Mock
}

type MockToken_Expecter struct {
	mock *mock.Mock
}

func (_m *MockToken) EXPECT() *MockToken_Expecter {
	return &MockToken_Expecter{mock: &_m.Mock}
}

// Done provides a mock function with given fields:
func (_m *MockToken) Done() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Done")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// MockToken_Done_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Done'
type MockToken_Done_Call struct {
	*mock.Call
}

// Done is a helper method to define mock.On call
func (_e *MockToken_Expecter) Done() *MockToken_Done_Call {
	return &MockToken_Done_Call{Call: _e.mock.On("Done")}
}

func (_c *MockToken_Done_Call) Run(run func()) *MockToken_Done_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockToken_Done_Call) Return(_a0 <-chan struct{}) *MockToken_Done_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockToken_Done_Call) RunAndReturn(run func() <-chan struct{}) *MockToken_Done_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields:
func (_m *MockToken) Error() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockToken_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type MockToken_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
func (_e *MockToken_Expecter) Error() *MockToken_Error_Call {
	return &MockToken_Error_Call{Call: _e.mock.On("Error")}
}

func (_c *MockToken_Error_Call) Run(run func()) *MockToken_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockToken_Error_Call) Return(_a0 error) *MockToken_Error_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockToken_Error_Call) RunAndReturn(run func() error) *MockToken_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Wait provides a mock function with given fields:
func (_m *MockToken) Wait() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Wait")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockToken_Wait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Wait'
type MockToken_Wait_Call struct {
	*mock.Call
}

// Wait is a helper method to define mock.On call
func (_e *MockToken_Expecter) Wait() *MockToken_Wait_Call {
	return &MockToken_Wait_Call{Call: _e.mock.On("Wait")}
}

func (_c *MockToken_Wait_Call) Run(run func()) *MockToken_Wait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockToken_Wait_Call) Return(_a0 bool) *MockToken_Wait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockToken_Wait_Call) RunAndReturn(run func() bool) *MockToken_Wait_Call {
	_c.Call.Return(run)
	return _c
}

// WaitTimeout provides a mock function with given fields: _a0
func (_m *MockToken) WaitTimeout(_a0 time.Duration) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for WaitTimeout")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(time.Duration) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockToken_WaitTimeout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WaitTimeout'
type MockToken_WaitTimeout_Call struct {
	*mock.Call
}

// WaitTimeout is a helper method to define mock.On call
//   - _a0 time.Duration
func (_e *MockToken_Expecter) WaitTimeout(_a0 interface{}) *MockToken_WaitTimeout_Call {
	return &MockToken_WaitTimeout_Call{Call: _e.mock.On("WaitTimeout", _a0)}
}

func (_c *MockToken_WaitTimeout_Call) Run(run func(_a0 time.Duration)) *MockToken_WaitTimeout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration))
	})
	return _c
}

func (_c *MockToken_WaitTimeout_Call) Return(_a0 bool) *MockToken_WaitTimeout_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockToken_WaitTimeout_Call) RunAndReturn(run func(time.Duration) bool) *MockToken_WaitTimeout_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockToken creates a new instance of MockToken. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockToken(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockToken {
	mock := &MockToken{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
