// Code generated by MockGen. DO NOT EDIT.
// Source: ./solution.go

// Package mock_solution is a generated GoMock package.
package mock_solution

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	reflect "reflect"
)

// MockGateway is a mock of Gateway interface
type MockGateway struct {
	ctrl     *gomock.Controller
	recorder *MockGatewayMockRecorder
}

// MockGatewayMockRecorder is the mock recorder for MockGateway
type MockGatewayMockRecorder struct {
	mock *MockGateway
}

// NewMockGateway creates a new mock instance
func NewMockGateway(ctrl *gomock.Controller) *MockGateway {
	mock := &MockGateway{ctrl: ctrl}
	mock.recorder = &MockGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGateway) EXPECT() *MockGatewayMockRecorder {
	return m.recorder
}

// GetSolution mocks base method
func (m *MockGateway) GetSolution(id int64) (*models.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolution", id)
	ret0, _ := ret[0].(*models.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolution indicates an expected call of GetSolution
func (mr *MockGatewayMockRecorder) GetSolution(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolution", reflect.TypeOf((*MockGateway)(nil).GetSolution), id)
}