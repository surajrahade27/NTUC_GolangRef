package dto

import (
	"campaign-mgmt/app/domain/entities"
	"net/http"
	"time"
)

// CampaignDTO ..
// swagger:response CampaignDTO
type CampaignDTO struct {
	// Campaign identifier
	ID int64 `json:"id"`
	// Campaign Title
	Title string `json:"campaign_title"`
	// Campaign Name
	Name string `json:"name"`
	// Campaign Status
	StatusCode int `json:"campaign_status_code"`
	// Campaign Type
	CampaignType string `json:"campaign_type"`
	// Campaign listing screen title
	ListingTitle string `json:"listing_title"`
	// Campaign listing screen description
	ListingDesc string `json:"listing_description"`
	// Campaign listing screen image path
	ListingImagePath string `json:"listing_image_path"`
	// Campaign onboarding screen title
	OnboardTitle string `json:"onboarding_title"`
	// Campaign obboarding screen description
	OnboardDesc string `json:"onboarding_description"`
	// Campaign onboarding screen image path
	OnboardImagePath string `json:"onboard_image_path"`
	// Campaign landing screen image path
	LandingImagePath string `json:"landing_image_path"`
	// Campaign order start date
	OrderStartDate string `json:"order_start_date"`
	// Campaign order end date
	OrderEndDate string `json:"order_end_date"`
	// Campaign collection start date
	CollectionStartDate string `json:"collection_start_date"`
	// Campaign collection end date
	CollectionEndDate string `json:"collection_end_date"`
	// Campaign lead time in days
	LeadTime int `json:"lead_time"`
	// Offer Identifier
	OfferID int64 `json:"offer_id"`
	// Tag Identifier
	TagID int64 `json:"tag_id"`
	// Is campaign published flag
	IsCampaignPublished bool `json:"is_campaign_published"`
	// Product Details.
	CampaignProducts []*CampaignProducts `json:"campaign_products,omitempty"`
	// Stores Details.
	CampaignStores []*CampaignStores `json:"campaign_stores,omitempty"`
}

func formatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	} else {
		return date.Format("2006-01-02 15:04:05")
	}
}

func ToCampaignDTO(campaignEntity entities.Campaign) CampaignDTO {
	return CampaignDTO{
		ID:                  campaignEntity.ID.ToInt64(),
		Title:               campaignEntity.Title,
		Name:                campaignEntity.Title,
		StatusCode:          int(campaignEntity.StatusCode),
		CampaignType:        string(campaignEntity.CampaignType),
		ListingTitle:        campaignEntity.ListingTitle,
		ListingDesc:         campaignEntity.ListingDesc,
		ListingImagePath:    campaignEntity.ListingImagePath,
		OnboardTitle:        campaignEntity.OnboardTitle,
		OnboardDesc:         campaignEntity.OnboardDesc,
		OnboardImagePath:    campaignEntity.OnboardImagePath,
		LandingImagePath:    campaignEntity.LandingImagePath,
		OrderStartDate:      formatDate(campaignEntity.OrderStartDate),
		OrderEndDate:        formatDate(campaignEntity.OrderEndDate),
		CollectionStartDate: formatDate(campaignEntity.CollectionStartDate),
		CollectionEndDate:   formatDate(campaignEntity.CollectionEndDate),
		LeadTime:            campaignEntity.LeadTime,
		OfferID:             campaignEntity.OfferID,
		TagID:               campaignEntity.TagID,
		IsCampaignPublished: campaignEntity.IsCampaignPublished,
	}
}

type CampaignListResponse struct {
	ListResponseFields
	Data DataList `json:"data"`
}

type CampaignResponse struct {
	ListResponseFields
	Data CampaignDTO `json:"data"`
}

type DataList struct {
	PaginationFields
	Campaigns []CampaignDTO `json:"campaigns"`
}

func ToCampaignDataList(entries []entities.Campaign, count int64, paginationData entities.PaginationConfig) DataList {
	var campaigns = make([]CampaignDTO, 0)
	for _, entry := range entries {
		campaign := ToCampaignDTO(entry)
		campaigns = append(campaigns, campaign)
	}

	return DataList{
		PaginationFields{Count: count, Limit: paginationData.Limit, Offset: 0},
		campaigns,
	}
}
func ToCampaignListResponse(dataList DataList) CampaignListResponse {
	return CampaignListResponse{
		ListResponseFields{http.StatusOK, "SUCCESS"},
		dataList,
	}
}

func ToCampaignResponse(Campaign CampaignDTO) CampaignResponse {
	return CampaignResponse{
		ListResponseFields: ListResponseFields{
			Code:   http.StatusOK,
			Status: "SUCCESS",
		},
		Data: Campaign,
	}
}
