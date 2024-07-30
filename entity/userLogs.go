package entity

import "time"

// User => Users (defaul convension dari nama table)
type UserLogs struct {
	ID        string `gorm:"primary_key;column:id;autoIncrement"`
	UserId  string `gorm:"column:user_id"`
	Action string `gorm:"column:action"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *UserLogs) TableName() string {
	return "user_logs"
}
