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
