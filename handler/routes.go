package handler

import (
	"goexample/configs"
	"goexample/model"
	"goexample/pkg"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

//CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

//NewValidator は新しいCustomValidatorを作成します。
func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

//Validate はCustomValidatorに定義されたvalidateタグの変数をvalidateします
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

//Context はecho.Contextをラップしたものです。
//
//contextの他にmodelやuserInfoを持ちます。modelは全てのmodelが実装しているBaseModelerを定義しています。
//
type Context struct {
	echo.Context
	userInfo model.User
}

//Register はHandlerではルーティングの設定を行なっています。
//
//例: GET:/user/home のリクエストがきた場合、h.GetHomeInfoを実行する等
//
//
func (h *Handler) Register(router *echo.Echo) {

	router.Use(middleware.BodyDump(bodyDumpHandler))

	v1 := router.Group("v1")

	// アクセストークンを用いて認可しないものは先に宣言
	login := v1.Group("/login")
	login.POST("", h.Login)

	//users
	users := v1.Group("/users")
	users.POST("", h.CreateUser)
	users.PUT("", h.UpdateUser)

	v1.Use(AuthMiddleware(h.db))

	router.Validator = NewValidator()

	//logout
	logout := v1.Group("/logout")
	logout.POST("", h.Logout)
}

//AuthMiddleware はアクセストークンの認可を行います。
//
//認可できなかった場合403を返します。
//
//アクセストークンの有効性の認可はマイカルテで行います。マイカルテの/patientへアクセスして成功の場合は認証済みです。
func AuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			header := c.Request().Header
			var accessToken string

			c.Logger().Debug("Target Header:", configs.Conf.Request.AuthHeaderKey)

			// header.Getは大文字小文字区別せずにチェックします
			accessToken = header.Get(configs.Conf.Request.AuthHeaderKey)

			if accessToken == "" {
				c.Echo().Logger.Error("Access Token Not found")
				return echo.NewHTTPError(http.StatusUnauthorized, "Access Token Not found")
			}

			c.Logger().Debug("accessToken:", accessToken)

			// jwtの認証
			claim := jwtClaim{}
			token, err := jwt.ParseWithClaims(accessToken, &claim, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv(configs.Conf.Authorization.Secret)), nil
			})
			if err != nil {
				c.Logger().Error("Access Token does not parse")
				er := pkg.NewValidateError(err)
				return c.JSON(er.Code, er)
			}
			if !token.Valid {
				c.Logger().Error("Access Token Not verified")
				er := pkg.NewInvalidTokenError()
				return c.JSON(er.Code, er)
			}
			c.Logger().Debug("Access Token Verified", claim)

			// DSAF get user
			user := model.NewUser()
			if err := user.FindByID(db, &claim.UserID); err != nil {
				if appErr, ok := err.(*pkg.AppError); ok && appErr.Code == http.StatusNotFound {
					c.Logger().Error("User not found")
					er := pkg.NewInvalidTokenError()
					return c.JSON(er.Code, er)
				}
				return err
			}
			// 有効期限チェック
			if user.IsExpired() {
				c.Logger().Error("Access Token expred")
				er := pkg.NewTokenExpiredError()
				return c.JSON(er.Code, er)
			}

			err = next(&Context{
				Context:  c,
				userInfo: *user,
			})

			return err
		}
	}
}

//bodyDumpHandler はRequestの内容を標準出力にログ出力します。
//
//出力内容は Request:httpMethod,URL RequestBody:RequestBody(ユーザが自由入力した箇所はマスク化してます) Response Status:ResponseStatus
func bodyDumpHandler(c echo.Context, reqBody, _ []byte) {
	c.Logger().Info("Request:", c.Request().Method, c.Request().URL)
	fmt.Printf("Request Body:%v\n", maskBody(reqBody))
	c.Logger().Info(fmt.Sprintf("Response Status:%v", c.Response().Status))
}

//マスク対象となっているkeyのvalueをマスクします。
func maskBody(b []byte) string {

	var s []interface{}
	var err error

	if err = json.Unmarshal(b, &s); err == nil {

		for k, v := range s {

			var bytes []byte
			if bytes, err = json.Marshal(v); err != nil {
				return ""
			}
			// call maskBody
			s[k] = maskBody(bytes)
		}
	} else {
		//  map[string]interface{}{}の場合 mask処理
		m := map[string]interface{}{}
		var err error
		if err = json.Unmarshal(b, &m); err != nil {
			return ""
		}

		for k, v := range m {
			// 配列になっている場合 maskBody
			if reflect.TypeOf(v) == reflect.TypeOf([]interface{}{}) {
				var bytes []byte
				if bytes, err = json.Marshal(v); err != nil {
					return ""
				}
				m[k] = maskBody(bytes)
			}
		}

		// configに設定されたkeyの場合mask
		for _, k := range configs.Conf.Log.Mask {
			if _, ok := m[k]; ok {
				m[k] = "***"
			}
		}
		return fmt.Sprint(m)
	}

	return fmt.Sprint(s)
}
