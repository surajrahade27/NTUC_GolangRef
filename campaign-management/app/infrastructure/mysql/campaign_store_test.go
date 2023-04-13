package mysql

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
	"regexp"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestCampaignStoreService_GetByStoreID(t *testing.T) {
	campaignID := 3
	storeID := 101

	t.Run("when store with given campaign id and store id fetched successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignStoreService := NewCampaignStoreService(gdb)

		rows := sqlmock.NewRows([]string{"campaign_store_id", "campaign_id", "store_id", "created_by", "created_at"}).
			AddRow(1, 3, 101, 9999, time.Now())

		const sqlSearch = `SELECT * FROM "campaign_stores" WHERE campaign_id = ? and store_id =?`
		mock.ExpectQuery(regexp.QuoteMeta(sqlSearch)).WithArgs(campaignID, storeID).WillReturnRows(rows)

		response, err := campaignStoreService.GetByStoreID(context.TODO(), valueobjects.CampaignID(1), int64(123))
		ShouldEqual(response, entities.CampaignStore{})
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}

func TestCampaignStoreService_DeleteByStoreID(t *testing.T) {
	campaignID := 3
	storeID := 101
	userID := 1234

	t.Run("when store with given campaign id and store id deleted successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		ShouldBeNil(err)

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		ShouldBeNil(err)

		campaignStoreService := NewCampaignStoreService(gdb)

		const sqlUpdate = `UPDATE "campaign_stores" SET "updated_by" = ? WHERE campaign_id = ? and store_id =?`
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).WithArgs(userID, campaignID, storeID).WillReturnResult(sqlmock.NewResult(0, 1))

		const sqlDelete = `DELETE FROM "campaign_stores" WHERE campaign_id = ? and store_id =?`
		mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).WithArgs(campaignID, storeID).WillReturnResult(sqlmock.NewResult(0, 1))

		err = campaignStoreService.DeleteByStoreID(context.TODO(), valueobjects.CampaignID(1), int64(123), int64(userID))
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}
