package usecases

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/services/mocks"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/dto"
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCampaignStoreUseCase_AddStores(t *testing.T) {
	campaignID := 101
	userID := 987654321
	storeEntities := []entities.CampaignStore{
		{
			StoreID:    1234,
			CampaignID: valueobjects.CampaignID(campaignID),
			CreatedBy:  int64(userID),
		},
		{
			StoreID:    5678,
			CampaignID: valueobjects.CampaignID(campaignID),
			CreatedBy:  int64(userID),
		},
	}
	t.Run("when stores for particular campaign added successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		storesDTO := []dto.CampaignStores{
			{
				ID:      1,
				StoreID: 1234,
			},
			{
				ID:      2,
				StoreID: 5678,
			},
		}
		mockCampaignStoreService.On("CreateMultiple", ctx, storeEntities).Return([]entities.CampaignStore{
			{
				ID:         1,
				StoreID:    1234,
				CampaignID: valueobjects.CampaignID(campaignID),
				CreatedBy:  int64(userID),
			},
			{
				ID:         2,
				StoreID:    5678,
				CampaignID: valueobjects.CampaignID(campaignID),
				CreatedBy:  int64(userID),
			},
		}, nil)

		response, err := storeUseCase.AddStores(ctx, storeEntities)
		ShouldEqual(response, storesDTO)
		ShouldBeNil(err)
	})
	t.Run("when error occured while saving store details in db", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		mockCampaignStoreService.On("CreateMultiple", ctx, storeEntities).Return(
			[]entities.CampaignStore{}, fmt.Errorf("%w: %v", valueobjects.ErrStoreCantCreate, errors.New("db error")))

		_, err := storeUseCase.AddStores(ctx, storeEntities)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrStoreCantCreate) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignStoreUseCase_UpdateStores(t *testing.T) {
	campaignID := 101
	userID := 987654321
	storeEntities := []entities.CampaignStore{
		{
			ID:         valueobjects.CampaignStoreID(1),
			StoreID:    1234,
			CampaignID: valueobjects.CampaignID(campaignID),
			CreatedBy:  int64(userID),
		},
		{
			ID:         valueobjects.CampaignStoreID(2),
			StoreID:    5678,
			CampaignID: valueobjects.CampaignID(campaignID),
			CreatedBy:  int64(userID),
		},
	}
	t.Run("when stores for particular campaign updated successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		for _, store := range storeEntities {
			mockCampaignStoreService.On("Update", ctx, store).Return(nil)
		}

		err := storeUseCase.UpdateStores(ctx, storeEntities)
		ShouldBeNil(err)
	})
	t.Run("when error occured while updating store details in db", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		mockCampaignStoreService.On("Update", ctx, storeEntities[0]).Return(
			fmt.Errorf("%w: %v", valueobjects.ErrStoreCantUpdate, errors.New("db error")))

		err := storeUseCase.UpdateStores(ctx, storeEntities)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrStoreCantUpdate) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignStoreUseCase_GetStores(t *testing.T) {
	t.Run("When campaign Store exist, it returns Store data", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		campaignID := 101
		userID := 987654321
		response := []entities.CampaignStore{
			{
				ID:         valueobjects.CampaignStoreID(1),
				StoreID:    1234,
				CampaignID: valueobjects.CampaignID(campaignID),
				CreatedBy:  int64(userID),
			},
			{
				ID:         valueobjects.CampaignStoreID(2),
				StoreID:    5678,
				CampaignID: valueobjects.CampaignID(campaignID),
				CreatedBy:  int64(userID),
			},
		}
		mockCampaignStoreService.On("GetByCampaignId", ctx, valueobjects.CampaignID(campaignID)).Return(
			response,
			nil,
		)
		actualValue, err := storeUseCase.GetStores(ctx, int64(campaignID))
		ShouldEqual(actualValue, response)
		ShouldBeNil(err)
	})
	t.Run("When campaign store details not exist, it returns error", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		campaignID := 101
		response := []entities.CampaignStore{}
		mockCampaignStoreService.On("GetByCampaignId", ctx, valueobjects.CampaignID(campaignID)).Return(
			response,
			errors.New("record Not Found"),
		)
		actualValue, err := storeUseCase.GetStores(ctx, int64(campaignID))
		ShouldEqual(actualValue, response)
		ShouldNotBeNil(err)
	})
}

