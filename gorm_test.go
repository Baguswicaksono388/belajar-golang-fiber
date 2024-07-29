package main

import (
	"strconv"
	"testing"

	"belajar-golang-fiber/entity"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(127.0.0.1:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //mengaktifkan logger pada GORM (semua perintah sql akan kelihatan)
	})
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

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func (tx *gorm.DB) error {
		err := tx.Create(&entity.User{ID: "10",Password: "rahasia",Name: entity.Name{FirstName: "User 10",}}).Error
		if err!= nil {
            return err
        }

		err = tx.Create(&entity.User{ID: "11",Password: "rahasia",Name: entity.Name{FirstName: "User 11",}}).Error
		if err != nil {
            return err
        }

		err = tx.Create(&entity.User{ID: "12",Password: "rahasia",Name: entity.Name{FirstName: "User 12",}}).Error
		if err != nil {
            return err
        }

		return nil
	})

	assert.Nil(t, err)
}


func TestTransactionError(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&entity.User{ID: "13",Password: "rahasia",Name: entity.Name{FirstName: "User 13",}}).Error
		if err!= nil {
            return err
        }

		err = tx.Create(&entity.User{ID: "11",Password: "rahasia",Name: entity.Name{FirstName: "User 11",}}).Error
		if err!= nil {
            return err
        }

		return nil
	})

	assert.NotNil(t, err) // ouputnya harus error
}

// DB Transaction manual tidak direkomendasikan
func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

    err := tx.Create(&entity.User{ID: "13", Password: "rahasia", Name: entity.Name{FirstName: "User 13",}}).Error
   	assert.Nil(t, err)

    err = tx.Create(&entity.User{ID: "14", Password: "rahasia", Name: entity.Name{FirstName: "User 14",}}).Error
    assert.Nil(t, err)
	
    if err == nil {
        tx.Commit()
    }
}

func TestManualTransactionError(t *testing.T) {
	tx := db.Begin()
    defer tx.Rollback()

    err := tx.Create(&entity.User{ID: "15", Password: "rahasia", Name: entity.Name{FirstName: "User 15",}}).Error
    assert.Nil(t, err)

    err = tx.Create(&entity.User{ID: "14", Password: "rahasia", Name: entity.Name{FirstName: "User 16",}}).Error
    assert.Nil(t, err)
	
	if err == nil {
    	tx.Commit()
    }
}

func TestQuerrySingleObject(t *testing.T) {
	user := entity.User{}
	err := db.First(&user).Error // mengambil 1 data pertama berdasarkan id
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)


	user = entity.User{}
	err = db.Last(&user).Error //mengambil 1 data terakhir berdasarkan id
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)
}

func TestQuerryInlineCondition(t *testing.T) {
	user := entity.User{}
	err := db.Take(&user, "id = ?", "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "User5", user.Name.FirstName)
}

func TestQueryAllObjects(t *testing.T) {
	var users []entity.User
	err := db.Find(&users, "id in ?", []string{"1","2","3","4"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []entity.User
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestQueryOrOperator(t *testing.T) {
	var users []entity.User
    err := db.Where("first_name like?", "%User%").Or("password =?", "rahasia").Find(&users).Error
    assert.Nil(t, err)
    assert.Equal(t, 14, len(users))
}

func TestQueryNotOperator(t *testing.T) {
	var users []entity.User
    err := db.Not("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
    assert.Nil(t, err)
    assert.Equal(t, 1, len(users))
}

func TestSelectFields(t *testing.T) {
	var users []entity.User
	err := db.Select("id", "first_name").Find(&users).Error
	assert.Nil(t, err)
	
	for _, user := range users {
		assert.NotNil(t, "id", user.ID)
        assert.NotEqual(t, "", user.Name.FirstName) 
	}

	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondtion := entity.User {
		Name: entity.Name{
			FirstName: "User5",
			LastName: "", // tidak bisa karena default value string
		},
		Password: "rahasia",
	}

	var users []entity.User
	err := db.Where(userCondtion).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}


func TestMapCondition(t *testing.T) {
	mapCondtion := map[string]interface{}{
		"middle_name" : "",
		"last_name" : "",
	}

	var users []entity.User
	err := db.Where(mapCondtion).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []entity.User
    err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
    assert.Nil(t, err)
    assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID string
	FirstName string
	LastName string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	err :=db.Model(&entity.User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
}
