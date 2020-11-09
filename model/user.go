package model

import (
	"goexample/pkg"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"
	"gopkg.in/go-playground/validator.v9"
)

type (
	//User はdtb_userの構造体です。
	User struct {
		BaseModel
		WelbyID        *int32     `json:"welby_user_id" gorm:"type:int(10);NOT NULL" validate:"required"`
		KartePatientID *int32     `json:"karte_patient_id" gorm:"type:int(10)"`
		Expired        *time.Time `json:"expired"`
	}
)

//NewUser はUserに初期値を設定します。
func NewUser() *User {
	return &User{
		BaseModel: NewBaseModel(),
	}
}

//BeforeSave はgorm.DB#Create、gorm.DB#Update、gorm.DB#Saveのコールバック。不要なら削除。ここはサンプル用
func (m *User) BeforeSave() error {
	validate := validator.New()
	if err := validate.Struct(m); err != nil {
		er := pkg.NewValidateError(err)
		return &er
	}
	return nil
}

//FindByID はidによってユーザ情報を取得します。
func (m *User) FindByID(db *gorm.DB, id *int32) error {
	rnf := db.Where("id = ?", id).First(m).RecordNotFound()
	if err := db.Error; err != nil {
		return xerrors.Errorf("query error: %w", err)
	}
	// data not found
	if rnf {
		fmt.Println("data not found")
		return &pkg.AppError{Code: http.StatusNotFound, Err: xerrors.New("data not found")}
	}
	return nil
}

//FindByWelbyID はwelbyIDによってユーザ情報を取得します。
func (m *User) FindByWelbyID(db *gorm.DB, welbyID *int32) (*User, error) {
	var r = NewUser()
	rnf := db.Where("welby_id = ?", welbyID).First(r).RecordNotFound()
	if err := db.Error; err != nil {
		return nil, xerrors.Errorf("query error: %w", err)
	}
	// data not found
	if rnf {
		return nil, &pkg.AppError{Code: http.StatusNotFound, Err: xerrors.New("data not found")}
	}
	return r, nil
}

//Create は新規ユーザ情報を作成します。
func (m *User) Create(db *gorm.DB) error {
	if err := db.Create(m).Error; err != nil {
		if pkg.IsAppError(err) {
			return err
		}
		er := pkg.NewQueryError(err)
		return &er
	}

	return nil
}

//Update は既存ユーザ情報を更新します。
func (m *User) Update(db *gorm.DB) error {
	if err := db.Model(m).Update(m).Error; err != nil {
		if pkg.IsAppError(err) {
			return err
		}
		er := pkg.NewQueryError(err)
		return &er
	}

	return nil
}

//ClearExpired は expired カラムに null を更新します。
func (m *User) ClearExpired(db *gorm.DB) error {
	if err := db.Model(m).Update("expired", gorm.Expr("NULL")).Error; err != nil {
		if pkg.IsAppError(err) {
			return err
		}
		er := pkg.NewQueryError(err)
		return &er
	}

	return nil
}

//IsExpired は期限切れかどうかを返す。
//
//Expiredが現在時刻より後ならtrue、そうでなければfalseを返す。Expiredがnullの場合はtrue
func (m *User) IsExpired() bool {
	if m.Expired == nil {
		return true
	}
	return m.Expired.Before(time.Now())
}
