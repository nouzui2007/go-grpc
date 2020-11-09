package router

import (
	"goexample/configs"
	"goexample/db"
	"goexample/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// New Echoを生成する。
//
// ログやDBなどの初期設定を行う。
func New() *echo.Echo {

	e := echo.New()

	// set Middleware
	e.Use(middleware.Recover())

	// load config
	configs.NewConfig()

	// setting logger
	l := configs.Conf.Log
	l.LoadLevel()

	// set LogLevel
	e.Logger.SetLevel(l.Lvl)

	// set LogHeader
	if lg, ok := e.Logger.(*log.Logger); ok {
		lg.SetHeader(l.LogHeader)
	}

	d := db.New()

	// // set SQL log
	d.LogMode(l.Lvl <= configs.DEBUG)

	// create Handler
	h := handler.NewHandler(d)

	h.Register(e)

	return e
}
