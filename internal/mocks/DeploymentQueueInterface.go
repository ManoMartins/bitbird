// Code generated by mockery v2.45.0. DO NOT EDIT.

package interfaces

import (
	model "github.com/manomartins/bitbird/internal/app/model"
	mock "github.com/stretchr/testify/mock"
)

// MockDeploymentQueueInterface is an autogenerated mock type for the DeploymentQueueInterface type
type MockDeploymentQueueInterface struct {
	mock.Mock
}

type MockDeploymentQueueInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeploymentQueueInterface) EXPECT() *MockDeploymentQueueInterface_Expecter {
	return &MockDeploymentQueueInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: data
func (_m *MockDeploymentQueueInterface) Create(data model.DeploymentQueueModel) error {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(model.DeploymentQueueModel) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeploymentQueueInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockDeploymentQueueInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - data model.DeploymentQueueModel
func (_e *MockDeploymentQueueInterface_Expecter) Create(data interface{}) *MockDeploymentQueueInterface_Create_Call {
	return &MockDeploymentQueueInterface_Create_Call{Call: _e.mock.On("Create", data)}
}

func (_c *MockDeploymentQueueInterface_Create_Call) Run(run func(data model.DeploymentQueueModel)) *MockDeploymentQueueInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.DeploymentQueueModel))
	})
	return _c
}

func (_c *MockDeploymentQueueInterface_Create_Call) Return(_a0 error) *MockDeploymentQueueInterface_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeploymentQueueInterface_Create_Call) RunAndReturn(run func(model.DeploymentQueueModel) error) *MockDeploymentQueueInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByCardKey provides a mock function with given fields: key
func (_m *MockDeploymentQueueInterface) GetByCardKey(key string) (*model.DeploymentQueueModel, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for GetByCardKey")
	}

	var r0 *model.DeploymentQueueModel
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.DeploymentQueueModel, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) *model.DeploymentQueueModel); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DeploymentQueueModel)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeploymentQueueInterface_GetByCardKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByCardKey'
type MockDeploymentQueueInterface_GetByCardKey_Call struct {
	*mock.Call
}

// GetByCardKey is a helper method to define mock.On call
//   - key string
func (_e *MockDeploymentQueueInterface_Expecter) GetByCardKey(key interface{}) *MockDeploymentQueueInterface_GetByCardKey_Call {
	return &MockDeploymentQueueInterface_GetByCardKey_Call{Call: _e.mock.On("GetByCardKey", key)}
}

func (_c *MockDeploymentQueueInterface_GetByCardKey_Call) Run(run func(key string)) *MockDeploymentQueueInterface_GetByCardKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDeploymentQueueInterface_GetByCardKey_Call) Return(_a0 *model.DeploymentQueueModel, _a1 error) *MockDeploymentQueueInterface_GetByCardKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeploymentQueueInterface_GetByCardKey_Call) RunAndReturn(run func(string) (*model.DeploymentQueueModel, error)) *MockDeploymentQueueInterface_GetByCardKey_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeploymentQueueInterface creates a new instance of MockDeploymentQueueInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeploymentQueueInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeploymentQueueInterface {
	mock := &MockDeploymentQueueInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
