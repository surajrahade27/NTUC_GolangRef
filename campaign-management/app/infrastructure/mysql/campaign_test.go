package mysql

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/util"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestCampaignService_GetCampaignsCount(t *testing.T) {
	t.Run("when campaigns count got successfully", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		gdb, _ := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		campaignService := NewCampaignService(gdb)
		list, count, err := campaignService.GetList(context.TODO(), entities.PaginationConfig{Limit: 20})
		campaignService.GetList(context.TODO(), entities.PaginationConfig{Limit: 20, Status: 2})
		campaignService.GetCampaignsCount()
		ShouldBeNil(err)
		ShouldNotBeNil(count)
		ShouldBeEmpty(list)
	})
}
func TestCampaignService_ToEntity(t *testing.T) {
	dateStr := time.Now().Format("2006-01-02 15:04:05")
	dateTime, _ := util.ToDateTime(dateStr)
	sqlDateTime := sql.NullTime{
		Time:  dateTime,
		Valid: true,
	}
	sqlOfferID := sql.NullInt64{
		Int64: 123,
		Valid: true,
	}
	sqlTagID := sql.NullInt64{
		Int64: 456,
		Valid: true,
	}
	sqlLeadTime := sql.NullInt32{
		Int32: 11,
		Valid: true,
	}
	t.Run("test conversion", func(t *testing.T) {
		db, _, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignService := NewCampaignService(gdb)
		isCampaignPublished := false

		entry := CampaignEntry{
			ID:                  1,
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
			OrderStartDate:      sqlDateTime,
			OrderEndDate:        sqlDateTime,
			CollectionStartDate: sqlDateTime,
			CollectionEndDate:   sqlDateTime,
			OfferID:             sqlOfferID,
			TagID:               sqlTagID,
			LeadTime:            sqlLeadTime,
			IsCampaignPublished: &isCampaignPublished,
		}
		expectedResponse := entities.Campaign{
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
			IsCampaignPublished: isCampaignPublished,
		}
		response := campaignService.ToEntity(entry)
		ShouldEqual(response, expectedResponse)
	})

	t.Run("when isCampaignPublished is nil", func(t *testing.T) {
		db, _, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignService := NewCampaignService(gdb)

		entry := CampaignEntry{
			ID:                  1,
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
			OrderStartDate:      sqlDateTime,
			OrderEndDate:        sqlDateTime,
			CollectionStartDate: sqlDateTime,
			CollectionEndDate:   sqlDateTime,
			OfferID:             sqlOfferID,
			TagID:               sqlTagID,
		}
		expectedResponse := entities.Campaign{
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
			IsCampaignPublished: false,
		}
		response := campaignService.ToEntity(entry)
		ShouldEqual(response, expectedResponse)
	})
}

func TestCampaignService_ToEntry(t *testing.T) {
	dateStr := time.Now().Format("2006-01-02 15:04:05")
	dateTime, _ := util.ToDateTime(dateStr)
	sqlDateTime := sql.NullTime{
		Time:  dateTime,
		Valid: true,
	}
	sqlOfferID := sql.NullInt64{
		Int64: 123,
		Valid: true,
	}
	sqlTagID := sql.NullInt64{
		Int64: 456,
		Valid: true,
	}
	sqlLeadTime := sql.NullInt32{
		Int32: 11,
		Valid: true,
	}
	t.Run("test conversion", func(t *testing.T) {
		db, _, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignService := NewCampaignService(gdb)
		isCampaignPublished := false

		entity := entities.Campaign{
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
			LeadTime:            11,
			IsCampaignPublished: isCampaignPublished,
		}
		expectedResponse := CampaignEntry{
			ID:                  1,
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
			OrderStartDate:      sqlDateTime,
			OrderEndDate:        sqlDateTime,
			CollectionStartDate: sqlDateTime,
			CollectionEndDate:   sqlDateTime,
			OfferID:             sqlOfferID,
			TagID:               sqlTagID,
			LeadTime:            sqlLeadTime,
			IsCampaignPublished: &isCampaignPublished,
		}
		response := campaignService.ToEntry(entity)
		ShouldEqual(response, expectedResponse)
	})
}

func TestCampaignService_publishCampaign(t *testing.T) {
	t.Run("success scenario", func(t *testing.T) {
		db, _, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		err = publishCampaigns(gdb)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected value : got - %v ; want - nil", err)
		}
	})
}

func TestCampaignService_deactivateCampaign(t *testing.T) {
	t.Run("success scenario", func(t *testing.T) {
		db, _, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		err = deactivateCampaigns(gdb)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected value : got - %v ; want - nil", err)
		}
	})
}

func TestCampaignService_UpdateStatus(t *testing.T) {
	t.Run("when campaigns updated successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignService := NewCampaignService(gdb)

		year, month, day := time.Now().Date()
		today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")

		const sqlUpdatetoActive = `UPDATE "campaigns" SET "status_code" = 2 WHERE CAST(order_start_date AS DATE)  <= ? and status_code = 3`
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdatetoActive)).WithArgs(today).WillReturnResult(sqlmock.NewResult(0, 1))

		const sqlUpdatetoInActive = `UPDATE "campaigns" SET "status_code" = 2 WHERE CAST(order_end_date AS DATE) < ? and status_code = 2`
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdatetoInActive)).WithArgs(today).WillReturnResult(sqlmock.NewResult(0, 1))

		err = campaignService.UpdateStatus(context.TODO())
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}

func TestCampaignService_Update(t *testing.T) {
	t.Run("when campaign with given id updated successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignService := NewCampaignService(gdb)

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

		const sqlUpdate = `UPDATE "campaigns" SET "status_code" = 2 WHERE campaign_id= ?`
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))

		err = campaignService.Update(context.TODO(), campaignEntity)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}
