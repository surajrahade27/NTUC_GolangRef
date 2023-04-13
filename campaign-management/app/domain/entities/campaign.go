package entities

import (
	"time"

	"campaign-mgmt/app/domain/valueobjects"
)

type Campaign struct {
	ID                  valueobjects.CampaignID
	Title               string
	OrderStartDate      time.Time
	OrderEndDate        time.Time
	CollectionStartDate time.Time
	CollectionEndDate   time.Time
	StatusCode          int64
	CampaignType        valueobjects.CampaignType
	ListingTitle        string
	ListingDesc         string
	ListingImagePath    string
	OnboardTitle        string
	OnboardDesc         string
	OnboardImagePath    string
	LandingImagePath    string
	LeadTime            int
	OfferID             int64
	TagID               int64
	IsCampaignPublished bool
	CreatedAt           time.Time
	CreatedBy           int64
	UpdatedAt           time.Time
	UpdatedBy           int64
	DeletedAt           time.Time
	DeletedBy           int64
}
