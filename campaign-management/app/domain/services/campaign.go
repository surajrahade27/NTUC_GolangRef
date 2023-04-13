package services

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
)

//go:generate mockery --name Campaigns --filename campaigns_services.go
type Campaigns interface {
	Get(ctx context.Context, campaignID valueobjects.CampaignID) (entities.Campaign, error)
	Create(ctx context.Context, campaignDetails entities.Campaign) (entities.Campaign, error)
	Exists(ctx context.Context, campaignID valueobjects.CampaignID, title string) (bool, error)
	Update(ctx context.Context, campaignDetails entities.Campaign) error
	GetList(ctx context.Context, paginationDetails entities.PaginationConfig) ([]entities.Campaign, int64, error)
	UpdateStatus(ctx context.Context) error
}
