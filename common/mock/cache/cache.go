// Code generated by mockery v2.20.2. DO NOT EDIT.

package mockcache

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

type Cache_Expecter struct {
	mock *mock.Mock
}

func (_m *Cache) EXPECT() *Cache_Expecter {
	return &Cache_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, key, value
func (_m *Cache) Get(ctx context.Context, key string, value interface{}) (bool, error) {
	ret := _m.Called(ctx, key, value)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) (bool, error)); ok {
		return rf(ctx, key, value)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) bool); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}) error); ok {
		r1 = rf(ctx, key, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Cache_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type Cache_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
func (_e *Cache_Expecter) Get(ctx interface{}, key interface{}, value interface{}) *Cache_Get_Call {
	return &Cache_Get_Call{Call: _e.mock.On("Get", ctx, key, value)}
}

func (_c *Cache_Get_Call) Run(run func(ctx context.Context, key string, value interface{})) *Cache_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *Cache_Get_Call) Return(exist bool, err error) *Cache_Get_Call {
	_c.Call.Return(exist, err)
	return _c
}

func (_c *Cache_Get_Call) RunAndReturn(run func(context.Context, string, interface{}) (bool, error)) *Cache_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: ctx, key
func (_m *Cache) Remove(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Cache_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type Cache_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *Cache_Expecter) Remove(ctx interface{}, key interface{}) *Cache_Remove_Call {
	return &Cache_Remove_Call{Call: _e.mock.On("Remove", ctx, key)}
}

func (_c *Cache_Remove_Call) Run(run func(ctx context.Context, key string)) *Cache_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Cache_Remove_Call) Return(_a0 error) *Cache_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Cache_Remove_Call) RunAndReturn(run func(context.Context, string) error) *Cache_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value, expiration
func (_m *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ret := _m.Called(ctx, key, value, expiration)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) error); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Cache_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type Cache_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
//   - expiration time.Duration
func (_e *Cache_Expecter) Set(ctx interface{}, key interface{}, value interface{}, expiration interface{}) *Cache_Set_Call {
	return &Cache_Set_Call{Call: _e.mock.On("Set", ctx, key, value, expiration)}
}

func (_c *Cache_Set_Call) Run(run func(ctx context.Context, key string, value interface{}, expiration time.Duration)) *Cache_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), args[3].(time.Duration))
	})
	return _c
}

func (_c *Cache_Set_Call) Return(err error) *Cache_Set_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Cache_Set_Call) RunAndReturn(run func(context.Context, string, interface{}, time.Duration) error) *Cache_Set_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewCache creates a new instance of Cache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCache(t mockConstructorTestingTNewCache) *Cache {
	mock := &Cache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
