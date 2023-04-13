package usecases

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/services/mocks"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/dto"
	"campaign-mgmt/app/usecases/util"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCampaignUseCase_ExistsOtherWay(t *testing.T) {
	campaignService := mocks.NewCampaigns(t)
	campaignUseCase := NewCampaignUseCase(campaignService)

	Convey("Given a campaign has(exists) use case", t, func() {
		ctx := context.Background()

		Convey("When campaign exists, it returns true", func() {
			campaignID := valueobjects.CampaignID(1)
			campaignTitle := ""
			campaignService.On("Exists", ctx, campaignID, campaignTitle).Return(
				true,
				nil,
			)
			actualValue, err := campaignUseCase.Exists(ctx, int64(1), campaignTitle)
			ShouldEqual(actualValue, true)
			ShouldBeNil(err)
		})

		Convey("When campaign does not exist, it returns false", func() {
			campaignID := valueobjects.CampaignID(2)
			campaignTitle := ""
			campaignService.On("Exists", ctx, campaignID, campaignTitle).Return(
				false,
				errors.New("something happened"),
			)
			actualValue, err := campaignUseCase.Exists(ctx, int64(2), campaignTitle)
			ShouldEqual(actualValue, false)
			ShouldNotBeNil(err)
		})

	})
}

func TestCampaignUseCase_Exists(t *testing.T) {
	t.Run("When campaign exists, it returns true", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignID := valueobjects.CampaignID(1)
		campaignTitle := ""
		campaignService.On("Exists", ctx, campaignID, campaignTitle).Return(
			true,
			nil,
		)
		actualValue, err := campaignUseCase.Exists(ctx, int64(1), campaignTitle)
		ShouldEqual(actualValue, true)
		ShouldBeNil(err)
	})

	t.Run("When campaign does not exist, it returns false", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignID := valueobjects.CampaignID(1)
		campaignTitle := ""
		campaignService.On("Exists", ctx, campaignID, campaignTitle).Return(
			true,
			nil,
		)
		actualValue, err := campaignUseCase.Exists(ctx, int64(1), campaignTitle)
		ShouldEqual(actualValue, true)
		ShouldBeNil(err)
	})

	t.Run("When some error occured", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignID := valueobjects.CampaignID(1)
		campaignTitle := ""
		campaignService.On("Exists", ctx, campaignID, campaignTitle).Return(
			true,
			valueobjects.ErrCampaignCantExist, errors.New("db error"),
		)
		actualValue, err := campaignUseCase.Exists(ctx, int64(1), campaignTitle)
		ShouldEqual(actualValue, true)
		ShouldBeNil(err)
	})
}

func TestCampaignUseCase_ExistsOtherWay1(t *testing.T) {
	campaignUseCase := NewCampaignUseCase(nil)
	ctx := context.Background()
	tests := []struct {
		name          string
		prepare       func()
		expectedErr   error
		campaignID    int64
		campaignTitle string
		expectedValue bool
	}{
		{
			name: "when the campaign exist",
			prepare: func() {
				campaignService := mocks.NewCampaigns(t)
				campaignUseCase = NewCampaignUseCase(campaignService)
				campaignService.On("Exists", ctx, valueobjects.CampaignID(1), "").Return(
					true,
					nil,
				)
			},
			expectedErr:   nil,
			expectedValue: true,
			campaignID:    1,
			campaignTitle: "",
		},

		{
			name: "when the campaign no exist",
			prepare: func() {
				campaignService := mocks.NewCampaigns(t)
				campaignUseCase = NewCampaignUseCase(campaignService)
				campaignService.On("Exists", ctx, valueobjects.CampaignID(2), "").Return(
					false,
					errors.New("something happenend"),
				)
			},
			expectedErr:   errors.New("something happenend"),
			expectedValue: false,
			campaignID:    2,
			campaignTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			actualValue, actualErr := campaignUseCase.Exists(ctx, tt.campaignID, tt.campaignTitle)
			ShouldEqual(actualValue, tt.expectedValue)
			ShouldEqual(actualErr, tt.expectedErr)
		})
	}
}

