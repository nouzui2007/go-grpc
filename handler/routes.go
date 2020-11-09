package handler

import (
	"goexample/configs"
	"goexample/model/entity"
	"goexample/pkg"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"

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
	userInfo entity.User
	karteAPI myKarte.Function
	vitalAPI *vital.API
}

//Register はHandlerではルーティングの設定を行なっています。
//
//例: GET:/user/home のリクエストがきた場合、h.GetHomeInfoを実行する等
//
//
func (h *Handler) Register(router *echo.Echo) {

	router.Use(middleware.BodyDump(bodyDumpHandler))

	v1 := router.Group("v1")
	v1.Use(CheckVersionMiddleware(h.db))

	appVersion := v1.Group("/app_version")
	appVersion.GET("", h.GetAppVersion)

	// アクセストークンを用いて認可しないものは先に宣言
	login := v1.Group("/login")
	login.POST("", h.Login)

	//users
	users := v1.Group("/users")
	users.POST("", h.CreateUser)
	users.PUT("", h.UpdateUser)
	users.POST("/confirmationcode", h.ConfirmationCode)

	v1.Use(AuthMiddleware(h.db))

	router.Validator = NewValidator()

	//logout
	logout := v1.Group("/logout")
	logout.POST("", h.Logout)

	//af_symptoms
	afSymptoms := v1.Group("/af_symptoms")
	afSymptoms.GET("", h.GetAfSymptoms)

	//articles
	articles := v1.Group("/articles")
	articles.GET("", h.GetArticles)

	//article_brows_histories
	articleBrowsHistories := v1.Group("/article_brows_histories")
	articleBrowsHistories.POST("", h.CreateArticleBrowsHistory)

	//dose_histories
	doseHistories := v1.Group("/dose_histories")
	doseHistories.GET("/:record_date_start/:record_date_to", h.GetDoseHistoriesBetweenRecordDates)
	doseHistories.POST("", h.CreateDoseHistory)
	doseHistories.PUT("/:id", h.UpdateDoseHistory)

	//dose_notifications
	doseNotifications := v1.Group("/dose_notifications")
	doseNotifications.GET("", h.GetDoseNotifications)
	doseNotifications.POST("", h.CreateDoseNotification)
	doseNotifications.PUT("/:id", h.UpdateDoseNotification)
	doseNotifications.PUT("", h.SaveDoseNotification)

	//dose_notification_times
	doseNotificationTimes := v1.Group("/dose_notification_times")
	doseNotificationTimes.GET("", h.GetDoseNotificationTimes)

	//dose_statuses
	doseStatuses := v1.Group("/dose_statuses")
	doseStatuses.GET("", h.GetDoseStatuses)

	//dose_timings
	doseTimings := v1.Group("/dose_timings")
	doseTimings.GET("", h.GetDoseTimings)

	//dose_patterns
	dosePatterns := v1.Group("/dose_patterns")
	dosePatterns.GET("", h.GetDosePatterns)

	//dose_units
	doseUnits := v1.Group("/dose_units")
	doseUnits.GET("", h.GetDoseUnits)

	//drugs
	drugs := v1.Group("/drugs")
	drugs.GET("", h.GetDrugs)

	//medicines
	medicines := v1.Group("/medicines")
	medicines.GET("", h.GetMedicinesOfTaking)
	medicines.POST("", h.CreateMedicine)
	medicines.PUT("/:id", h.UpdateMedicine)
	medicines.PATCH("/:id/delete", h.FinishedTaking)

	//medicine_side_effects
	medicineSideEffects := v1.Group("/medicine_side_effects")
	medicineSideEffects.GET("", h.GetMedicineSideEffects)

	//profiles
	profiles := v1.Group("/profiles")
	profiles.GET("", h.GetProfile)
	profiles.POST("", h.CreateProfile)
	profiles.PUT("", h.UpdateProfile)

	//symptoms
	symptoms := v1.Group("/symptoms")
	symptoms.GET("/:record_date_start/:record_date_to", h.GetSymptomsBetweenRecordDates)
	symptoms.POST("", h.CreateSymptom)
	symptoms.PUT("/:id", h.UpdateSymptom)
	symptoms.DELETE("/:id", h.DeleteSymptom)

	vitals := v1.Group("/vitals")
	vitals.GET("/:record_date_start/:record_date_to", h.GetVitalsBetweenRecordDates)
	vitals.POST("", h.CreateVital)
	vitals.PUT("", h.UpdateVital)
	vitals.PUT("/delete", h.DeleteVital)
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
			user := entity.NewUser()
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
				karteAPI: &myKarte.API{
					Token: claim.MyKarteToken,
				},
				vitalAPI: vital.NewAPI(user.WelbyID),
			})

			return err
		}
	}
}

//CheckVersionMiddleware はバージョンをチェックしてレスポンスヘッダに結果を追加します。
func CheckVersionMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			header := c.Request().Header
			device := header.Get(configs.Conf.Request.DeviceKey)
			version := header.Get(configs.Conf.Request.VersionKey)

			appVersion := entity.NewAppVersion()
			b, err := appVersion.ExistsUpperVersion(db, device, version)
			if err != nil {
				return err
			}
			c.Response().Header().Set(configs.Conf.Response.UpdateRequiredKey, strconv.FormatBool(b))

			return next(c)
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
