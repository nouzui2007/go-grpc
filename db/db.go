/*
Package db はDataBase接続用のパッケージです。

コネクション上限数、コネクションライフサイクル等を設定しています。
*/
package db

import (
	"fmt"
	"time"

	"goexample/configs"

	"github.com/jinzhu/gorm"
)

//New はDB接続オブジェクトを作成します。
//
// コネクション上限数、コネクションライフサイクルをconfigの値によって設定します。
func New() *gorm.DB {

	dbConf := configs.Conf.Database

	db, err := gorm.Open(dbConf.Dbms, dbConf.DSN())
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(dbConf.MaxIdleConnections)
	db.DB().SetMaxOpenConns(dbConf.MaxOpenConnections)
	db.DB().SetConnMaxLifetime(time.Duration(dbConf.MaxLifeTime) * time.Second)
	db.LogMode(dbConf.LogMode)

	db.Callback().Query().Before("gorm:query").Register("new_before_query_callback", newBeforeQueryFunction)
	db.Callback().RowQuery().Before("gorm:row_query").Register("new_before_row_query_callback", newBeforeQueryFunction)

	return db
}

//newBeforeQueryFunction は論理削除されたレコードを検索対象から除外する
func newBeforeQueryFunction(scope *gorm.Scope) {
	var (
		quotedTableName               = scope.QuotedTableName()
		deletedField, hasDeletedField = scope.FieldByName("Deleted")
		defaultUnixTime               = 0
	)

	if !scope.Search.Unscoped && hasDeletedField {
		scope.Search.Unscoped = true
		sql := fmt.Sprintf("%v.%v = '%v'", quotedTableName, scope.Quote(deletedField.DBName), defaultUnixTime)
		scope.Search.Where(sql)
	}
}