func TestCampaignUseCase_Get(t *testing.T) {
	t.Run("When campaign details exist, it returns campaign Details", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignID := valueobjects.CampaignID(1)
		response := entities.Campaign{
			ID:                  campaignID,
			Title:               "Test Campaign 1",
			CampaignType:        "deli",
			ListingTitle:        "listing screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OnboardTitle:        "New Year:2023",
			OnboardDesc:         "this is sample campaign 1",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			CreatedBy:           1,
			OrderStartDate:      time.Date(2000, 1, 3, 12, 30, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2000, 1, 31, 12, 30, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2000, 1, 5, 12, 30, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2000, 2, 2, 12, 30, 0, 0, time.UTC),
			LeadTime:            10,
			OfferID:             1234,
		}
		campaignService.On("Get", ctx, campaignID).Return(
			response,
			nil,
		)
		actualValue, err := campaignUseCase.Get(ctx, int64(1))
		ShouldEqual(actualValue, response)
		ShouldBeNil(err)
	})
	t.Run("When campaign details not exist, it returns error", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignID := valueobjects.CampaignID(1000)
		response := entities.Campaign{}
		campaignService.On("Get", ctx, campaignID).Return(
			response,
			errors.New("record Not Found"),
		)
		actualValue, err := campaignUseCase.Get(ctx, int64(1000))
		ShouldEqual(actualValue, response)
		ShouldNotBeNil(err)
	})
}

func TestCampaignUseCase_GetList(t *testing.T) {
	t.Run("When campaign details exist, it returns campaigns list", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignDetails1 := entities.Campaign{
			ID:                  1,
			StatusCode:          int64(1),
			Title:               "Test Campaign 1",
			CampaignType:        "deli",
			ListingTitle:        "listing screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OnboardTitle:        "New Year:2023",
			OnboardDesc:         "this is sample campaign 1",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			CreatedBy:           int64(12),
			OrderStartDate:      time.Date(2000, 1, 3, 12, 30, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2000, 1, 31, 12, 30, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2000, 1, 5, 12, 30, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2000, 2, 2, 12, 30, 0, 0, time.UTC),
			LeadTime:            10,
			OfferID:             1234,
		}
		campaignDetails2 := entities.Campaign{
			ID:                  2,
			StatusCode:          int64(1),
			Title:               "Test Campaign 2",
			CampaignType:        "deli",
			ListingTitle:        "listing screen title 2",
			ListingDesc:         "test description ",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OnboardTitle:        "New Year:2023",
			OnboardDesc:         "this is sample campaign 1",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			CreatedBy:           int64(12),
			OrderStartDate:      time.Date(2001, 1, 3, 12, 30, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2001, 1, 31, 12, 30, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2001, 1, 5, 12, 30, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2001, 2, 2, 12, 30, 0, 0, time.UTC),
			LeadTime:            10,
			OfferID:             12345,
		}
		campaignEntities := []entities.Campaign{
			campaignDetails1,
			campaignDetails2,
		}
		responseDTO := dto.CampaignListResponse{
			ListResponseFields: dto.ListResponseFields{Code: 200, Status: "SUCCESS"},
			Data:               dto.ToCampaignDataList(campaignEntities, 2, entities.PaginationConfig{Limit: 20, Page: 1}),
		}
		campaignService.On("GetList", ctx, entities.PaginationConfig{Limit: 20, Page: 1}).Return(
			campaignEntities,
			int64(2),
			nil,
		)
		actualValue, err := campaignUseCase.GetList(ctx, entities.PaginationConfig{Limit: 20, Page: 1})
		ShouldEqual(actualValue, responseDTO)
		ShouldBeNil(err)
	})

	t.Run("When campaign details does not exist, it returns error", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		var response []entities.Campaign
		campaignService.On("GetList", ctx, entities.PaginationConfig{Limit: 20, Page: 1}).Return(
			response, int64(0),
			errors.New("records not found"),
		)
		actualValue, err := campaignUseCase.GetList(ctx, entities.PaginationConfig{Limit: 20, Page: 1})
		ShouldEqual(actualValue, response)
		ShouldNotBeNil(err)
	})
}

