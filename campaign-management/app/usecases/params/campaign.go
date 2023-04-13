package params

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/util"
	"time"
)

// CampaignCreationForm ..
// swagger:model CampaignCreationForm
type CampaignCreationForm struct {
	// Campaign Title
	Title string `json:"title" validate:"required"`
	// Campaign Status code
	StatusCode int `json:"campaign_status_code" validate:"oneof=1 2 3"`
	// Campaign Type
	CampaignType string `json:"campaign_type"  validate:"omitempty,oneof=deli cash&carry"`
	// Listing screen title
	ListingTitle string `json:"listing_title"`
	// Listing screen description
	ListingDesc string `json:"listing_description"`
	// Listing screen image path
	ListingImagePath string `json:"listing_image_path" validate:"omitempty,url"`
	// Onboarding title
	OnboardTitle string `json:"onboarding_title"`
	// Onboarding Description
	OnboardDesc string `json:"onboarding_description"`
	// Onboarding image path
	OnboardImagePath string `json:"onboarding_image_path" validate:"omitempty,url"`
	// Landing screen image path
	LandingImagePath string `json:"landing_image_path"  validate:"omitempty,url"`
	// Order start date
	OrderStartDate string `json:"order_start_date" example:"2023-12-31 12:00:00"`
	// Order end date
	OrderEndDate string `json:"order_end_date" example:"2023-12-31 12:00:00"`
	// Collection start date
	CollectionStartDate string `json:"collection_start_date" example:"2023-12-31 12:00:00"`
	// Collection end date
	CollectionEndDate string `json:"collection_end_date" example:"2023-12-31 12:00:00"`
	// Lead time in days
	LeadTime int `json:"lead_time"`
	// Offer Id
	OfferID int64 `json:"offer_id"`
	// Tag Id
	TagID int64 `json:"tag_id"`
	// Is campaign published flag
	IsCampaignPublished bool `json:"is_campaign_published"`
	// List of campaign stores
	Stores []int64 `json:"stores"`
	// List of campaign products
	Products []CampaignProduct `json:"products"`
}

// CampaignUpdateForm ..
// swagger:model CampaignUpdateForm
type CampaignUpdateForm struct {
	// Campaign Title
	Title string `json:"title"`
	// Campaign Status code
	StatusCode int `json:"campaign_status_code" validate:"oneof=1 2 3"`
	// Campaign Type
	CampaignType string `json:"campaign_type" validate:"omitempty,oneof=deli cash&carry"`
	// Listing screen title
	ListingTitle string `json:"listing_title"`
	// Listing screen description
	ListingDesc string `json:"listing_description"`
	// Listing screen image path
	ListingImagePath string `json:"listing_image_path" validate:"omitempty,url"`
	// Onboarding title
	OnboardTitle string `json:"onboarding_title"`
	// Onboarding Description
	OnboardDesc string `json:"onboarding_description"`
	// Onboarding image path
	OnboardImagePath string `json:"onboarding_image_path" validate:"omitempty,url"`
	// Landing screen image path
	LandingImagePath string `json:"landing_image_path"  validate:"omitempty,url"`
	// Order start date
	OrderStartDate string `json:"order_start_date" example:"2023-12-31 12:00:00"`
	// Order end date
	OrderEndDate string `json:"order_end_date" example:"2023-12-31 12:00:00"`
	// Collection start date
	CollectionStartDate string `json:"collection_start_date" example:"2023-12-31 12:00:00"`
	// Collection end date
	CollectionEndDate string `json:"collection_end_date" example:"2023-12-31 12:00:00"`
	// Lead time in days
	LeadTime int `json:"lead_time"`
	// Offer Id
	OfferID int64 `json:"offer_id"`
	// Tag Id
	TagID int64 `json:"tag_id"`
	// Is campaign published flag
	IsCampaignPublished bool `json:"is_campaign_published"`
	// List of campaign stores
	Stores []int64 `json:"stores"`
	// List of campaign products
	Products []UpdateCampaignProduct `json:"products" validate:"dive"`
}

