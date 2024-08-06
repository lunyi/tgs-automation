// Code generated by MockGen. DO NOT EDIT.
// Source: tgs-automation/pkg/namecheap (interfaces: NamecheapApi)
//
// Generated by this command:
//
//	mockgen -destination pkg/namecheap/mocks/mock_namecheap.go -package=namecheap tgs-automation/pkg/namecheap NamecheapApi
//

// Package namecheap is a generated GoMock package.
package namecheap

import (
	context "context"
	reflect "reflect"
	namecheap "tgs-automation/pkg/namecheap"

	gomock "go.uber.org/mock/gomock"
)

// MockNamecheapApi is a mock of NamecheapApi interface.
type MockNamecheapApi struct {
	ctrl     *gomock.Controller
	recorder *MockNamecheapApiMockRecorder
}

// MockNamecheapApiMockRecorder is the mock recorder for MockNamecheapApi.
type MockNamecheapApiMockRecorder struct {
	mock *MockNamecheapApi
}

// NewMockNamecheapApi creates a new mock instance.
func NewMockNamecheapApi(ctrl *gomock.Controller) *MockNamecheapApi {
	mock := &MockNamecheapApi{ctrl: ctrl}
	mock.recorder = &MockNamecheapApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNamecheapApi) EXPECT() *MockNamecheapApiMockRecorder {
	return m.recorder
}

// CheckDomainAvailable mocks base method.
func (m *MockNamecheapApi) CheckDomainAvailable(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDomainAvailable", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDomainAvailable indicates an expected call of CheckDomainAvailable.
func (mr *MockNamecheapApiMockRecorder) CheckDomainAvailable(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDomainAvailable", reflect.TypeOf((*MockNamecheapApi)(nil).CheckDomainAvailable), arg0, arg1)
}

// CreateDomain mocks base method.
func (m *MockNamecheapApi) CreateDomain(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDomain", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDomain indicates an expected call of CreateDomain.
func (mr *MockNamecheapApiMockRecorder) CreateDomain(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDomain", reflect.TypeOf((*MockNamecheapApi)(nil).CreateDomain), arg0, arg1, arg2)
}

// GetBalance mocks base method.
func (m *MockNamecheapApi) GetBalance(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockNamecheapApiMockRecorder) GetBalance(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockNamecheapApi)(nil).GetBalance), arg0)
}

// GetCouponCode mocks base method.
func (m *MockNamecheapApi) GetCouponCode(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCouponCode", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCouponCode indicates an expected call of GetCouponCode.
func (mr *MockNamecheapApiMockRecorder) GetCouponCode(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCouponCode", reflect.TypeOf((*MockNamecheapApi)(nil).GetCouponCode), arg0)
}

// GetDomainPrice mocks base method.
func (m *MockNamecheapApi) GetDomainPrice(arg0 context.Context, arg1 string) (*namecheap.CheckDomainPriceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainPrice", arg0, arg1)
	ret0, _ := ret[0].(*namecheap.CheckDomainPriceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainPrice indicates an expected call of GetDomainPrice.
func (mr *MockNamecheapApiMockRecorder) GetDomainPrice(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainPrice", reflect.TypeOf((*MockNamecheapApi)(nil).GetDomainPrice), arg0, arg1)
}

// GetExpiredDomains mocks base method.
func (m *MockNamecheapApi) GetExpiredDomains() ([]namecheap.FilteredDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpiredDomains")
	ret0, _ := ret[0].([]namecheap.FilteredDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExpiredDomains indicates an expected call of GetExpiredDomains.
func (mr *MockNamecheapApiMockRecorder) GetExpiredDomains() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpiredDomains", reflect.TypeOf((*MockNamecheapApi)(nil).GetExpiredDomains))
}

// UpdateNameServer mocks base method.
func (m *MockNamecheapApi) UpdateNameServer(arg0, arg1 string) (*namecheap.UpdateNameServerApiResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNameServer", arg0, arg1)
	ret0, _ := ret[0].(*namecheap.UpdateNameServerApiResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNameServer indicates an expected call of UpdateNameServer.
func (mr *MockNamecheapApiMockRecorder) UpdateNameServer(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNameServer", reflect.TypeOf((*MockNamecheapApi)(nil).UpdateNameServer), arg0, arg1)
}