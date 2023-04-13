package mysql

import (
	"time"

	"gorm.io/gorm"
)

type StoreSpecificTimeSlotService struct {
	db *gorm.DB
}

type StoreSpecificTimeSlotEntry struct {
	ID        int64     `gorm:"primary_key;autoIncrement;column:Specific_time_slot_id"`
	StoreID   int64     `gorm:"column:store_id;type:bigint;not null"`
	Date      time.Time `gorm:"column:date;type:date"`
	StartTime time.Time `gorm:"column:start_time;type:time"`
	EndTime   time.Time `gorm:"column:end_time;type:time"`
	Quota     int       `gorm:"column:quota;type:smallint"`
	UserID    int64     `gorm:"column:user_id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime"`
	CreatedBy int64     `gorm:"column:created_by;type:bigint"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime"`
	UpdatedBy int64     `gorm:"column:updated_by;type:bigint"`
	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime"`
	DeletedBy int64     `gorm:"column:deleted_by;type:bigint"`
}

func NewStoreSpecificTimeSlotService(db *gorm.DB) *StoreSpecificTimeSlotService {
	return &StoreSpecificTimeSlotService{db: db}
}

func (c *StoreSpecificTimeSlotEntry) TableName() string {
	return "store_specific_time_slots"
}

func (c *StoreSpecificTimeSlotService) Migrate() error {
	err := c.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&StoreSpecificTimeSlotEntry{})
	return err
}
