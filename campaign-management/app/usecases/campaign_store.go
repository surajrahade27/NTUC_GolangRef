package usecases

import (
	"campaign-mgmt/app/domain/valueobjects"
	"context"

	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/services"
	"campaign-mgmt/app/usecases/dto"
)

type CampaignStoreUseCase struct {
	campaignStoreRepo services.CampaignStores
}

func NewCampaignStoreUseCase(campaignStoreRepo services.CampaignStores) *CampaignStoreUseCase {
	return &CampaignStoreUseCase{
		campaignStoreRepo: campaignStoreRepo,
	}
}

func (c *CampaignStoreUseCase) AddStores(ctx context.Context, stores []entities.CampaignStore) ([]*dto.CampaignStores, error) {
	response, err := c.campaignStoreRepo.CreateMultiple(ctx, stores)
	if err != nil {
		return nil, err
	}
	var campaignStores []*dto.CampaignStores
	for _, store := range response {
		campaignStores = append(campaignStores, dto.ToCampaignStoreDTO(store))
	}
	return campaignStores, nil
}

func (c *CampaignStoreUseCase) GetStores(ctx context.Context, campaignID int64) ([]*dto.CampaignStores, error) {
	campaignStoresEntities, err := c.campaignStoreRepo.GetByCampaignId(ctx, valueobjects.CampaignID(campaignID))
	if err != nil {
		return nil, err
	}
	var campaignStores []*dto.CampaignStores
	for _, store := range campaignStoresEntities {
		campaignStores = append(campaignStores, dto.ToCampaignStoreDTO(store))
	}
	return campaignStores, nil
}

func (c *CampaignStoreUseCase) UpdateStores(ctx context.Context, campaignStoreDetails []entities.CampaignStore) error {
	for _, store := range campaignStoreDetails {
		err := c.campaignStoreRepo.Update(ctx, store)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CampaignStoreUseCase) DeleteStores(ctx context.Context, campaignID, userID int64) error {
	return c.campaignStoreRepo.DeleteByCampaignID(ctx, valueobjects.CampaignID(campaignID), userID)
}

func (c *CampaignStoreUseCase) DeleteStore(ctx context.Context, campaignID, storeID, userID int64) error {
	return c.campaignStoreRepo.Delete(ctx, valueobjects.CampaignID(campaignID), valueobjects.CampaignStoreID(storeID), userID)
}

func (c *CampaignStoreUseCase) GetByStoreID(ctx context.Context, campaignID int64, StoreID int64) (dto.CampaignStores, error) {
	data, err := c.campaignStoreRepo.GetByStoreID(ctx, valueobjects.CampaignID(campaignID), StoreID)
	if err != nil {
		return dto.CampaignStores{}, err
	}
	response := dto.ToCampaignStoreDTO(data)
	return *response, nil
}

func (c *CampaignStoreUseCase) DeleteByStoreID(ctx context.Context, campaignID, storeID, userID int64) error {
	return c.campaignStoreRepo.DeleteByStoreID(ctx, valueobjects.CampaignID(campaignID), storeID, userID)
}
