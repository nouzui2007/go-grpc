package handler

import (
	"goexample/configs"
	"goexample/model"
	"goexample/pkg"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	jwtClaim struct {
		UserID       int32
		jwt.StandardClaims
	}

	responseLogin struct {
		APIToken string `json:"api_token"`
	}

	//LoginRequest はマイカルテのログインに必要な構造体です。
	LoginRequest struct {
		UserName *string `json:"username" validate:"required"`
		Password *string `json:"password" validate:"required"`
	}
)

//NewLoginRequest はマイカルテログイン構造体に値を詰めます。
func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

//Login はログイン認証して、トークンを返します。
func (h *Handler) Login(c echo.Context) error {
	r := NewLoginRequest()
	if err := c.Bind(r); err != nil {
		er := pkg.NewParameterError(err)
		c.Echo().Logger.Error(err)
		return c.JSON(er.Code, &er)
	}

	// get user
	user := model.NewUser()
	if err := user.FindByName(h.db, r.UserName); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// トークン発行
	user.Expired = pkg.TimeToTPtr(time.Now().Add(time.Duration(configs.Conf.Authorization.Expired) * time.Second))
	if er := user.Update(h.db); er != nil {
		c.Echo().Logger.Error(er)
		return c.JSON(http.StatusBadRequest, er)
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwtClaim{
		UserID:       *user.ID,
	})
	// Secretで文字列にする. このSecretはサーバだけが知っている
	tokenstring, err := token.SignedString([]byte(os.Getenv(configs.Conf.Authorization.Secret)))
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusForbidden, &err)
	}

	response := responseLogin{APIToken: tokenstring}
	return c.JSON(http.StatusOK, response)
}

//Logout は認証済みのユーザーの認証期間を空にする。
//
//認証期間を空にすることで、トークンの有効期間をなくす。
func (h *Handler) Logout(c echo.Context) error {
	cc := c.(*Context)

	cc.userInfo.Expired = nil
	if er := cc.userInfo.ClearExpired(h.db); er != nil {
		c.Echo().Logger.Error(er)
		return c.JSON(er.(*pkg.AppError).Code, er)
	}
	return c.NoContent(http.StatusOK)
}

// CreateUser はユーザを生成します。この段階では仮登録です。
func (h *Handler) CreateUser(c echo.Context) error {
	r := NewLoginRequest()
	if err := c.Bind(r); err != nil {
		er := pkg.NewParameterError(err)
		c.Echo().Logger.Error(err)
		return c.JSON(er.Code, &er)
	}

	user := model.NewUser()
	user.Name = r.UserName
	user.Password = r.Password

	if err := user.Create(h.db); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser はユーザを更新します。本登録に利用します。
func (h *Handler) UpdateUser(c echo.Context) error {
	r := NewLoginRequest()
	if err := c.Bind(r); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	user := model.NewUser()
	if err := user.FindByName(h.db, r.UserName); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	user.Password = r.Password
	if err := user.Update(h.db); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	return c.NoContent(http.StatusOK)
}