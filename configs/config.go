/*
Package configs は設定値用のパッケージです。

将来的にconst.goは削除したいと思っています。
*/
package configs

import (
	"fmt"
	"sync/atomic"

	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
)

//Conf はconfig.tomlの値を持つ構造体です。
var Conf Config

type (
	//Config はconfig.tomlに記載されている設定をまとめたstruct
	//
	//実際の値はそれぞれのstructが保持します。
	Config struct {
		Database      DatabaseConfig
		Message       MessageConfig
		Log           LogConfig
		TestDatabase  TestDatabaseConfig
		Request       Request
		Response      Response
		Authorization Authorization
	}
	//DatabaseConfig はDB接続用の情報です。
	DatabaseConfig struct {
		Dbms               string
		Server             string
		Port               string
		Database           string
		User               string
		Password           string
		LogMode            bool
		ParseTime          string
		MaxIdleConnections int
		MaxOpenConnections int
		MaxLifeTime        int
		TestEnvFlg         bool
	}
	//MessageConfig はエラーメッセージ用の情報です。 (利用していないので削除)
	MessageConfig struct {
		E0001 string
	}
	//LogConfig はログ出力用の情報です。
	LogConfig struct {
		Level     string
		LogHeader string
		Lvl       log.Lvl
		Mask      []string
	}
	//TestDatabaseConfig はunitテスト用DB用の情報です。
	TestDatabaseConfig struct {
		Dbms               string
		Server             string
		Port               string
		Database           string
		User               string
		Password           string
		LogMode            bool
		ParseTime          string
		MaxIdleConnections int
		MaxOpenConnections int
		MaxLifeTime        int
	}
	//Request はユーザ認可用の情報です。
	Request struct {
		AuthHeaderKey string
		DeviceKey     string
		VersionKey    string
	}
	Response struct {
		UpdateRequiredKey string
	}
	//MyKarte はマイカルテアクセス用の情報です。
	MyKarte struct {
		AuthHeaderKey        string
		Domain               string
		Version              string
		Login                string
		Logout               string
		Patient              string
		User                 string
		UserConfirm          string
		UserConfirmationCode string
	}
	Vital struct {
		APIKey         string
		BasicAuthMode  bool
		BasicAuthValue string
		Domain         string
		Version        string
		Vitals         string
		Search         string
	}
	//Medicine は薬剤サーバアクセス用の情報です。
	Medicine struct {
		BasicAuthMode  bool
		BasicAuthValue string
		Domain         string
		NumberPerPage  int
		SearchEntry    string
		SearchType     string
		DrugSearch     string
	}
	//Authorization はユーザ認証
	Authorization struct {
		Secret  string
		Expired int
	}
)

//NewConfig は configオブジェクトを作成します。
func NewConfig() {
	var config Config

	if _, cErr := toml.DecodeFile("./configs/config.toml", &config); cErr != nil {
		fmt.Println(cErr)
	}

	Conf = config
}

//NewTestConfig は unitテスト用 configオブジェクトを作成します。 (利用していないので削除)
func NewTestConfig() {
	var config Config

	if _, cErr := toml.DecodeFile("../configs/config.toml", &config); cErr != nil {
		fmt.Println(cErr)
	}

	Conf = config
}

//NewSpecifiedTestConfig はunitテスト用 config オブジェクトを作成します。
func NewSpecifiedTestConfig(s string) {
	var config Config

	if _, cErr := toml.DecodeFile(s, &config); cErr != nil {
		fmt.Println(cErr)
	}

	Conf = config
}

//DSN はconfigからDSNを整形
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s", d.User, d.Password, d.Server, d.Port, d.Database, d.ParseTime)
}

//DSN はunitテスト用 configからDSNを整形
func (td TestDatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s", td.User, td.Password, td.Server, td.Port, td.Database, td.ParseTime)
}

//LoadLevel はログレベル設定
func (l *LogConfig) LoadLevel() {
	var u uint32

	switch l.Level {
	case "DEBUG":
		u = uint32(DEBUG)
	case "INFO":
		u = uint32(INFO)
	case "WARN":
		u = uint32(WARN)
	case "ERROR":
		u = uint32(ERROR)
	case "FATAL":
		u = uint32(FATAL)
	}

	l.Lvl = log.Lvl(atomic.LoadUint32(&u))
}
