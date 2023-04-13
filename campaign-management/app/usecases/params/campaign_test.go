package params

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ToCampaignEntity(t *testing.T) {
	t.Run("test conversion : success scenario", func(t *testing.T) {
		request := CampaignCreationForm{
			Title:               "new campaign",
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

		expectedResponse := entities.Campaign{
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
			LeadTime:            20,
		}

		response, err := ToCampaignEntity(request)
		ShouldBeNil(err)
		ShouldEqual(response, expectedResponse)
	})

	t.Run("test conversion : failure scenario - order start date", func(t *testing.T) {
		request := CampaignCreationForm{
			Title:               "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      "2023-03-01 12:99:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-10 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-03-01 12:99:00": minute out of range`
		_, err := ToCampaignEntity(request)
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("test conversion : failure scenario - order end date", func(t *testing.T) {
		request := CampaignCreationForm{
			Title:               "new campaign",
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
			OrderEndDate:        "2023-03-31 33:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-10 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-03-31 33:00:00": hour out of range`
		_, err := ToCampaignEntity(request)
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("test conversion : failure scenario - collection start date", func(t *testing.T) {
		request := CampaignCreationForm{
			Title:               "new campaign",
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
			CollectionStartDate: "2023-03-05 00",
			CollectionEndDate:   "2023-04-10 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-03-05 00" as "2006-01-02 15:04:05": cannot parse "" as ":"`
		_, err := ToCampaignEntity(request)
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("test conversion : failure scenario - collection end date", func(t *testing.T) {
		request := CampaignCreationForm{
			Title:               "new campaign",
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
			CollectionEndDate:   "2023-04-10 12:00:78",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-04-10 12:00:78": second out of range`
		_, err := ToCampaignEntity(request)
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
}

func Test_ToUpdateCampaignEntity(t *testing.T) {
	t.Run("test conversion : success scenario", func(t *testing.T) {
		request := CampaignUpdateForm{
			Title:               "new campaign",
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

		expectedResponse := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
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
			LeadTime:            20,
		}

		response, err := ToUpdateCampaignEntity(request, int64(1))
		ShouldBeNil(err)
		ShouldEqual(response, expectedResponse)
	})

	t.Run("test conversion : failure scenario - order start date", func(t *testing.T) {
		request := CampaignUpdateForm{
			Title:               "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      "2023-03-01 12:99:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-10 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-03-01 12:99:00": minute out of range`
		_, err := ToUpdateCampaignEntity(request, int64(1))
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("test conversion : failure scenario - order end date", func(t *testing.T) {
		request := CampaignUpdateForm{
			Title:               "new campaign",
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
			OrderEndDate:        "2023-03-31 33:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-10 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-03-31 33:00:00": hour out of range`
		_, err := ToUpdateCampaignEntity(request, int64(1))
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("test conversion : failure scenario - collection start date", func(t *testing.T) {
		request := CampaignUpdateForm{
			Title:               "new campaign",
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
			CollectionStartDate: "2023-03-05 00",
			CollectionEndDate:   "2023-04-10 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-03-05 00" as "2006-01-02 15:04:05": cannot parse "" as ":"`
		_, err := ToUpdateCampaignEntity(request, int64(1))
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("test conversion : failure scenario - collection end date", func(t *testing.T) {
		request := CampaignUpdateForm{
			Title:               "new campaign",
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
			CollectionEndDate:   "2023-04-10 12:00:78",
			OfferID:             123,
			TagID:               456,
			LeadTime:            20,
		}

		expectedErr := `parsing time "2023-04-10 12:00:78": second out of range`
		_, err := ToUpdateCampaignEntity(request, int64(1))
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
}
