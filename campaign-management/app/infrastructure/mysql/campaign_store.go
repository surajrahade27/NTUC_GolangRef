package mysql

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
	"database/sql"
	"fmt"

	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignStoreService struct {
	db *gorm.DB
}

type CampaignStoreEntry struct {
	ID         int64          `gorm:"primary_key;autoIncrement;column:campaign_store_id"`
	CampaignID int64          `gorm:"column:campaign_id;"`
	StoreID    int64          `gorm:"column:store_id"`
	CreatedAt  sql.NullTime   `gorm:"column:created_at;type:datetime"`
	CreatedBy  int64          `gorm:"column:created_by"`
	UpdatedAt  sql.NullTime   `gorm:"column:updated_at;type:datetime"`
	UpdatedBy  int64          `gorm:"column:updated_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;type:datetime"`
	DeletedBy  int64          `gorm:"column:deleted_by"`
}

func NewCampaignStoreService(db *gorm.DB) *CampaignStoreService {
	return &CampaignStoreService{db: db}
}

func (c *CampaignStoreEntry) TableName() string {
	return "campaign_stores"
}

func (c *CampaignStoreService) Migrate() error {
	err := c.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&CampaignStoreEntry{})
	return err
}

func (c *CampaignStoreService) CreateMultiple(ctx context.Context, stores []entities.CampaignStore) ([]entities.CampaignStore, error) {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}
	entries := []CampaignStoreEntry{}
	for i := range stores {
		entries = append(entries, c.ToEntry(stores[i]))
	}
	err := db.Create(&entries).Error
	if err != nil {
		return []entities.CampaignStore{}, fmt.Errorf("%w: %v", valueobjects.ErrStoreCantCreate, err)
	}
	logger.Infof("stores added for campaign id : %v", entries[0].CampaignID)

	var storeEntities []entities.CampaignStore
	for _, store := range entries {
		storeEntities = append(storeEntities, c.ToEntity(store))
	}
	return storeEntities, nil
}

func (c *CampaignStoreService) ToEntry(storeEntity entities.CampaignStore) CampaignStoreEntry {
	return CampaignStoreEntry{
		ID:         storeEntity.ID.ToInt64(),
		CampaignID: storeEntity.CampaignID.ToInt64(),
		StoreID:    storeEntity.StoreID,
		CreatedBy:  storeEntity.CreatedBy,
		UpdatedBy:  storeEntity.UpdatedBy,
	}
}

func (c *CampaignStoreService) ToEntity(storeEntry CampaignStoreEntry) entities.CampaignStore {
	return entities.CampaignStore{
		ID:         valueobjects.CampaignStoreID(storeEntry.ID),
		CampaignID: valueobjects.CampaignID(storeEntry.CampaignID),
		StoreID:    storeEntry.StoreID,
	}
}

func (c *CampaignStoreService) GetByCampaignId(ctx context.Context, CampaignID valueobjects.CampaignID) ([]entities.CampaignStore, error) {
	var entry []CampaignStoreEntry
	err := c.db.Where("campaign_id = ?", CampaignID).Find(&entry).Error
	return c.ToEntityList(entry), err
}

func (c *CampaignStoreService) ToEntityList(entries []CampaignStoreEntry) []entities.CampaignStore {
	var stores []entities.CampaignStore
	for _, entry := range entries {
		store := c.ToEntity(entry)
		stores = append(stores, store)
	}
	return stores
}

func (c *CampaignStoreService) Update(ctx context.Context, campaignStore entities.CampaignStore) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}

	campaignStoreEntry := c.ToEntry(campaignStore)
	err := db.Model(&campaignStoreEntry).Where("campaign_store_id = ? and campaign_id = ?",
		campaignStoreEntry.ID, campaignStoreEntry.CampaignID).Updates(&campaignStoreEntry).Error
	if err != nil {
		logger.Errorf("error occured : %v", err.Error())
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantUpdate, err)
	}
	logger.Infof("store with id %v updated successfully", campaignStoreEntry.ID)
	return nil
}

func (c *CampaignStoreService) DeleteByCampaignID(ctx context.Context, campaignID valueobjects.CampaignID, userID int64) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}
	var entry CampaignStoreEntry
	err := db.Model(&CampaignStoreEntry{}).Where("campaign_id = ?", campaignID.ToInt64()).Update("deleted_by", userID).Error
	if err != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantUpdate, err)
	}

	response := db.Where("campaign_id = ?", campaignID.ToInt64()).Delete(&entry)
	if response.Error != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantDelete, response.Error)
	}
	if response.RowsAffected < 1 {
		logger.Infof("stores with campaign id %v not exists", campaignID.ToInt64())
		return nil
	}
	logger.Infof("stores with  campaign id %v deleted successfully", campaignID.ToInt64())
	return nil
}

func (c *CampaignStoreService) Delete(ctx context.Context, campaignID valueobjects.CampaignID, campaignStoreID valueobjects.CampaignStoreID, userID int64) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}
	var entry CampaignStoreEntry
	err := db.Model(&CampaignStoreEntry{}).Where("campaign_store_id = ? and campaign_id = ?",
		campaignStoreID, campaignID).Update("deleted_by", userID).Error
	if err != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantUpdate, err)
	}

	response := db.Where("campaign_store_id = ? and campaign_id = ?", campaignStoreID.ToInt64(),
		campaignID.ToInt64()).Delete(&entry)
	if response.Error != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantDelete, response.Error)
	}
	if response.RowsAffected < 1 {
		logger.Infof("campaign store with id %v not exists", campaignStoreID.ToInt64())
		return nil
	}
	logger.Infof("campaign store with id %v deleted successfully", campaignStoreID.ToInt64())
	return nil
}

func (c *CampaignStoreService) GetByStoreID(ctx context.Context, campaignID valueobjects.CampaignID, storeID int64) (entities.CampaignStore, error) {
	var entry CampaignStoreEntry
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}
	err := db.Where("campaign_id = ? and store_id = ?", campaignID, storeID).First(&entry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("entity not found")
			return entities.CampaignStore{}, fmt.Errorf("%w", valueobjects.ErrStoreNotExists)
		}
		return entities.CampaignStore{}, fmt.Errorf("%w: %v", valueobjects.ErrStoreCantGet, err)
	}
	return c.ToEntity(entry), nil
}

func (c *CampaignStoreService) DeleteByStoreID(ctx context.Context, campaignID valueobjects.CampaignID, storeID, userID int64) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}
	var entry CampaignStoreEntry
	err := db.Model(&CampaignStoreEntry{}).Where("store_id = ? and campaign_id = ?",
		storeID, campaignID).Update("deleted_by", userID).Error
	if err != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantUpdate, err)
	}

	response := db.Where("store_id = ? and campaign_id = ?", storeID, campaignID.ToInt64()).Delete(&entry)
	if response.Error != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrStoreCantDelete, response.Error)
	}
	if response.RowsAffected < 1 {
		logger.Infof("campaign store with store id %v not exists", storeID)
		return nil
	}
	logger.Infof("campaign store with store id %v deleted successfully", storeID)
	return nil
}
