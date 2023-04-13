package services

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
)

//go:generate mockery --name CampaignProducts --filename campaign_products_services.go
type CampaignProducts interface {
	CreateMultiple(ctx context.Context, products []entities.CampaignProduct) ([]entities.CampaignProduct, error)
	GetByCampaignId(ctx context.Context, CampaignID valueobjects.CampaignID) ([]entities.CampaignProduct, error)
	Update(ctx context.Context, product entities.CampaignProduct) error
	DeleteByCampaignId(ctx context.Context, campaignID int64, productID int64) error
	DeleteAllByCampaignId(ctx context.Context, campaignID int64) error
}
