package entities

import (
	"campaign-mgmt/app/domain/valueobjects"
	"time"
)

type CampaignProduct struct {
	ID          valueobjects.CampaignProductID
	CampaignID  valueobjects.CampaignID
	ProductID   int64
	SKUNo       int64
	SerialNo    int
	SequenceNo  int
	ProductType string
	CreatedAt   time.Time
	CreatedBy   int64
	UpdatedAt   time.Time
	UpdatedBy   int64
	DeletedAt   time.Time
	DeletedBy   int64
}
