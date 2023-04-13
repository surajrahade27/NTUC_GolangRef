package mysql

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
	"database/sql"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignService struct {
	db *gorm.DB
}

type CampaignEntry struct {
	ID                  int64          `gorm:"primary_key;autoIncrement;column:campaign_id"`
	Title               string         `gorm:"column:title;type:varchar(1024)"`
	OrderStartDate      sql.NullTime   `gorm:"column:order_start_date;type:datetime"`
	OrderEndDate        sql.NullTime   `gorm:"column:order_end_date;type:datetime"`
	CollectionStartDate sql.NullTime   `gorm:"column:collection_start_date;type:datetime"`
	CollectionEndDate   sql.NullTime   `gorm:"column:collection_end_date;type:datetime"`
	StatusCode          int64          `gorm:"column:status_code"`
	CampaignType        string         `gorm:"column:campaign_type;type:varchar(1024)"`
	ListingTitle        string         `gorm:"column:listing_title;type:varchar(1024)"`
	ListingDesc         string         `gorm:"column:listing_description;type:text"`
	ListingImagePath    string         `gorm:"column:listing_image_path;type:text"`
	OnboardTitle        string         `gorm:"column:onboard_title;type:varchar(1024)"`
	OnboardDesc         string         `gorm:"column:onboard_description;type:text"`
	OnboardImagePath    string         `gorm:"column:onboard_image_path;type:text"`
	LandingImagePath    string         `gorm:"column:landing_image_path;type:text"`
	LeadTime            sql.NullInt32  `gorm:"column:lead_time;type:smallint;default:NULL"`
	OfferID             sql.NullInt64  `gorm:"column:offer_id;default:NULL"`
	TagID               sql.NullInt64  `gorm:"column:tag_id;default:NULL"`
	IsCampaignPublished *bool          `gorm:"column:is_campaign_published;type:boolean;default:false"`
	CreatedAt           sql.NullTime   `gorm:"column:created_at;type:datetime"`
	CreatedBy           int64          `gorm:"column:created_by"`
	UpdatedAt           sql.NullTime   `gorm:"column:updated_at;type:datetime"`
	UpdatedBy           int64          `gorm:"column:updated_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;type:datetime"`
	DeletedBy           int64          `gorm:"column:deleted_by"`
}

func NewCampaignService(db *gorm.DB) *CampaignService {
	return &CampaignService{db: db}
}

func (c *CampaignEntry) TableName() string {
	return "campaigns"
}

func (c *CampaignService) Migrate() error {
	err := c.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&CampaignEntry{})
	return err
}

func (c *CampaignService) Get(ctx context.Context, id valueobjects.CampaignID) (entities.Campaign, error) {
	entry := CampaignEntry{}
	err := c.db.First(&entry, id).Error
	return c.ToEntity(entry), err
}

func (c *CampaignService) GetList(ctx context.Context, pagination entities.PaginationConfig) ([]entities.Campaign, int64, error) {
	var entries []CampaignEntry
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := c.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	var result *gorm.DB

	if pagination.Status > int64(0) {
		result = queryBuilder.Model(&CampaignEntry{}).Where("title like ? and status_code = ?", pagination.Name+"%", pagination.Status).Find(&entries)
	} else {
		result = queryBuilder.Model(&CampaignEntry{}).Where("title like ? ", pagination.Name+"%").Find(&entries)
	}

	return c.ToEntityList(entries), c.GetCampaignsCount(), result.Error
}

func (c *CampaignService) GetCampaignsCount() int64 {
	var count int64
	c.db.Table("campaigns").Count(&count)
	return count
}