func TestCampaignUseCase_Create(t *testing.T) {
	dateStr := time.Now().Format("2006-01-02 15:04:05")
	dateTime, _ := util.ToDateTime(dateStr)
	campaignEntity := entities.Campaign{
		Title:               "test_campaign",
		StatusCode:          int64(1),
		CampaignType:        "deli",
		ListingTitle:        "listing screen title",
		ListingDesc:         "test description",
		ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
		OnboardTitle:        "New Year:2023",
		OnboardDesc:         "this is sample campaign 1",
		OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
		LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
		OrderStartDate:      dateTime,
		OrderEndDate:        dateTime,
		CollectionStartDate: dateTime,
		CollectionEndDate:   dateTime,
		OfferID:             int64(123),
		TagID:               int64(456),
		CreatedBy:           int64(12121212),
		LeadTime:            10,
	}
	t.Run("when campaign creation is successful", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignDetails := dto.CampaignDTO{
			ID:                  1,
			Title:               "test_campaign",
			Name:                "test_campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "listing screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OnboardTitle:        "New Year:2023",
			OnboardDesc:         "this is sample campaign 1",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      dateStr,
			OrderEndDate:        dateStr,
			CollectionStartDate: dateStr,
			CollectionEndDate:   dateStr,
			LeadTime:            10,
			OfferID:             int64(123),
			TagID:               int64(456),
		}
		campaignService.On("Create", ctx, campaignEntity).Return(
			entities.Campaign{
				ID:                  valueobjects.CampaignID(1),
				Title:               "test_campaign",
				StatusCode:          int64(1),
				CampaignType:        "deli",
				ListingTitle:        "listing screen title",
				ListingDesc:         "test description",
				ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
				OnboardTitle:        "New Year:2023",
				OnboardDesc:         "this is sample campaign 1",
				OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
				LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
				OrderStartDate:      dateTime,
				OrderEndDate:        dateTime,
				CollectionStartDate: dateTime,
				CollectionEndDate:   dateTime,
				OfferID:             int64(123),
				TagID:               int64(456),
			}, nil)
		response, err := campaignUseCase.Create(ctx, campaignEntity)
		ShouldEqual(response, campaignDetails)
		ShouldBeNil(err)
	})
	t.Run("when error occured while campaign creation", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignService.On("Create", ctx, campaignEntity).Return(
			entities.Campaign{}, fmt.Errorf("%w: %v", valueobjects.ErrCampaignCantCreate, errors.New("db error")))
		_, err := campaignUseCase.Create(ctx, campaignEntity)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrCampaignCantCreate) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignUseCase_Update(t *testing.T) {
	dateStr := time.Now().Format("2006-01-02 15:04:05")
	dateTime, _ := util.ToDateTime(dateStr)
	campaignEntity := entities.Campaign{
		ID:                  valueobjects.CampaignID(1),
		Title:               "test_campaign",
		StatusCode:          int64(1),
		CampaignType:        "deli",
		ListingTitle:        "listing screen title",
		ListingDesc:         "test description",
		ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
		OnboardTitle:        "New Year:2023",
		OnboardDesc:         "this is sample campaign 1",
		OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
		LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
		OrderStartDate:      dateTime,
		OrderEndDate:        dateTime,
		CollectionStartDate: dateTime,
		CollectionEndDate:   dateTime,
		OfferID:             int64(123),
		TagID:               int64(456),
		UpdatedBy:           int64(12121212),
	}
	t.Run("when campaign update is successful", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)

		campaignService.On("Update", ctx, campaignEntity).Return(nil)
		err := campaignUseCase.Update(ctx, campaignEntity)
		ShouldBeNil(err)
	})
	t.Run("when error occured while updating campaign details", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignService.On("Update", ctx, campaignEntity).Return(fmt.Errorf("%w: %v",
			valueobjects.ErrCampaignCantUpdate, errors.New("db error")))
		err := campaignUseCase.Update(ctx, campaignEntity)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrCampaignCantUpdate) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignUseCase_UpdateStatus(t *testing.T) {
	t.Run("when campaign updated successfully", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)

		campaignService.On("UpdateStatus", ctx).Return(nil)
		err := campaignUseCase.UpdateStatus(ctx)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
	t.Run("when error occured while updating campaign  status", func(t *testing.T) {
		ctx := context.Background()
		campaignService := mocks.NewCampaigns(t)
		campaignUseCase := NewCampaignUseCase(campaignService)
		campaignService.On("UpdateStatus", ctx).Return(fmt.Errorf("%w: %v", valueobjects.ErrCampaignStatusCantUpdate, errors.New("db error")))
		err := campaignUseCase.UpdateStatus(ctx)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrCampaignStatusCantUpdate) {
			t.Error("invalid error type")
		}
	})
}
