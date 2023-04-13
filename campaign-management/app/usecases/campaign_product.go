package usecases

import (
	"campaign-mgmt/app/domain/valueobjects"
	"context"

	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/services"
	"campaign-mgmt/app/usecases/dto"
)

type CampaignProductUseCase struct {
	campaignProductRepo services.CampaignProducts
}

func NewCampaignProductUseCase(campaignProductRepo services.CampaignProducts) *CampaignProductUseCase {
	return &CampaignProductUseCase{
		campaignProductRepo: campaignProductRepo,
	}
}

func (c *CampaignProductUseCase) AddProducts(ctx context.Context, products []entities.CampaignProduct) ([]*dto.CampaignProducts, error) {
	response, err := c.campaignProductRepo.CreateMultiple(ctx, products)
	if err != nil {
		return nil, err
	}
	var campaignProducts []*dto.CampaignProducts
	for _, product := range response {
		campaignProducts = append(campaignProducts, dto.ToCampaignProductDTO(product))
	}
	return campaignProducts, nil
}

func (c *CampaignProductUseCase) GetProducts(ctx context.Context, campaignID int64) ([]*dto.CampaignProducts, error) {
	campaignProductsEntities, err := c.campaignProductRepo.GetByCampaignId(ctx, valueobjects.CampaignID(campaignID))
	if err != nil {
		return nil, err
	}
	var campaignProducts []*dto.CampaignProducts
	for _, product := range campaignProductsEntities {
		campaignProducts = append(campaignProducts, dto.ToCampaignProductDTO(product))
	}
	return campaignProducts, nil
}

func (c *CampaignProductUseCase) UpdateProducts(ctx context.Context, campaignProductDetails []entities.CampaignProduct) error {
	for _, product := range campaignProductDetails {
		err := c.campaignProductRepo.Update(ctx, product)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CampaignProductUseCase) DeleteByCampaignId(ctx context.Context, campaignID int64, productID int64) error {
	err := c.campaignProductRepo.DeleteByCampaignId(ctx, campaignID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CampaignProductUseCase) DeleteAllByCampaignId(ctx context.Context, campaignID int64) error {
	err := c.campaignProductRepo.DeleteAllByCampaignId(ctx, campaignID)
	if err != nil {
		return err
	}
	return nil
}
