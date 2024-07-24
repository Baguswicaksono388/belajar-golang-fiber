package entity

import "time"

// User => Users (defaul convension dari nama table)
type User struct {
	ID        string `gorm:"primary_key;column:id;<-:create"` // Jika penulisan ID Gorm sudah mempresentasikan primary key, jadi tidak perlu menggunakan tag. <-:create hanya untuk dibuat saja pada Bab Field Permission
	Password  string `gorm:"column:password"`
	Name      Name `gorm:"embedded"` // Grouping Field Name
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"` // penulisan ini sudah sesuai dengan conversation dari GORM, jd ini sudah autoCreatedTime ketika data dibuat tambah menambahkan tag
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Information string `gorm:"-"` // field permission: tidak ada read/write permission
}

func (u *User) TableName() string {
	return "users" // nama table yang diinginkan
    // return "users_table" // jika ingin menggunakan nama table yang berbeda
}

type Name struct {
	FirstName string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName string `gorm:"column:last_name"`
}
