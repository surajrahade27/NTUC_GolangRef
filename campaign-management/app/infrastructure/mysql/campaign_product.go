package mysql

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignProductService struct {
	db *gorm.DB
}

type CampaignProductEntry struct {
	ID          int64          `gorm:"primary_key;autoIncrement;column:campaign_product_id"`
	CampaignID  int64          `gorm:"column:campaign_id"`
	ProductID   int64          `gorm:"column:product_id"`
	SKUNo       int64          `gorm:"column:SKU_no"`
	SerialNo    int            `gorm:"column:serial_no;type:smallint"`
	SequenceNo  int            `gorm:"column:sequence_no;type:smallint"`
	ProductType string         `gorm:"column:product_type;type:varchar(20)"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:datetime"`
	CreatedBy   int64          `gorm:"column:created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:datetime"`
	UpdatedBy   int64          `gorm:"column:updated_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:datetime"`
	DeletedBy   int64          `gorm:"column:deleted_by"`
}

func NewCampaignProductService(db *gorm.DB) *CampaignProductService {
	return &CampaignProductService{db: db}
}

func (c *CampaignProductEntry) TableName() string {
	return "campaign_products"
}

func (c *CampaignProductService) Migrate() error {
	err := c.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&CampaignProductEntry{})
	return err
}

func (c *CampaignProductService) CreateMultiple(ctx context.Context, products []entities.CampaignProduct) ([]entities.CampaignProduct, error) {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}

	productEntries := []CampaignProductEntry{}
	for i := range products {
		productEntries = append(productEntries, c.ToEntry(products[i]))
	}
	err := db.Create(&productEntries).Error
	if err != nil {
		return []entities.CampaignProduct{}, fmt.Errorf("%w: %v", valueobjects.ErrProductCantCreate, err)
	}
	logger.Infof("products added for campaign id : %v", productEntries[0].CampaignID)

	var productEntities []entities.CampaignProduct
	for _, product := range productEntries {
		productEntities = append(productEntities, c.ToEntity(product))
	}
	return productEntities, nil
}

func (c *CampaignProductService) ToEntry(productEntity entities.CampaignProduct) CampaignProductEntry {
	return CampaignProductEntry{
		ID:          productEntity.ID.ToInt64(),
		CampaignID:  productEntity.CampaignID.ToInt64(),
		ProductID:   productEntity.ProductID,
		SKUNo:       productEntity.SKUNo,
		SerialNo:    productEntity.SerialNo,
		SequenceNo:  productEntity.SequenceNo,
		ProductType: productEntity.ProductType,
		CreatedBy:   productEntity.CreatedBy,
		UpdatedBy:   productEntity.UpdatedBy,
	}
}

func (c *CampaignProductService) ToEntity(productEntry CampaignProductEntry) entities.CampaignProduct {
	return entities.CampaignProduct{
		ID:          valueobjects.CampaignProductID(productEntry.ID),
		CampaignID:  valueobjects.CampaignID(productEntry.CampaignID),
		ProductID:   productEntry.ProductID,
		SKUNo:       productEntry.SKUNo,
		SerialNo:    productEntry.SerialNo,
		SequenceNo:  productEntry.SequenceNo,
		ProductType: productEntry.ProductType,
	}
}

func (c *CampaignProductService) GetByCampaignId(ctx context.Context, CampaignID valueobjects.CampaignID) ([]entities.CampaignProduct, error) {
	var entry []CampaignProductEntry
	err := c.db.Where("campaign_id = ?", CampaignID).Find(&entry).Error
	return c.ToEntityList(entry), err
}

func (c *CampaignProductService) ToEntityList(entries []CampaignProductEntry) []entities.CampaignProduct {
	var products []entities.CampaignProduct
	for _, entry := range entries {
		product := c.ToEntity(entry)
		products = append(products, product)
	}
	return products
}

func (c *CampaignProductService) Update(ctx context.Context, campaignProduct entities.CampaignProduct) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}

	campaignProductEntry := c.ToEntry(campaignProduct)
	err := db.Model(&campaignProductEntry).Where("campaign_product_id = ? and campaign_id = ?",
		campaignProductEntry.ID, campaignProductEntry.CampaignID).Updates(&campaignProductEntry).Error
	if err != nil {
		logger.Errorf("error occured : %v", err.Error())
		return fmt.Errorf("%w: %v", valueobjects.ErrProductCantUpdate, err)
	}
	logger.Infof("product with id %v updated successfully", campaignProductEntry.ID)
	return nil
}

func (c *CampaignProductService) DeleteByCampaignId(ctx context.Context, campaignID int64, productID int64) error {
	result := c.db.Where("campaign_id = ? and product_id =?", campaignID, productID).Delete(&CampaignProductEntry{})
	if result.Error != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrProductCantDelete, result.Error)
	}
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Info("entity not found")
			return nil
		}
		return result.Error
	}
	return nil
}

func (c *CampaignProductService) DeleteAllByCampaignId(ctx context.Context, campaignID int64) error {
	result := c.db.Where("campaign_id = ?", campaignID).Delete(&CampaignProductEntry{})
	if result.Error != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrProductCantDelete, result.Error)
	}
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Info("entity not found")
			return nil
		}
		return result.Error
	}
	return nil
}
