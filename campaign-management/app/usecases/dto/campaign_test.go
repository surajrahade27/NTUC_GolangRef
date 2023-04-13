package dto

import (
	"campaign-mgmt/app/domain/entities"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ToCampaignDTO(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		campaignEntity := entities.Campaign{
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			CreatedBy:           12345,
			LeadTime:            20,
		}

		expectedResponse := CampaignDTO{
			ID:                  int64(1),
			Title:               "new campaign",
			Name:                "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}
		response := ToCampaignDTO(campaignEntity)
		ShouldEqual(response, expectedResponse)
	})

}

func Test_ToCampaignDataList(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		campaignEntity := []entities.Campaign{
			{
				Title:               "new campaign",
				StatusCode:          int64(1),
				CampaignType:        "deli",
				ListingTitle:        "test screen title",
				ListingDesc:         "test description",
				ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
				OnboardTitle:        "test campaign",
				OnboardDesc:         "test desc",
				OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
				LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
				OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
				OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
				CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
				CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
				OfferID:             123,
				TagID:               456,
				CreatedBy:           12345,
				LeadTime:            20,
			},
		}

		paginationData := entities.PaginationConfig{
			Limit:  20,
			Page:   1,
			Sort:   "asc",
			Name:   "campaign",
			Status: 1,
		}

		expectedResponse := DataList{
			PaginationFields: PaginationFields{
				Count:  1,
				Limit:  20,
				Offset: 1,
			},
			Campaigns: []CampaignDTO{
				{
					ID:                  int64(1),
					Title:               "new campaign",
					Name:                "new campaign",
					StatusCode:          1,
					CampaignType:        "deli",
					ListingTitle:        "test screen title",
					ListingDesc:         "test description",
					ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
					OnboardTitle:        "test campaign",
					OnboardDesc:         "test desc",
					OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
					LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
					OrderStartDate:      "2023-03-01 12:00:00",
					OrderEndDate:        "2023-03-31 12:00:00",
					CollectionStartDate: "2023-03-05 12:00:00",
					CollectionEndDate:   "2023-04-05 12:00:00",
					OfferID:             123,
					TagID:               456,
					LeadTime:            20,
				},
			},
		}
		response := ToCampaignDataList(campaignEntity, 1, paginationData)
		ShouldEqual(response, expectedResponse)
	})

}

func Test_FormatDate(t *testing.T) {
	t.Run("when valid date is given", func(t *testing.T) {
		date := time.Now()
		expectedDate := time.Now().Format("2006-01-02 15:04:05")

		response := formatDate(date)
		ShouldEqual(response, expectedDate)
	})

	t.Run("when empty date is given", func(t *testing.T) {
		date := time.Time{}

		response := formatDate(date)
		ShouldEqual(response, "")
	})
}

func Test_ToCampaignResponse(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		campaignDTO := CampaignDTO{
			ID:                  int64(1),
			Title:               "new campaign",
			Name:                "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedResponse := CampaignResponse{
			ListResponseFields: ListResponseFields{
				Code:   200,
				Status: "SUCCESS",
			},
			Data: CampaignDTO{
				ID:                  int64(1),
				Title:               "new campaign",
				Name:                "new campaign",
				StatusCode:          1,
				CampaignType:        "deli",
				ListingTitle:        "test screen title",
				ListingDesc:         "test description",
				ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
				OnboardTitle:        "test campaign",
				OnboardDesc:         "test desc",
				OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
				LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
				OrderStartDate:      "2023-03-01 12:00:00",
				OrderEndDate:        "2023-03-31 12:00:00",
				CollectionStartDate: "2023-03-05 12:00:00",
				CollectionEndDate:   "2023-04-05 12:00:00",
				OfferID:             123,
				TagID:               456,
				LeadTime:            20,
			},
		}
		response := ToCampaignResponse(campaignDTO)
		ShouldEqual(response, expectedResponse)
	})
}

func Test_ToCampaignListResponse(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		dataList := DataList{
			PaginationFields: PaginationFields{
				Count:  1,
				Limit:  20,
				Offset: 1,
			},
			Campaigns: []CampaignDTO{
				{
					ID:                  int64(1),
					Title:               "new campaign",
					Name:                "new campaign",
					StatusCode:          1,
					CampaignType:        "deli",
					ListingTitle:        "test screen title",
					ListingDesc:         "test description",
					ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
					OnboardTitle:        "test campaign",
					OnboardDesc:         "test desc",
					OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
					LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
					OrderStartDate:      "2023-03-01 12:00:00",
					OrderEndDate:        "2023-03-31 12:00:00",
					CollectionStartDate: "2023-03-05 12:00:00",
					CollectionEndDate:   "2023-04-05 12:00:00",
					OfferID:             123,
					TagID:               456,
					LeadTime:            20,
				},
			},
		}

		expectedResponse := CampaignListResponse{
			ListResponseFields: ListResponseFields{
				Code:   200,
				Status: "SUCCESS",
			},
			Data: dataList,
		}
		response := ToCampaignListResponse(dataList)
		ShouldEqual(response, expectedResponse)
	})
}
