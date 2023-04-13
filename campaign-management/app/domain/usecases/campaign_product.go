package usecases

import (
	"context"

	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/usecases/dto"
)

//go:generate mockery --name CampaignProductUseCases --filename campaign_product_usecases.go
type CampaignProductUseCases interface {
	AddProducts(ctx context.Context, products []entities.CampaignProduct) ([]*dto.CampaignProducts, error)
	GetProducts(ctx context.Context, campaignID int64) ([]*dto.CampaignProducts, error)
	UpdateProducts(ctx context.Context, products []entities.CampaignProduct) error
	DeleteByCampaignId(ctx context.Context, campaignID int64, productID int64) error
	DeleteAllByCampaignId(ctx context.Context, campaignID int64) error
}
