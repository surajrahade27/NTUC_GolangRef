package mysql

import (
	"time"

	"gorm.io/gorm"
)

type StoreDailyTimeSlotService struct {
	db *gorm.DB
}

type StoreDailyTimeSlotEntry struct {
	ID              int64     `gorm:"primary_key;autoIncrement;column:daily_time_slot_id"`
	StoreID         int64     `gorm:"column:store_id;type:bigint;not null"`
	StartTime       time.Time `gorm:"column:start_time;type:time"`
	EndTime         time.Time `gorm:"column:end_time;type:time"`
	Quota           int       `gorm:"column:quota;type:smallint"`
	DayofWeek       string    `gorm:"column:day_of_week;type:varchar(20)"`
	IsSlotAvailable bool      `gorm:"column:is_slot_available;type:boolean"`
	UserID          int64     `gorm:"column:user_id;type:bigint"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime"`
	CreatedBy       int64     `gorm:"column:created_by;type:bigint"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime"`
	UpdatedBy       int64     `gorm:"column:updated_by;type:bigint"`
	DeletedAt       time.Time `gorm:"column:deleted_at;type:datetime"`
	DeletedBy       int64     `gorm:"column:deleted_by;type:bigint"`
}

func NewStoreDailyTimeSlotService(db *gorm.DB) *StoreDailyTimeSlotService {
	return &StoreDailyTimeSlotService{db: db}
}

func (c *StoreDailyTimeSlotEntry) TableName() string {
	return "store_daily_time_slots"
}

func (c *StoreDailyTimeSlotService) Migrate() error {
	err := c.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&StoreDailyTimeSlotEntry{})
	return err
}