type CampaignDates struct {
	OrderStartDate      string
	OrderEndDate        string
	CollectionStartDate string
	CollectionEndDate   string
}

func ToCampaignEntity(campaign CampaignCreationForm) (entities.Campaign, error) {
	var orderStartDate, orderEndDate, collectionStartDate, collectionEndDate time.Time
	var err error
	orderStartDate, err = util.ToDateTime(campaign.OrderStartDate)
	if err != nil {
		return entities.Campaign{}, err
	}
	orderEndDate, err = util.ToDateTime(campaign.OrderEndDate)
	if err != nil {
		return entities.Campaign{}, err
	}
	collectionStartDate, err = util.ToDateTime(campaign.CollectionStartDate)
	if err != nil {
		return entities.Campaign{}, err
	}
	collectionEndDate, err = util.ToDateTime(campaign.CollectionEndDate)
	if err != nil {
		return entities.Campaign{}, err
	}

	return entities.Campaign{
		Title:               campaign.Title,
		StatusCode:          int64(campaign.StatusCode),
		CampaignType:        valueobjects.CampaignType(campaign.CampaignType),
		ListingTitle:        campaign.ListingTitle,
		ListingDesc:         campaign.ListingDesc,
		ListingImagePath:    campaign.ListingImagePath,
		OnboardTitle:        campaign.OnboardTitle,
		OnboardDesc:         campaign.OnboardDesc,
		OnboardImagePath:    campaign.OnboardImagePath,
		LandingImagePath:    campaign.LandingImagePath,
		OrderStartDate:      orderStartDate,
		OrderEndDate:        orderEndDate,
		CollectionStartDate: collectionStartDate,
		CollectionEndDate:   collectionEndDate,
		OfferID:             campaign.OfferID,
		TagID:               campaign.TagID,
		LeadTime:            campaign.LeadTime,
		IsCampaignPublished: campaign.IsCampaignPublished,
	}, nil
}

func ToUpdateCampaignEntity(campaign CampaignUpdateForm, campaignID int64) (entities.Campaign, error) {
	var orderStartDate, orderEndDate, collectionStartDate, collectionEndDate time.Time
	var err error
	if campaign.OrderStartDate != "" {
		orderStartDate, err = util.ToDateTime(campaign.OrderStartDate)
		if err != nil {
			return entities.Campaign{}, err
		}
	}
	if campaign.OrderEndDate != "" {
		orderEndDate, err = util.ToDateTime(campaign.OrderEndDate)
		if err != nil {
			return entities.Campaign{}, err
		}
	}
	if campaign.CollectionStartDate != "" {
		collectionStartDate, err = util.ToDateTime(campaign.CollectionStartDate)
		if err != nil {
			return entities.Campaign{}, err
		}
	}
	if campaign.CollectionEndDate != "" {
		collectionEndDate, err = util.ToDateTime(campaign.CollectionEndDate)
		if err != nil {
			return entities.Campaign{}, err
		}
	}

	return entities.Campaign{
		ID:                  valueobjects.CampaignID(campaignID),
		Title:               campaign.Title,
		StatusCode:          int64(campaign.StatusCode),
		CampaignType:        valueobjects.CampaignType(campaign.CampaignType),
		ListingTitle:        campaign.ListingTitle,
		ListingDesc:         campaign.ListingDesc,
		ListingImagePath:    campaign.ListingImagePath,
		OnboardTitle:        campaign.OnboardTitle,
		OnboardDesc:         campaign.OnboardDesc,
		OnboardImagePath:    campaign.OnboardImagePath,
		LandingImagePath:    campaign.LandingImagePath,
		OrderStartDate:      orderStartDate,
		OrderEndDate:        orderEndDate,
		CollectionStartDate: collectionStartDate,
		CollectionEndDate:   collectionEndDate,
		OfferID:             campaign.OfferID,
		TagID:               campaign.TagID,
		LeadTime:            campaign.LeadTime,
		IsCampaignPublished: campaign.IsCampaignPublished,
	}, nil
}
