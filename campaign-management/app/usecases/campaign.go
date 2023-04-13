package usecases

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/services"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/dto"
	"context"
)

type CampaignUseCase struct {
	campaignRepo services.Campaigns
}

func NewCampaignUseCase(campaignRepo services.Campaigns) *CampaignUseCase {
	return &CampaignUseCase{
		campaignRepo: campaignRepo,
	}
}

func (c *CampaignUseCase) Get(ctx context.Context, campaignID int64) (*dto.CampaignDTO, error) {
	data, err := c.campaignRepo.Get(ctx, valueobjects.CampaignID(campaignID))
	if err != nil {
		return nil, err
	}
	response := dto.ToCampaignDTO(data)
	return &response, nil
}

func (c *CampaignUseCase) Create(ctx context.Context, campaignDetails entities.Campaign) (*dto.CampaignDTO, error) {
	// add logic for offer id genration
	// campaignDetails.OfferID = offerID
	campaign, err := c.campaignRepo.Create(ctx, campaignDetails)
	if err != nil {
		return nil, err
	}
	response := dto.ToCampaignDTO(campaign)
	return &response, nil
}

func (c *CampaignUseCase) Exists(ctx context.Context, campaignID int64, title string) (bool, error) {
	return c.campaignRepo.Exists(ctx, valueobjects.CampaignID(campaignID), title)
}

func (c *CampaignUseCase) Update(ctx context.Context, campaignDetails entities.Campaign) error {
	return c.campaignRepo.Update(ctx, campaignDetails)
}

func (c *CampaignUseCase) GetList(ctx context.Context, pagination entities.PaginationConfig) (*dto.CampaignListResponse, error) {
	data, count, err := c.campaignRepo.GetList(ctx, pagination)
	if err != nil {
		return nil, err
	}
	campaignList := dto.ToCampaignDataList(data, count, pagination)
	response := dto.ToCampaignListResponse(campaignList)
	return &response, nil
}

func (c *CampaignUseCase) UpdateStatus(ctx context.Context) error {
	return c.campaignRepo.UpdateStatus(ctx)
}
