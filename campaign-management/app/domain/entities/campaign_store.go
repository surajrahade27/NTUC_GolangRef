package entities

import (
	"campaign-mgmt/app/domain/valueobjects"
	"time"
)

type CampaignStore struct {
	ID         valueobjects.CampaignStoreID
	CampaignID valueobjects.CampaignID
	StoreID    int64
	CreatedAt  time.Time
	CreatedBy  int64
	UpdatedAt  time.Time
	UpdatedBy  int64
	DeletedAt  time.Time
	DeletedBy  int64
}
