package mysql

import (
	"time"

	"gorm.io/gorm"
)

type CampaignCodeService struct {
	db *gorm.DB
}

type CampaignCodeEntry struct {
	StatusCode  int64     `gorm:"primary_key;autoIncrement;column:status_code"`
	StatusValue string    `gorm:"column:status_value;type:varchar(100)"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
	CreatedBy   int64     `gorm:"column:created_by"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
	UpdatedBy   int64     `gorm:"column:updated_by"`
	DeletedAt   time.Time `gorm:"column:deleted_at;type:datetime"`
	DeletedBy   int64     `gorm:"column:deleted_by"`
}

func NewCampaignCodeService(db *gorm.DB) *CampaignCodeService {
	return &CampaignCodeService{db: db}
}

func (c *CampaignCodeEntry) TableName() string {
	return "campaign_codes"
}

func (c *CampaignCodeService) Migrate(defaultCampaignStatusDBEntry []map[string]interface{}) error {
	err := c.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&CampaignCodeEntry{})
	c.CreateCampainStatuses(defaultCampaignStatusDBEntry)
	return err
}

func (c *CampaignCodeService) CreateCampainStatuses(defaultCampaignStatusDBEntry []map[string]interface{}) {
	count := c.db.Model(&CampaignCodeEntry{}).Find(&CampaignCodeEntry{}).RowsAffected
	if count == 0 {
		c.db.Model(&CampaignCodeEntry{}).Create(defaultCampaignStatusDBEntry)
	}
}
