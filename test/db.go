package test

import (
	"goexample/configs"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"

	"github.com/jinzhu/gorm"
)

// dbConn はテスト時に使用するDBコネクションです。
var dbConn *gorm.DB

type cleanKind struct{ value string }

var cleanKinds = struct {
	drop     cleanKind
	truncate cleanKind
}{
	drop:     cleanKind{"DROP TABLE"},
	truncate: cleanKind{"TRUNCATE"},
}

// SetupDBConn ... testパッケージにDBへの接続を持っておく
//
//返り値の funcを実行するとdb接続close
func SetupDBConn() func() {
	// testDataBaseに対する接続情報取得
	_, pwd, _, _ := runtime.Caller(0)
	configPath := fmt.Sprintf("%s/../configs/config.toml", path.Dir(pwd))
	configs.NewSpecifiedTestConfig(configPath)
	dbConf := configs.Conf.TestDatabase

	db, err := gorm.Open(dbConf.Dbms, dbConf.DSN())
	if err != nil {
		fmt.Println("storage err: ", err)
	}

	//set connection pool
	db.DB().SetMaxIdleConns(dbConf.MaxIdleConnections)
	db.DB().SetMaxOpenConns(dbConf.MaxOpenConnections)
	db.LogMode(dbConf.LogMode)

	db.Callback().Query().Before("gorm:query").Register("new_before_query_callback", newBeforeQueryFunction)
	db.Callback().RowQuery().Before("gorm:row_query").Register("new_before_row_query_callback", newBeforeQueryFunction)

	dbConn = db

	createTablesIfNotExist()

	return func() {
		// cleanAllTables(cleanKinds.drop)
		dbConn.Close()
	}
}

// GetTestDBConn ... プールしてあるテスト用のDBコネクションを返す
//
//返り値のfunc()実行で table truncate
func GetTestDBConn() (*gorm.DB, func()) {
	if dbConn == nil {
		fmt.Println("db not connection")
	}
	return dbConn, func() {
		cleanAllTables(cleanKinds.truncate)
	}
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

func execSchema(fPath string) {
	b, err := ioutil.ReadFile(fPath)
	if err != nil {
		fmt.Printf("schema reading error: %v", err)
	}

	queries := strings.Split(string(b), ";")

	for _, query := range queries[:len(queries)-1] {
		dbConn.Exec(query)
		if err != nil {
			fmt.Printf("exec schema error: %v, query: %s", err, query)
		}
	}
}

func createTablesIfNotExist() {
	_, pwd, _, _ := runtime.Caller(0)
	schemaPath := fmt.Sprintf("%s/../schema/1_create_tables.sql", path.Dir(pwd))
	fmt.Printf("create tables path => %s", schemaPath)
	execSchema(schemaPath)
}

func cleanAllTables(kind cleanKind) {

	rows, err := dbConn.Raw("SHOW TABLES").Rows()

	if err != nil {
		fmt.Printf("show tables error: %#v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			fmt.Printf("show table error: %#v", err)
			continue
		}

		cmds := []string{
			"SET FOREIGN_KEY_CHECKS = 0",
			fmt.Sprintf("%s `"+configs.Conf.TestDatabase.Database+"`.%s", kind.value, tableName),
			"SET FOREIGN_KEY_CHECKS = 1",
		}
		for _, cmd := range cmds {
			if err := dbConn.Exec(cmd).Error; err != nil {
				fmt.Printf("drop error: %#v\n", err)
				continue
			}
		}
	}

}