func (c *CampaignService) ToEntity(entry CampaignEntry) entities.Campaign {
	var isCampaignPublished bool
	var collectionStartDate, collectionEndDate, orderStartDate, orderEndDate time.Time
	var leadTime int
	var offerID, tagID int64
	if entry.IsCampaignPublished != nil {
		isCampaignPublished = *entry.IsCampaignPublished
	} else {
		isCampaignPublished = false
	}
	if entry.CollectionStartDate.Valid {
		collectionStartDate = entry.CollectionStartDate.Time.UTC()
	}
	if entry.CollectionEndDate.Valid {
		collectionEndDate = entry.CollectionEndDate.Time.UTC()
	}
	if entry.OrderStartDate.Valid {
		orderStartDate = entry.OrderStartDate.Time.UTC()
	}
	if entry.OrderEndDate.Valid {
		orderEndDate = entry.OrderEndDate.Time.UTC()
	}
	if entry.LeadTime.Valid {
		leadTime = int(entry.LeadTime.Int32)
	}
	if entry.OfferID.Valid {
		offerID = entry.OfferID.Int64
	}
	if entry.TagID.Valid {
		tagID = entry.TagID.Int64
	}
	return entities.Campaign{
		ID:                  valueobjects.CampaignID(entry.ID),
		Title:               entry.Title,
		CollectionStartDate: collectionStartDate,
		CollectionEndDate:   collectionEndDate,
		OrderStartDate:      orderStartDate,
		OrderEndDate:        orderEndDate,
		StatusCode:          entry.StatusCode,
		ListingTitle:        entry.ListingTitle,
		ListingDesc:         entry.ListingDesc,
		ListingImagePath:    entry.ListingImagePath,
		OnboardTitle:        entry.OnboardTitle,
		OnboardDesc:         entry.OnboardDesc,
		OnboardImagePath:    entry.OnboardImagePath,
		LandingImagePath:    entry.LandingImagePath,
		CampaignType:        valueobjects.CampaignType(entry.CampaignType),
		LeadTime:            leadTime,
		OfferID:             offerID,
		TagID:               tagID,
		IsCampaignPublished: isCampaignPublished,
	}
}

func (c *CampaignService) ToEntry(campaignEntity entities.Campaign) CampaignEntry {
	var collectionStartDate, collectionEndDate, orderStartDate, orderEndDate sql.NullTime
	var leadTime sql.NullInt32
	var offerID, tagID sql.NullInt64
	if !campaignEntity.CollectionStartDate.IsZero() {
		collectionStartDate = sql.NullTime{
			Time:  campaignEntity.CollectionStartDate,
			Valid: true,
		}
	}
	if !campaignEntity.CollectionEndDate.IsZero() {
		collectionEndDate = sql.NullTime{
			Time:  campaignEntity.CollectionEndDate,
			Valid: true,
		}
	}
	if !campaignEntity.OrderStartDate.IsZero() {
		orderStartDate = sql.NullTime{
			Time:  campaignEntity.OrderStartDate,
			Valid: true,
		}
	}
	if !campaignEntity.OrderEndDate.IsZero() {
		orderEndDate = sql.NullTime{
			Time:  campaignEntity.OrderEndDate,
			Valid: true,
		}
	}
	if campaignEntity.LeadTime != 0 {
		leadTime = sql.NullInt32{
			Int32: int32(campaignEntity.LeadTime),
			Valid: true,
		}
	}
	if campaignEntity.OfferID != 0 {
		offerID = sql.NullInt64{
			Int64: campaignEntity.OfferID,
			Valid: true,
		}
	}
	if campaignEntity.TagID != 0 {
		tagID = sql.NullInt64{
			Int64: campaignEntity.TagID,
			Valid: true,
		}
	}
	return CampaignEntry{
		ID:                  campaignEntity.ID.ToInt64(),
		Title:               campaignEntity.Title,
		CollectionStartDate: collectionStartDate,
		CollectionEndDate:   collectionEndDate,
		OrderStartDate:      orderStartDate,
		OrderEndDate:        orderEndDate,
		StatusCode:          campaignEntity.StatusCode,
		ListingTitle:        campaignEntity.ListingTitle,
		ListingDesc:         campaignEntity.ListingDesc,
		ListingImagePath:    campaignEntity.ListingImagePath,
		OnboardTitle:        campaignEntity.OnboardTitle,
		OnboardDesc:         campaignEntity.OnboardDesc,
		OnboardImagePath:    campaignEntity.OnboardImagePath,
		LandingImagePath:    campaignEntity.LandingImagePath,
		CampaignType:        campaignEntity.CampaignType.String(),
		LeadTime:            leadTime,
		OfferID:             offerID,
		TagID:               tagID,
		CreatedBy:           campaignEntity.CreatedBy,
		UpdatedBy:           campaignEntity.UpdatedBy,
		IsCampaignPublished: &campaignEntity.IsCampaignPublished,
	}
}

func (c *CampaignService) ToEntityList(entries []CampaignEntry) []entities.Campaign {
	var campaigns []entities.Campaign
	for _, entry := range entries {
		campaign := c.ToEntity(entry)
		campaigns = append(campaigns, campaign)
	}
	return campaigns
}

