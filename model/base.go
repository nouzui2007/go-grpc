package model

import (
	"goexample/pkg"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type (
	//BaseModel は全てのModelが保持するカラムが設定されています。
	BaseModel struct {
		ID        *int32     `json:"id" gorm:"primary_key"`
		CreatedAt *time.Time `json:"-"`
		UpdatedAt *time.Time `json:"-"`
		Deleted   *int32     `json:"-"`
	}

	//NamedBaseModel は名称付きのベースモデル
	NamedBaseModel struct {
		BaseModel
		Name *string `json:"name" gorm:"type:varchar(45);NOT NULL`
	}
)

//NewBaseModel は初期値をつめてオブジェクトを生成します。
func NewBaseModel() BaseModel {
	return BaseModel{
		// Created:  pkg.TimeToTPtr(time.Now()),
		// Modified: pkg.TimeToTPtr(time.Now()),
		Deleted: pkg.IntToIPtr(0),
	}
}

//NewNamedBaseModel は初期値をつめてオブジェクトを生成します。
func NewNamedBaseModel() NamedBaseModel {
	return NamedBaseModel{
		BaseModel: NewBaseModel(),
	}
}

//BeforeSave はBaseModelがSaveされる前にバリデーションを実施します。
func (b *BaseModel) BeforeSave() error {
	validate := validator.New()
	return validate.Struct(b)
}

//BeforeCreate はBaseModelがCreateされる前にバリデーションを実施します。
func (b *BaseModel) BeforeCreate() error {
	validate := validator.New()
	return validate.Struct(b)
}

//BeforeUpdate はBaseModelがUpdateされる前にバリデーションを実施します。
func (b *BaseModel) BeforeUpdate() error {
	validate := validator.New()
	return validate.Struct(b)
}
