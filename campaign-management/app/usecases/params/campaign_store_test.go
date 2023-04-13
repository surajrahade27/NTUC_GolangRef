package params

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ToCampaignStoreEntity(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		expectedResponse := entities.CampaignStore{
			CampaignID: valueobjects.CampaignID(1),
			StoreID:    123,
			CreatedBy:  999,
		}

		response := ToCampaignStoreEntity(int64(123), int64(1), int64(999))
		ShouldEqual(response, expectedResponse)
	})
}

func Test_ToUpdateCampaignStoreEntity(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		store := UpdateCampaignStore{
			ID:      101,
			StoreID: 123,
		}

		expectedResponse := entities.CampaignStore{
			ID: valueobjects.CampaignStoreID(101),
			CampaignID: valueobjects.CampaignID(1),
			StoreID:    123,
			UpdatedBy:  999,
		}


		response := ToUpdateCampaignStoreEntity(store, int64(1), int64(999))
		ShouldEqual(response, expectedResponse)
	})
}