func (c *CampaignService) Create(ctx context.Context, campaign entities.Campaign) (entities.Campaign, error) {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}

	entry := c.ToEntry(campaign)
	err := db.Create(&entry).Error
	if err != nil {
		return entities.Campaign{}, fmt.Errorf("%w: %v", valueobjects.ErrCampaignCantCreate, err)
	}
	logger.Infof("campaign created with id : %v", entry.ID)
	return c.ToEntity(entry), nil
}

func (c *CampaignService) Exists(ctx context.Context, id valueobjects.CampaignID, title string) (bool, error) {
	campaign := CampaignEntry{}
	if id != 0 {
		err := c.db.Where("campaign_id = ?", id).First(&campaign).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				logger.Info("entity not found")
				return false, nil
			}
			return false, fmt.Errorf("%w:  %v", valueobjects.ErrCampaignCantExist, err)
		}
		return true, nil
	} else {
		err := c.db.Where("title = ?", title).First(&campaign).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				logger.Info("entity not found")
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
}

func (c *CampaignService) Update(ctx context.Context, campaign entities.Campaign) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}

	var exists bool
	err := db.Model(CampaignEntry{}).Select("count(*) > 0").Where("campaign_id != ? and title = ?", campaign.ID, campaign.Title).Find(&exists).Error
	if err != nil {
		return fmt.Errorf("error occured while checking campaign existence with given name: %v", err)
	}

	if exists {
		return fmt.Errorf("campaign with title '%s' already exists. please provide another title", campaign.Title)
	}

	campaignEntry := c.ToEntry(campaign)
	err = db.Model(campaignEntry).Where("campaign_id = ?", campaignEntry.ID).Updates(map[string]interface{}{
		"title":                 campaignEntry.Title,
		"order_start_date":      campaignEntry.OrderStartDate,
		"order_end_date":        campaignEntry.OrderEndDate,
		"collection_start_date": campaignEntry.CollectionStartDate,
		"collection_end_date":   campaignEntry.CollectionEndDate,
		"status_code":           campaignEntry.StatusCode,
		"campaign_type":         campaignEntry.CampaignType,
		"listing_title":         campaignEntry.ListingTitle,
		"listing_description":   campaignEntry.ListingDesc,
		"listing_image_path":    campaignEntry.ListingImagePath,
		"onboard_title":         campaignEntry.OnboardTitle,
		"onboard_description":   campaignEntry.OnboardDesc,
		"onboard_image_path":    campaignEntry.OnboardImagePath,
		"landing_image_path":    campaignEntry.LandingImagePath,
		"lead_time":             campaignEntry.LeadTime,
		"offer_id":              campaignEntry.OfferID,
		"tag_id":                campaignEntry.TagID,
		"is_campaign_published": campaignEntry.IsCampaignPublished}).Error
	if err != nil {
		return fmt.Errorf("%w: %v", valueobjects.ErrCampaignCantUpdate, err)
	}
	logger.Infof("campaign with id : %v updated successfully", campaignEntry.ID)
	return nil
}

func (c *CampaignService) UpdateStatus(ctx context.Context) error {
	db := c.db
	if ctxDB := DBTransaction(ctx); ctxDB != nil {
		db = ctxDB
	}
	err := publishCampaigns(db)
	if err != nil {
		return err
	}
	err = deactivateCampaigns(db)
	if err != nil {
		return err
	}
	return nil
}

func publishCampaigns(db *gorm.DB) error {
	logger.Info("publishing campaigns")
	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
	campaign := CampaignEntry{}
	response := db.Model(&campaign).Where("CAST(order_start_date AS DATE)  <= ? and status_code = 3", today).Update("status_code", 2)
	if response.Error != nil {
		return response.Error
	}
	if response.RowsAffected > 0 {
		logger.Infof("published %d campaigns", response.RowsAffected)
		return nil
	} else {
		logger.Info("no campaign to publish")
		return nil
	}
}

func deactivateCampaigns(db *gorm.DB) error {
	logger.Info("deactivating campaigns")
	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
	campaign := CampaignEntry{}
	response := db.Model(&campaign).Where("CAST(order_end_date AS DATE) < ? and status_code = 2", today).Update("status_code", 1)
	if response.Error != nil {
		return response.Error
	}
	if response.RowsAffected > 0 {
		logger.Infof("deactivated %d campaigns", response.RowsAffected)
		return nil
	} else {
		logger.Info("no campaign to deactivate")
		return nil
	}
}
