package model

import (
	"goexample/pkg"
	"goexample/test"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestUser_Create(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	t.Run("正常に生成", func(t *testing.T) {
		user := NewUser()
		user.WelbyID = pkg.IntToIPtr(1)
		user.KartePatientID = pkg.IntToIPtr(10)
		user.Create(db)
		assert.NotNil(t, *user.ID)
	})

	t.Run("WelbyIDなし", func(t *testing.T) {
		user := NewUser()
		user.KartePatientID = pkg.IntToIPtr(10)
		err := user.Create(db)
		assert.Error(t, err)
	})

	t.Run("KartePatientIDなし", func(t *testing.T) {
		user := NewUser()
		user.WelbyID = pkg.IntToIPtr(2)
		user.Create(db)
		assert.NotNil(t, *user.ID)
	})
}

func TestUser_Update(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	createInitialUserData(db)

	t.Run("正常に更新", func(t *testing.T) {
		user := NewUser()
		user.ID = pkg.IntToIPtr(1)
		db.First(user)

		user.WelbyID = pkg.IntToIPtr(1000)
		user.Update(db)
		if err := db.Where("welby_id = ?", 1000).First(user).Error; err != nil {
			assert.Fail(t, "更新した値での検索不可")
		}
	})

	t.Run("WelbyIDなし", func(t *testing.T) {
		user := NewUser()
		user.ID = pkg.IntToIPtr(1)
		db.First(user)

		user.WelbyID = nil
		err := user.Update(db)
		assert.Error(t, err)
	})

	t.Run("KartePatientIDなし", func(t *testing.T) {
		user := NewUser()
		user.ID = pkg.IntToIPtr(1)
		db.First(user)

		user.WelbyID = pkg.IntToIPtr(2000)
		user.KartePatientID = nil
		user.Update(db)
		if err := db.Where("welby_id = ?", 2000).First(user).Error; err != nil {
			assert.Fail(t, "更新した値での検索不可")
		}
	})

}

func TestUser_FindByID(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	createInitialUserData(db)

	t.Run("正常に取得できる", func(t *testing.T) {
		user := NewUser()
		user.FindByID(db, pkg.IntToIPtr(1))
		assert.Equal(t, user.ID, pkg.IntToIPtr(1))
	})

	t.Run("削除されたID", func(t *testing.T) {
		err := user.FindByID(db, pkg.IntToIPtr(2))
		assert.Error(t, err)
	})

	t.Run("存在しないID", func(t *testing.T) {
		user := NewUser()
		err := user.FindByID(db, pkg.IntToIPtr(3))
		assert.Error(t, err)
	})
}

func TestUser_FindByWelbyID(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	createInitialUserData(db)

	user := NewUser()

	t.Run("正常に取得できる", func(t *testing.T) {
		result, _ := user.FindByWelbyID(db, pkg.IntToIPtr(100))
		assert.Equal(t, result.WelbyID, pkg.IntToIPtr(100))
		assert.Equal(t, result.KartePatientID, pkg.IntToIPtr(101))
	})

	t.Run("削除されたWelbyID", func(t *testing.T) {
		_, err := user.FindByWelbyID(db, pkg.IntToIPtr(200))
		assert.Error(t, err)
	})

	t.Run("存在しないWelbyID", func(t *testing.T) {
		_, err := user.FindByWelbyID(db, pkg.IntToIPtr(300))
		assert.Error(t, err)
	})
}

func TestUser_ClearExpired(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	expired := time.Date(2020, time.June, 1, 15, 00, 00, 0, time.UTC)
	user := User{
		WelbyID:        pkg.IntToIPtr(100),
		KartePatientID: pkg.IntToIPtr(101),
		Expired:        pkg.TimeToTPtr(expired),
		BaseModel: BaseModel{
			ID:        pkg.IntToIPtr(1),
			CreatedAt: pkg.TimeToTPtr(time.Now()),
			UpdatedAt: pkg.TimeToTPtr(time.Now()),
			Deleted:   pkg.IntToIPtr(0),
		},
	}
	db.Create(&user)
	assert.NotNil(t, user.Expired)

	t.Run("正常にクリア", func(t *testing.T) {
		user.ClearExpired(db)
		confirm := NewUser()
		confirm.ID = pkg.IntToIPtr(1)
		db.First(confirm)
		assert.Nil(t, confirm.Expired)
	})
}

func createInitialUserData(db *gorm.DB) {
	user := User{
		WelbyID:        pkg.IntToIPtr(100),
		KartePatientID: pkg.IntToIPtr(101),
		Expired:        pkg.TimeToTPtr(time.Now()),
		BaseModel: BaseModel{
			ID:        pkg.IntToIPtr(1),
			CreatedAt: pkg.TimeToTPtr(time.Now()),
			UpdatedAt: pkg.TimeToTPtr(time.Now()),
			Deleted:   pkg.IntToIPtr(0),
		},
	}
	db.Create(&user)

	deletedUser := User{
		WelbyID:        pkg.IntToIPtr(200),
		KartePatientID: pkg.IntToIPtr(201),
		Expired:        pkg.TimeToTPtr(time.Now()),
		BaseModel: BaseModel{
			ID:        pkg.IntToIPtr(2),
			CreatedAt: pkg.TimeToTPtr(time.Now()),
			UpdatedAt: pkg.TimeToTPtr(time.Now()),
			Deleted:   pkg.IntToIPtr(1),
		},
	}
	db.Create(&deletedUser)
}
