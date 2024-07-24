package main

import (
	"strconv"
	"testing"

	"belajar-golang-fiber/entity"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(127.0.0.1:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("insert into sample(id, name) values (?, ?)", "1", "Bagus").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id, name) values (?, ?)", "2", "Budi").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id, name) values (?, ?)", "3", "Joko").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id, name) values (?, ?)", "4", "Rully").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id string
	Name string
}

func TestRawSql(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Bagus", sample.Name)


	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

func TestSqlRow(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close() // Agar tidak memory lick

	var samples []Sample
	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{
			Id: id,
			Name: name,
		})
	}
	assert.Equal(t, 4, len(samples))
}

func TestScanRow(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close() // Agar tidak memory lick

	var samples []Sample
	for rows.Next() {
		err := db.ScanRows(rows, &samples) // Scan rows dari Gorm
		assert.Nil(t, err)
	}
	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	// import Struct User from folder Entity
	user := entity.User{ID: "1", Password: "rahasia", Name: entity.Name{
		FirstName: "Bagus",
		LastName:  "Wicaksono",
		MiddleName: "Testing",
	}, Information: "ini akan di ignore",}

    response := db.Create(&user)
    assert.Nil(t, response.Error)
    assert.Equal(t, int64(1), response.RowsAffected)
}

// Memasukkan data lebih dari 1
func TestBatchInsert(t *testing.T) {
	var users []entity.User
	for i :=2; i < 10; i++ {
		users = append(users, entity.User{
			ID: strconv.Itoa(i),
			Password: "rahasia",
			Name: entity.Name{
				FirstName: "User" + strconv.Itoa(i),
			},
		})
	}

	result := db.Create(&users)
	assert.Nil(t,  result.Error)
	assert.Equal(t, 8, int(result.RowsAffected))
}