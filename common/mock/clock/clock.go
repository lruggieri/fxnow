// Code generated by mockery v2.20.2. DO NOT EDIT.

package mockclock

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Clock is an autogenerated mock type for the Clock type
type Clock struct {
	mock.Mock
}

type Clock_Expecter struct {
	mock *mock.Mock
}

func (_m *Clock) EXPECT() *Clock_Expecter {
	return &Clock_Expecter{mock: &_m.Mock}
}

// Now provides a mock function with given fields:
func (_m *Clock) Now() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// Clock_Now_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Now'
type Clock_Now_Call struct {
	*mock.Call
}

// Now is a helper method to define mock.On call
func (_e *Clock_Expecter) Now() *Clock_Now_Call {
	return &Clock_Now_Call{Call: _e.mock.On("Now")}
}

func (_c *Clock_Now_Call) Run(run func()) *Clock_Now_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Clock_Now_Call) Return(_a0 time.Time) *Clock_Now_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Clock_Now_Call) RunAndReturn(run func() time.Time) *Clock_Now_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewClock interface {
	mock.TestingT
	Cleanup(func())
}

// NewClock creates a new instance of Clock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClock(t mockConstructorTestingTNewClock) *Clock {
	mock := &Clock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
