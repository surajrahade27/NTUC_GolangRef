// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	entities "campaign-mgmt/app/domain/entities"
	dto "campaign-mgmt/app/usecases/dto"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CampaignStoreUseCases is an autogenerated mock type for the CampaignStoreUseCases type
type CampaignStoreUseCases struct {
	mock.Mock
}

// AddStores provides a mock function with given fields: ctx, stores
func (_m *CampaignStoreUseCases) AddStores(ctx context.Context, stores []entities.CampaignStore) ([]*dto.CampaignStores, error) {
	ret := _m.Called(ctx, stores)

	var r0 []*dto.CampaignStores
	if rf, ok := ret.Get(0).(func(context.Context, []entities.CampaignStore) []*dto.CampaignStores); ok {
		r0 = rf(ctx, stores)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dto.CampaignStores)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []entities.CampaignStore) error); ok {
		r1 = rf(ctx, stores)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByStoreID provides a mock function with given fields: ctx, campaignID, storeID, userID
func (_m *CampaignStoreUseCases) DeleteByStoreID(ctx context.Context, campaignID int64, storeID int64, userID int64) error {
	ret := _m.Called(ctx, campaignID, storeID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64) error); ok {
		r0 = rf(ctx, campaignID, storeID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteStore provides a mock function with given fields: ctx, campaignID, campaignStoreID, userID
func (_m *CampaignStoreUseCases) DeleteStore(ctx context.Context, campaignID int64, campaignStoreID int64, userID int64) error {
	ret := _m.Called(ctx, campaignID, campaignStoreID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64) error); ok {
		r0 = rf(ctx, campaignID, campaignStoreID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteStores provides a mock function with given fields: ctx, campaignID, userID
func (_m *CampaignStoreUseCases) DeleteStores(ctx context.Context, campaignID int64, userID int64) error {
	ret := _m.Called(ctx, campaignID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, campaignID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByStoreID provides a mock function with given fields: ctx, campaignID, storeID
func (_m *CampaignStoreUseCases) GetByStoreID(ctx context.Context, campaignID int64, storeID int64) (dto.CampaignStores, error) {
	ret := _m.Called(ctx, campaignID, storeID)

	var r0 dto.CampaignStores
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) dto.CampaignStores); ok {
		r0 = rf(ctx, campaignID, storeID)
	} else {
		r0 = ret.Get(0).(dto.CampaignStores)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, campaignID, storeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStores provides a mock function with given fields: ctx, campaignID
func (_m *CampaignStoreUseCases) GetStores(ctx context.Context, campaignID int64) ([]*dto.CampaignStores, error) {
	ret := _m.Called(ctx, campaignID)

	var r0 []*dto.CampaignStores
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*dto.CampaignStores); ok {
		r0 = rf(ctx, campaignID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dto.CampaignStores)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, campaignID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStores provides a mock function with given fields: ctx, stores
func (_m *CampaignStoreUseCases) UpdateStores(ctx context.Context, stores []entities.CampaignStore) error {
	ret := _m.Called(ctx, stores)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []entities.CampaignStore) error); ok {
		r0 = rf(ctx, stores)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCampaignStoreUseCases interface {
	mock.TestingT
	Cleanup(func())
}

// NewCampaignStoreUseCases creates a new instance of CampaignStoreUseCases. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCampaignStoreUseCases(t mockConstructorTestingTNewCampaignStoreUseCases) *CampaignStoreUseCases {
	mock := &CampaignStoreUseCases{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
