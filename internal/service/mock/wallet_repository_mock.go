// Code generated by MockGen. DO NOT EDIT.
// Source: wallet_repository.go
//
// Generated by this command:
//
//	mockgen -source=wallet_repository.go -destination=mock/wallet_repository_mock.go -package service_mock
//

// Package service_mock is a generated GoMock package.
package service_mock

import (
	context "context"
	reflect "reflect"

	repository "github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockWalletRepository is a mock of WalletRepository interface.
type MockWalletRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWalletRepositoryMockRecorder
}

// MockWalletRepositoryMockRecorder is the mock recorder for MockWalletRepository.
type MockWalletRepositoryMockRecorder struct {
	mock *MockWalletRepository
}

// NewMockWalletRepository creates a new mock instance.
func NewMockWalletRepository(ctrl *gomock.Controller) *MockWalletRepository {
	mock := &MockWalletRepository{ctrl: ctrl}
	mock.recorder = &MockWalletRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepository) EXPECT() *MockWalletRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWalletRepository) Create(ctx context.Context) (*repository.DBWallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx)
	ret0, _ := ret[0].(*repository.DBWallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockWalletRepositoryMockRecorder) Create(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWalletRepository)(nil).Create), ctx)
}

// GetByWalletID mocks base method.
func (m *MockWalletRepository) GetByWalletID(ctx context.Context, walletID string) (*repository.DBWallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByWalletID", ctx, walletID)
	ret0, _ := ret[0].(*repository.DBWallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByWalletID indicates an expected call of GetByWalletID.
func (mr *MockWalletRepositoryMockRecorder) GetByWalletID(ctx, walletID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByWalletID", reflect.TypeOf((*MockWalletRepository)(nil).GetByWalletID), ctx, walletID)
}