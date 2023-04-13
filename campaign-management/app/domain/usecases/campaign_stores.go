package usecases

import (
	"context"

	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/usecases/dto"
)

//go:generate mockery --name CampaignStoreUseCases --filename campaign_store_usecases.go
type CampaignStoreUseCases interface {
	AddStores(ctx context.Context, stores []entities.CampaignStore) ([]*dto.CampaignStores, error)
	GetStores(ctx context.Context, campaignID int64) ([]*dto.CampaignStores, error)
	UpdateStores(ctx context.Context, stores []entities.CampaignStore) error
	DeleteStores(ctx context.Context, campaignID, userID int64) error
	DeleteStore(ctx context.Context, campaignID, campaignStoreID, userID int64) error
	GetByStoreID(ctx context.Context, campaignID int64, storeID int64) (dto.CampaignStores, error)
	DeleteByStoreID(ctx context.Context, campaignID int64, storeID, userID int64) error
}
