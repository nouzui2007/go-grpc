/*
Package handler は各モデルに対する処理の呼び出しを行うパッケージです。

*/
package handler

import (
	"github.com/jinzhu/gorm"
)

//Handler はDB接続を持つ構造体です。
//Handlerが保持しているDB接続情報を用いて各処理を行います。
type Handler struct {
	db *gorm.DB
}

//NewHandler はHandlerを新しく生成します。
func NewHandler(d *gorm.DB) *Handler {
	return &Handler{
		db: d,
	}
}
