package dto

import (
	"campaign-mgmt/app/domain/entities"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ToCampaignStoreDTO(t *testing.T) {
	t.Run("test conversion", func(t *testing.T) {
		campaignStoreEntity := entities.CampaignStore{
			CampaignID: 1,
			StoreID:    123,
		}

		expectedResponse := CampaignStores{
			ID:      int64(1),
			StoreID: int64(123),
		}
		response := ToCampaignStoreDTO(campaignStoreEntity)
		ShouldEqual(response, expectedResponse)
	})
}
