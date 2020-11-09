package handler

import (
	"goexample/model/entity"
	"goexample/model/myKarte"
	"goexample/test"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMain(m *testing.M) {
	//DB接続 終了時にクローズ
	shutdown := test.SetupDBConn()
	defer shutdown()

	m.Run()
}

func createContext(method string, uri string, reader io.Reader) (*Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, uri, reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	u := entity.User{}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return &Context{
		Context:  c,
		userInfo: u,
	}, rec
}

func addRequestHeader(request *http.Request, key, value string) {
	request.Header.Set(key, value)
}
