package handler

import (
	"goexample/model/entity"
	"goexample/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCheckVersionMiddleware(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	h := CheckVersionMiddleware(db)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	createVersion(db, map[string]string{
		"iOS":     "1.0.0",
		"Android": "1.0.0",
	})

	t.Run("バージョンアップなし(iOS)", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "iOS", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "false", rec.HeaderMap["Update-Required"][0])
		}
	})
	t.Run("バージョンアップなし(Android)", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "Android", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "false", rec.HeaderMap["Update-Required"][0])
		}
	})

	createVersion(db, map[string]string{
		"iOS":     "1.0.1",
		"Android": "1.0.1",
	})
	t.Run("バージョンアップあり(iOS)", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "iOS", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "true", rec.HeaderMap["Update-Required"][0])
		}
	})
	t.Run("バージョンアップあり(ios)", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "ios", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "true", rec.HeaderMap["Update-Required"][0])
		}
	})
	t.Run("バージョンアップあり(Android)", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "Android", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "true", rec.HeaderMap["Update-Required"][0])
		}
	})
	t.Run("バージョンアップあり(android)", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "android", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "true", rec.HeaderMap["Update-Required"][0])
		}
	})

	t.Run("ヘッダなし", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", nil)
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "false", rec.HeaderMap["Update-Required"][0])
		}
	})
	t.Run("Android/iOS以外", func(t *testing.T) {
		c, rec := createContextWithHeader(http.MethodGet, "/", map[string]string{"Device": "Windows", "Version": "1.0.0"})
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "false", rec.HeaderMap["Update-Required"][0])
		}
	})

}

func createContextWithHeader(method string, uri string, headers map[string]string) (*Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, uri, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	u := entity.User{}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return &Context{
		Context:  c,
		userInfo: u,
	}, rec
}

func createVersion(db *gorm.DB, rec map[string]string) {
	for device, version := range rec {
		appVersion := entity.AppVersion{
			Device:    &device,
			Version:   &version,
			BaseModel: entity.NewBaseModel(),
		}
		db.Create(&appVersion)
	}
}
