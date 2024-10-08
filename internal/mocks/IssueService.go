// Code generated by mockery v2.45.0. DO NOT EDIT.

package interfaces

import (
	jira "github.com/andygrunwald/go-jira"
	mock "github.com/stretchr/testify/mock"

	work "github.com/manomartins/bitbird/internal/app/work"
)

// MockIssueService is an autogenerated mock type for the IssueService type
type MockIssueService struct {
	mock.Mock
}

type MockIssueService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIssueService) EXPECT() *MockIssueService_Expecter {
	return &MockIssueService_Expecter{mock: &_m.Mock}
}

// GetFirstIssueByCodeBase provides a mock function with given fields: base
func (_m *MockIssueService) GetFirstIssueByCodeBase(base work.CodeBase) *jira.Issue {
	ret := _m.Called(base)

	if len(ret) == 0 {
		panic("no return value specified for GetFirstIssueByCodeBase")
	}

	var r0 *jira.Issue
	if rf, ok := ret.Get(0).(func(work.CodeBase) *jira.Issue); ok {
		r0 = rf(base)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jira.Issue)
		}
	}

	return r0
}

// MockIssueService_GetFirstIssueByCodeBase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFirstIssueByCodeBase'
type MockIssueService_GetFirstIssueByCodeBase_Call struct {
	*mock.Call
}

// GetFirstIssueByCodeBase is a helper method to define mock.On call
//   - base work.CodeBase
func (_e *MockIssueService_Expecter) GetFirstIssueByCodeBase(base interface{}) *MockIssueService_GetFirstIssueByCodeBase_Call {
	return &MockIssueService_GetFirstIssueByCodeBase_Call{Call: _e.mock.On("GetFirstIssueByCodeBase", base)}
}

func (_c *MockIssueService_GetFirstIssueByCodeBase_Call) Run(run func(base work.CodeBase)) *MockIssueService_GetFirstIssueByCodeBase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(work.CodeBase))
	})
	return _c
}

func (_c *MockIssueService_GetFirstIssueByCodeBase_Call) Return(_a0 *jira.Issue) *MockIssueService_GetFirstIssueByCodeBase_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssueService_GetFirstIssueByCodeBase_Call) RunAndReturn(run func(work.CodeBase) *jira.Issue) *MockIssueService_GetFirstIssueByCodeBase_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIssueService creates a new instance of MockIssueService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIssueService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIssueService {
	mock := &MockIssueService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
