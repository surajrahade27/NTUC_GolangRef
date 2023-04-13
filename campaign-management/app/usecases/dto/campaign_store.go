package dto

import "campaign-mgmt/app/domain/entities"

type CampaignStoresDTO struct {
	CampaignID int64             `json:"campaign_id"`
	Stores     []*CampaignStores `json:"stores"`
}

type CampaignStores struct {
	ID      int64 `json:"campaign_store_id"`
	StoreID int64 `json:"store_id"`
}

func ToCampaignStoreDTO(storeEntity entities.CampaignStore) *CampaignStores {
	return &CampaignStores{
		ID:      storeEntity.ID.ToInt64(),
		StoreID: storeEntity.StoreID,
	}
}
