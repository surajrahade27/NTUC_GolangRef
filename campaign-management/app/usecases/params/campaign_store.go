package params

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
)

// CampaignStoresForm ..
// swagger:model CampaignStoresForm
type CampaignStoresForm struct {
	// List of campaign stores
	Stores []int64 `json:"stores" validate:"required"`
}

type UpdateCampaignStore struct {
	ID      int64 `json:"campaign_store_id" validate:"required"`
	StoreID int64 `json:"store_id" validate:"required"`
}

func ToCampaignStoreEntity(storeID, campaignID, userID int64) entities.CampaignStore {
	return entities.CampaignStore{
		CampaignID: valueobjects.CampaignID(campaignID),
		StoreID:    storeID,
		CreatedBy:  userID,
	}
}

func ToUpdateCampaignStoreEntity(store UpdateCampaignStore, campaignID, userID int64) entities.CampaignStore {
	return entities.CampaignStore{
		ID:         valueobjects.CampaignStoreID(store.ID),
		CampaignID: valueobjects.CampaignID(campaignID),
		StoreID:    store.StoreID,
		UpdatedBy:  userID,
	}
}
