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

func TestCampaignProductUseCase_AddProducts(t *testing.T) {
	campaignID := 101
	userID := 987654321
	productEntities := []entities.CampaignProduct{
		{
			ProductID:   123,
			SKUNo:       1111,
			SerialNo:    1,
			SequenceNo:  1,
			ProductType: "cd",
			CampaignID:  valueobjects.CampaignID(campaignID),
			CreatedBy:   int64(userID),
		},
		{
			ProductID:   456,
			SKUNo:       2222,
			SerialNo:    2,
			SequenceNo:  2,
			ProductType: "cd",
			CampaignID:  valueobjects.CampaignID(campaignID),
			CreatedBy:   int64(userID),
		},
	}
	t.Run("when products for particular campaign added successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)

		productsDTO := []dto.CampaignProducts{
			{
				ID:          1,
				ProductID:   456,
				SKUNo:       2222,
				SerialNo:    2,
				SequenceNo:  2,
				ProductType: "cd",
			},
			{
				ID:          2,
				ProductID:   456,
				SKUNo:       2222,
				SerialNo:    2,
				SequenceNo:  2,
				ProductType: "cd",
			},
		}
		mockCampaignProductService.On("CreateMultiple", ctx, productEntities).Return([]entities.CampaignProduct{
			{
				ID:          1,
				ProductID:   456,
				SKUNo:       2222,
				SerialNo:    2,
				SequenceNo:  2,
				ProductType: "cd",
				CampaignID:  valueobjects.CampaignID(campaignID),
				CreatedBy:   int64(userID),
			},
			{
				ID:          2,
				ProductID:   456,
				SKUNo:       2222,
				SerialNo:    2,
				SequenceNo:  2,
				ProductType: "cd",
				CampaignID:  valueobjects.CampaignID(campaignID),
				CreatedBy:   int64(userID),
			},
		}, nil)

		response, err := productUseCase.AddProducts(ctx, productEntities)
		ShouldEqual(response, productsDTO)
		ShouldBeNil(err)
	})
	t.Run("when error occured while saving product details in db", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		mockCampaignProductService.On("CreateMultiple", ctx, productEntities).Return(
			[]entities.CampaignProduct{}, fmt.Errorf("%w: %v", valueobjects.ErrProductCantCreate, errors.New("db error")))

		_, err := productUseCase.AddProducts(ctx, productEntities)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrProductCantCreate) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignProductUseCase_UpdateProducts(t *testing.T) {
	campaignID := 101
	userID := 987654321
	productEntities := []entities.CampaignProduct{
		{
			ID:          valueobjects.CampaignProductID(111),
			ProductID:   123,
			SKUNo:       1111,
			SerialNo:    1,
			SequenceNo:  1,
			ProductType: "ncd",
			CampaignID:  valueobjects.CampaignID(campaignID),
			CreatedBy:   int64(userID),
		},
		{
			ID:          valueobjects.CampaignProductID(222),
			ProductID:   456,
			SKUNo:       2222,
			SerialNo:    2,
			SequenceNo:  2,
			ProductType: "ncd",
			CampaignID:  valueobjects.CampaignID(campaignID),
			CreatedBy:   int64(userID),
		},
	}
	t.Run("when products for particular campaign updated successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		for _, product := range productEntities {
			mockCampaignProductService.On("Update", ctx, product).Return(nil)
		}

		err := productUseCase.UpdateProducts(ctx, productEntities)
		ShouldBeNil(err)
	})
	t.Run("when error occured while updating product details in db", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		mockCampaignProductService.On("Update", ctx, productEntities[0]).Return(
			fmt.Errorf("%w: %v", valueobjects.ErrProductCantUpdate, errors.New("db error")))

		err := productUseCase.UpdateProducts(ctx, productEntities)
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
		if !errors.As(err, &valueobjects.ErrProductCantUpdate) {
			t.Error("invalid error type")
		}
	})
}

