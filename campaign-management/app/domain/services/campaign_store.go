package services

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
)

//go:generate mockery --name CampaignStores --filename campaign_stores_services.go
type CampaignStores interface {
	CreateMultiple(ctx context.Context, stores []entities.CampaignStore) ([]entities.CampaignStore, error)
	GetByCampaignId(ctx context.Context, campaignID valueobjects.CampaignID) ([]entities.CampaignStore, error)
	Update(ctx context.Context, store entities.CampaignStore) error
	DeleteByCampaignID(ctx context.Context, campaignID valueobjects.CampaignID, userID int64) error
	Delete(ctx context.Context, campaignID valueobjects.CampaignID, campaignStoreID valueobjects.CampaignStoreID, userID int64) error
	GetByStoreID(ctx context.Context, campaignID valueobjects.CampaignID, storeID int64) (entities.CampaignStore, error)
	DeleteByStoreID(ctx context.Context, campaignID valueobjects.CampaignID, storeID, userID int64) error
}
