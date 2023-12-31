// Code generated by mockery v2.20.2. DO NOT EDIT.

package mockfxsource

import (
	context "context"

	fxsource "github.com/lruggieri/fxnow/common/fxsource"
	mock "github.com/stretchr/testify/mock"
)

// FXSource is an autogenerated mock type for the FXSource type
type FXSource struct {
	mock.Mock
}

type FXSource_Expecter struct {
	mock *mock.Mock
}

func (_m *FXSource) EXPECT() *FXSource_Expecter {
	return &FXSource_Expecter{mock: &_m.Mock}
}

// FetchAllRates provides a mock function with given fields: _a0, _a1
func (_m *FXSource) FetchAllRates(_a0 context.Context, _a1 fxsource.FetchAllRatesRequest) (*fxsource.FetchAllRatesResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *fxsource.FetchAllRatesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, fxsource.FetchAllRatesRequest) (*fxsource.FetchAllRatesResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, fxsource.FetchAllRatesRequest) *fxsource.FetchAllRatesResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fxsource.FetchAllRatesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, fxsource.FetchAllRatesRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FXSource_FetchAllRates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchAllRates'
type FXSource_FetchAllRates_Call struct {
	*mock.Call
}

// FetchAllRates is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 fxsource.FetchAllRatesRequest
func (_e *FXSource_Expecter) FetchAllRates(_a0 interface{}, _a1 interface{}) *FXSource_FetchAllRates_Call {
	return &FXSource_FetchAllRates_Call{Call: _e.mock.On("FetchAllRates", _a0, _a1)}
}

func (_c *FXSource_FetchAllRates_Call) Run(run func(_a0 context.Context, _a1 fxsource.FetchAllRatesRequest)) *FXSource_FetchAllRates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(fxsource.FetchAllRatesRequest))
	})
	return _c
}

func (_c *FXSource_FetchAllRates_Call) Return(_a0 *fxsource.FetchAllRatesResponse, _a1 error) *FXSource_FetchAllRates_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FXSource_FetchAllRates_Call) RunAndReturn(run func(context.Context, fxsource.FetchAllRatesRequest) (*fxsource.FetchAllRatesResponse, error)) *FXSource_FetchAllRates_Call {
	_c.Call.Return(run)
	return _c
}

// FetchRate provides a mock function with given fields: _a0, _a1
func (_m *FXSource) FetchRate(_a0 context.Context, _a1 fxsource.FetchRateRequest) (*fxsource.FetchRateResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *fxsource.FetchRateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, fxsource.FetchRateRequest) (*fxsource.FetchRateResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, fxsource.FetchRateRequest) *fxsource.FetchRateResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fxsource.FetchRateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, fxsource.FetchRateRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FXSource_FetchRate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchRate'
type FXSource_FetchRate_Call struct {
	*mock.Call
}

// FetchRate is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 fxsource.FetchRateRequest
func (_e *FXSource_Expecter) FetchRate(_a0 interface{}, _a1 interface{}) *FXSource_FetchRate_Call {
	return &FXSource_FetchRate_Call{Call: _e.mock.On("FetchRate", _a0, _a1)}
}

func (_c *FXSource_FetchRate_Call) Run(run func(_a0 context.Context, _a1 fxsource.FetchRateRequest)) *FXSource_FetchRate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(fxsource.FetchRateRequest))
	})
	return _c
}

func (_c *FXSource_FetchRate_Call) Return(_a0 *fxsource.FetchRateResponse, _a1 error) *FXSource_FetchRate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FXSource_FetchRate_Call) RunAndReturn(run func(context.Context, fxsource.FetchRateRequest) (*fxsource.FetchRateResponse, error)) *FXSource_FetchRate_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewFXSource interface {
	mock.TestingT
	Cleanup(func())
}

// NewFXSource creates a new instance of FXSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFXSource(t mockConstructorTestingTNewFXSource) *FXSource {
	mock := &FXSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
