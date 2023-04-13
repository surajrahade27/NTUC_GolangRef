// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	entities "campaign-mgmt/app/domain/entities"
	dto "campaign-mgmt/app/usecases/dto"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CampaignProductUseCases is an autogenerated mock type for the CampaignProductUseCases type
type CampaignProductUseCases struct {
	mock.Mock
}

// AddProducts provides a mock function with given fields: ctx, products
func (_m *CampaignProductUseCases) AddProducts(ctx context.Context, products []entities.CampaignProduct) ([]*dto.CampaignProducts, error) {
	ret := _m.Called(ctx, products)

	var r0 []*dto.CampaignProducts
	if rf, ok := ret.Get(0).(func(context.Context, []entities.CampaignProduct) []*dto.CampaignProducts); ok {
		r0 = rf(ctx, products)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dto.CampaignProducts)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []entities.CampaignProduct) error); ok {
		r1 = rf(ctx, products)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAllByCampaignId provides a mock function with given fields: ctx, campaignID
func (_m *CampaignProductUseCases) DeleteAllByCampaignId(ctx context.Context, campaignID int64) error {
	ret := _m.Called(ctx, campaignID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, campaignID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByCampaignId provides a mock function with given fields: ctx, campaignID, productID
func (_m *CampaignProductUseCases) DeleteByCampaignId(ctx context.Context, campaignID int64, productID int64) error {
	ret := _m.Called(ctx, campaignID, productID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, campaignID, productID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProducts provides a mock function with given fields: ctx, campaignID
func (_m *CampaignProductUseCases) GetProducts(ctx context.Context, campaignID int64) ([]*dto.CampaignProducts, error) {
	ret := _m.Called(ctx, campaignID)

	var r0 []*dto.CampaignProducts
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*dto.CampaignProducts); ok {
		r0 = rf(ctx, campaignID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dto.CampaignProducts)
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

// UpdateProducts provides a mock function with given fields: ctx, products
func (_m *CampaignProductUseCases) UpdateProducts(ctx context.Context, products []entities.CampaignProduct) error {
	ret := _m.Called(ctx, products)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []entities.CampaignProduct) error); ok {
		r0 = rf(ctx, products)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCampaignProductUseCases interface {
	mock.TestingT
	Cleanup(func())
}

// NewCampaignProductUseCases creates a new instance of CampaignProductUseCases. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCampaignProductUseCases(t mockConstructorTestingTNewCampaignProductUseCases) *CampaignProductUseCases {
	mock := &CampaignProductUseCases{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}