func TestCampaignStoreUseCase_DeleteStores(t *testing.T) {
	campaignID := int64(101)
	userID := int64(232323)
	t.Run("when all the stores for particular campaign deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockStoreService)

		mockStoreService.On("DeleteByCampaignID", ctx, valueobjects.CampaignID(campaignID), userID).Return(nil)

		err := storeUseCase.DeleteStores(ctx, campaignID, userID)
		ShouldBeNil(err)
	})

	t.Run("error occured while deleting stores entries from db", func(t *testing.T) {
		ctx := context.Background()
		mockStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockStoreService)

		mockStoreService.On("DeleteByCampaignID", ctx, valueobjects.CampaignID(campaignID), userID).Return(
			fmt.Errorf("%w: %v", valueobjects.ErrStoreCantDelete, errors.New("db error")))
		err := storeUseCase.DeleteStores(ctx, campaignID, userID)

		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrStoreCantDelete) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignStoreUseCase_DeleteStore(t *testing.T) {
	campaignID := int64(101)
	storeID := int64(2121)
	userID := int64(232323)
	t.Run("when store with given id deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockStoreService)

		mockStoreService.On("Delete", ctx, valueobjects.CampaignID(campaignID), valueobjects.CampaignStoreID(storeID),
			userID).Return(nil)

		err := storeUseCase.DeleteStore(ctx, campaignID, storeID, userID)
		ShouldBeNil(err)
	})

	t.Run("error occured while deleting particular store entry from db", func(t *testing.T) {
		ctx := context.Background()
		mockStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockStoreService)

		mockStoreService.On("Delete", ctx, valueobjects.CampaignID(campaignID), valueobjects.CampaignStoreID(storeID),
			userID).Return(fmt.Errorf("%w: %v", valueobjects.ErrStoreCantDelete, errors.New("db error")))
		err := storeUseCase.DeleteStore(ctx, campaignID, storeID, userID)

		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrStoreCantDelete) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignStoreUseCase_GetByStoreID(t *testing.T) {
	campaignID := valueobjects.CampaignID(123)
	userID := 987654321
	t.Run("when store with given campaign and store id fetched successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		storeDTO := dto.CampaignStores{
			ID:      1,
			StoreID: 1234,
		}
		mockCampaignStoreService.On("GetByStoreID", ctx, campaignID, int64(1234)).Return(entities.CampaignStore{
			ID:         1,
			StoreID:    1234,
			CampaignID: valueobjects.CampaignID(campaignID),
			CreatedBy:  int64(userID),
		}, nil)

		response, err := storeUseCase.GetByStoreID(ctx, int64(campaignID), int64(1234))
		ShouldEqual(response, storeDTO)
		ShouldBeNil(err)
	})

	t.Run("error occured while getting store details", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockCampaignStoreService)
		mockCampaignStoreService.On("GetByStoreID", ctx, campaignID, int64(1234)).Return(entities.CampaignStore{},
			fmt.Errorf("%w: %v", valueobjects.ErrStoreCantGet, errors.New("db error")))
		_, err := storeUseCase.GetByStoreID(ctx, int64(campaignID), int64(1234))
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrStoreCantGet) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignStoreUseCase_DeleteByStoreID(t *testing.T) {
	campaignID := int64(101)
	storeID := int64(2121)
	userID := int64(232323)
	t.Run("when store with given store id deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockStoreService)

		mockStoreService.On("DeleteByStoreID", ctx, valueobjects.CampaignID(campaignID), storeID,
			userID).Return(nil)

		err := storeUseCase.DeleteByStoreID(ctx, campaignID, storeID, userID)
		ShouldBeNil(err)
	})

	t.Run("error occured while deleting store entry from db", func(t *testing.T) {
		ctx := context.Background()
		mockStoreService := mocks.NewCampaignStores(t)
		storeUseCase := NewCampaignStoreUseCase(mockStoreService)

		mockStoreService.On("DeleteByStoreID", ctx, valueobjects.CampaignID(campaignID), storeID,
			userID).Return(fmt.Errorf("%w: %v", valueobjects.ErrStoreCantDelete, errors.New("db error")))
		err := storeUseCase.DeleteByStoreID(ctx, campaignID, storeID, userID)

		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrStoreCantDelete) {
			t.Error("invalid error type")
		}
	})
}