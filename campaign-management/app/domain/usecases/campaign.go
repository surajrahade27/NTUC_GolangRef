package usecases

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/usecases/dto"
	"context"
)

//go:generate mockery --name CampaignUseCases --filename campaign_usecases.go
type CampaignUseCases interface {
	Get(ctx context.Context, campaignID int64) (*dto.CampaignDTO, error)
	Create(ctx context.Context, campaignData entities.Campaign) (*dto.CampaignDTO, error)
	Exists(ctx context.Context, campaignID int64, title string) (bool, error)
	Update(ctx context.Context, campaignData entities.Campaign) error
	UpdateStatus(ctx context.Context) error
	GetList(ctx context.Context, paginationData entities.PaginationConfig) (*dto.CampaignListResponse, error)
}