func TestCampaignProductUseCase_DeleteProductByCampaignId(t *testing.T) {
	campaignID := 101
	userID := 987654321
	productEntity := entities.CampaignProduct{
		ID:          valueobjects.CampaignProductID(1000),
		ProductID:   100,
		SKUNo:       1111,
		SerialNo:    1,
		SequenceNo:  1,
		ProductType: "ncd",
		CampaignID:  valueobjects.CampaignID(campaignID),
		CreatedBy:   int64(userID),
	}

	t.Run("when product for particular campaign deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		mockCampaignProductService.On("DeleteByCampaignId", ctx, int64(productEntity.CampaignID), productEntity.ProductID).Return(nil)
		err := productUseCase.DeleteByCampaignId(ctx, int64(productEntity.CampaignID), productEntity.ProductID)
		ShouldBeNil(err)
	})
	t.Run("when product for particular campaign deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		mockCampaignProductService.On("DeleteByCampaignId", ctx, int64(productEntity.CampaignID), productEntity.ProductID).Return(errors.New("record Not Found"))
		err := productUseCase.DeleteByCampaignId(ctx, int64(productEntity.CampaignID), productEntity.ProductID)
		ShouldNotBeNil(err)
	})
}

func TestCampaignProductUseCase_GetProducts(t *testing.T) {
	t.Run("When campaign product exist, it returns product data", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductsService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductsService)
		campaignID := 101
		userID := 987654321
		response := []entities.CampaignProduct{
			{
				ID:          valueobjects.CampaignProductID(111),
				ProductID:   123,
				SKUNo:       1111,
				SerialNo:    1,
				SequenceNo:  1,
				ProductType: "cd",
				CampaignID:  valueobjects.CampaignID(campaignID),
				CreatedBy:   int64(userID),
			},
			{
				ID:          valueobjects.CampaignProductID(222),
				ProductID:   456,
				SKUNo:       2222,
				SerialNo:    2,
				SequenceNo:  2,
				ProductType: "cd",
				CampaignID:  valueobjects.CampaignID(campaignID),
				CreatedBy:   int64(userID),
			},
		}
		mockCampaignProductsService.On("GetByCampaignId", ctx, valueobjects.CampaignID(campaignID)).Return(
			response,
			nil,
		)
		actualValue, err := productUseCase.GetProducts(ctx, int64(campaignID))
		ShouldEqual(actualValue, response)
		ShouldBeNil(err)
	})
	t.Run("When campaign details not exist, it returns error", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductsService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductsService)
		campaignID := 101
		response := []entities.CampaignProduct{}
		mockCampaignProductsService.On("GetByCampaignId", ctx, valueobjects.CampaignID(campaignID)).Return(
			response,
			errors.New("record Not Found"),
		)
		actualValue, err := productUseCase.GetProducts(ctx, int64(campaignID))
		ShouldEqual(actualValue, response)
		ShouldNotBeNil(err)
	})
}
func TestCampaignProductUseCase_DeleteAllProductByCampaignId(t *testing.T) {
	campaignID := 101
	userID := 987654321
	productEntity := entities.CampaignProduct{
		ID:          valueobjects.CampaignProductID(1000),
		ProductID:   100,
		SKUNo:       1111,
		SerialNo:    1,
		SequenceNo:  1,
		ProductType: "ncd",
		CampaignID:  valueobjects.CampaignID(campaignID),
		CreatedBy:   int64(userID),
	}

	t.Run("when all products for particular campaign deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		mockCampaignProductService.On("DeleteAllByCampaignId", ctx, int64(productEntity.CampaignID)).Return(nil)
		err := productUseCase.DeleteAllByCampaignId(ctx, int64(productEntity.CampaignID))
		ShouldBeNil(err)
	})
	t.Run("when all products for particular campaign deleted successfully", func(t *testing.T) {
		ctx := context.Background()
		mockCampaignProductService := mocks.NewCampaignProducts(t)
		productUseCase := NewCampaignProductUseCase(mockCampaignProductService)
		mockCampaignProductService.On("DeleteAllByCampaignId", ctx, int64(productEntity.CampaignID)).Return(errors.New("record Not Found"))
		err := productUseCase.DeleteAllByCampaignId(ctx, int64(productEntity.CampaignID))
		ShouldNotBeNil(err)
	})
}
