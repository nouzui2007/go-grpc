package handler

import (
	"goexample/model/entity"
	"goexample/model/myKarte"
	"goexample/pkg"
	"goexample/test"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"

	"github.com/stretchr/testify/assert"
)

type returnMockValue int

const (
	returnNull returnMockValue = iota
	returnRequestError
	returnAuthError
	returnValidateError
)

var mode returnMockValue

func (m *KarteAPIMock) Login(request *myKarte.LoginRequest) error {
	switch mode {
	case returnRequestError:
		err := pkg.NewMyKarteRequestError(xerrors.Errorf("mock"))
		return &err
	case returnAuthError:
		err := pkg.NewMyKarteUnauthorizedError()
		return &err
	case returnValidateError:
		err := pkg.NewValidateError(xerrors.Errorf("mock"))
		return &err
	default:
		return nil
	}
}

func (m *KarteAPIMock) Logout() error {
	switch mode {
	case returnRequestError:
		err := pkg.NewMyKarteRequestError(xerrors.Errorf("mock"))
		return &err
	case returnAuthError:
		err := pkg.NewMyKarteUnauthorizedError()
		return &err
	default:
		return nil
	}
}

func (m *KarteAPIMock) GetProfile() (*myKarte.Profile, error) {
	switch mode {

	default:
		return &myKarte.Profile{
			ID: pkg.IntToIPtr(1),
			User: &myKarte.User{
				ID:          pkg.IntToIPtr(1),
				WelbyUserID: pkg.IntToIPtr(100),
			},
		}, nil
	}
}

func TestHandler_Login(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	h := NewHandler(db)
	// u := createUser(db)

	t.Run("usernameなし", func(t *testing.T) {
		mode = returnValidateError
		var json = `{"password":"pass"}`
		c, rec := createContext(http.MethodPost, "/", strings.NewReader(json))
		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("passwordなし", func(t *testing.T) {
		mode = returnValidateError
		var json = `{"username":"user"}`
		c, rec := createContext(http.MethodPost, "/", strings.NewReader(json))
		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ログイン可能", func(t *testing.T) {
		mode = returnNull
		var json = `{"username":"user", "password":"pass"}`
		c, rec := createContext(http.MethodPost, "/", strings.NewReader(json))
		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestHandler_Logout(t *testing.T) {
	db, teardown := test.GetTestDBConn()
	defer teardown()

	h := NewHandler(db)
	u := createUser(db)

	t.Run("正常にログアウト", func(t *testing.T) {
		mode = returnNull
		c, rec := createContext(http.MethodPost, "/", nil)
		c.userInfo = u
		if assert.NoError(t, h.Logout(c)) {
			assert.Equal(t, http.StatusOK, rec.Code, "HTTPステータスコード")
			assert.Nil(t, c.userInfo.Expired, "有効期限")
		}
	})

	t.Run("リクエストエラー", func(t *testing.T) {
		mode = returnRequestError
		c, rec := createContext(http.MethodPost, "/", nil)
		c.userInfo = u
		if assert.NoError(t, h.Logout(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "HTTPステータスコード")
		}
	})

	t.Run("MyKarteAPIがない", func(t *testing.T) {
		c, rec := createContext(http.MethodPost, "/", nil)
		c.userInfo.WelbyID = pkg.IntToIPtr(1)
		c.karteAPI = nil
		if assert.NoError(t, h.Logout(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code, "HTTPステータスコード")
		}
	})
}

func createUser(db *gorm.DB) entity.User {
	user := entity.User{
		WelbyID:        pkg.IntToIPtr(100),
		KartePatientID: pkg.IntToIPtr(101),
		Expired:        pkg.TimeToTPtr(time.Now()),
		BaseModel: entity.BaseModel{
			ID:        pkg.IntToIPtr(1),
			CreatedAt: pkg.TimeToTPtr(time.Now()),
			UpdatedAt: pkg.TimeToTPtr(time.Now()),
			Deleted:   pkg.IntToIPtr(0),
		},
	}
	db.Create(&user)
	return user
}
