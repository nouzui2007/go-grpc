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
		MyKarteToken myKarte.Token
		jwt.StandardClaims
	}

	responseLogin struct {
		APIToken string `json:"api_token"`
	}
)

//Login はログイン認証して、トークンを返します。
func (h *Handler) Login(c echo.Context) error {
	r := myKarte.NewLoginRequest()
	if err := c.Bind(r); err != nil {
		er := pkg.NewParameterError(err)
		c.Echo().Logger.Error(err)
		return c.JSON(er.Code, &er)
	}

	// マイカルテへログインして、プロファイルを取得する
	karte := myKarte.NewAPI()
	if err := karte.Login(r); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}
	// マイカルテ /patient
	profile, errp := karte.GetProfile()
	if errp != nil {
		c.Echo().Logger.Error(errp)
		return c.JSON(errp.(*pkg.AppError).Code, errp)
	}

	welbyID := profile.User.WelbyUserID
	kartePatientID := profile.ID

	// DSAF get user
	userObject := entity.NewUser()
	user, err := userObject.FindByWelbyID(h.db, welbyID)
	if err != nil {
		// create user if DSAF user not found
		if err.(*pkg.AppError).Code == http.StatusNotFound {
			c.Logger().Debug("User data is not found.")
			user = entity.NewUser()
			user.WelbyID = welbyID
			user.KartePatientID = kartePatientID
			if er := user.Create(h.db); er != nil {
				c.Echo().Logger.Error(err)
				return c.JSON(er.(*pkg.AppError).Code, er)
			}
		} else {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
	} else if (user.KartePatientID == nil) && (kartePatientID != nil) {
		c.Logger().Debug("Karte Patient IDの更新")
		user.KartePatientID = kartePatientID
		if er := user.Update(h.db); er != nil {
			c.Echo().Logger.Error(er)
			return c.JSON(er.(*pkg.AppError).Code, er)
		}
	}

	// トークン発行
	user.Expired = pkg.TimeToTPtr(time.Now().Add(time.Duration(configs.Conf.Authorization.Expired) * time.Second))
	if er := user.Update(h.db); er != nil {
		c.Echo().Logger.Error(er)
		return c.JSON(http.StatusBadRequest, er)
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwtClaim{
		UserID:       *user.ID,
		MyKarteToken: karte.Token,
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

	// マイカルテからログアウト
	if cc.karteAPI == nil {
		err := pkg.NewMyKarteUnauthorizedError()
		c.Echo().Logger.Error(err)
		return c.JSON(err.Code, &err)
	}

	if err := cc.karteAPI.Logout(); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	cc.userInfo.Expired = nil
	if er := cc.userInfo.ClearExpired(h.db); er != nil {
		c.Echo().Logger.Error(er)
		return c.JSON(er.(*pkg.AppError).Code, er)
	}
	return c.NoContent(http.StatusOK)
}

// CreateUser はユーザを生成します。この段階では仮登録です。
func (h *Handler) CreateUser(c echo.Context) error {
	r := myKarte.NewTemporaryRegistration()
	if err := c.Bind(r); err != nil {
		er := pkg.NewParameterError(err)
		c.Echo().Logger.Error(err)
		return c.JSON(er.Code, &er)
	}

	// マイカルテ仮登録
	karte := myKarte.NewAPI()
	result, err := karte.RegistTemporary(r)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	// 正常終了したら、dsafユーザ作成
	user := entity.NewUser()
	user.WelbyID = result.WelbyUserID

	if err := user.Create(h.db); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser はユーザを更新します。本登録に利用します。
func (h *Handler) UpdateUser(c echo.Context) error {
	r := myKarte.NewFullRegistration()
	if err := c.Bind(r); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// マイカルテ 本登録
	karte := myKarte.NewAPI()
	if err := karte.RegistFull(r); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	return c.NoContent(http.StatusOK)
}

//ConfirmationCode は認証コードの再発行を依頼します。
func (h *Handler) ConfirmationCode(c echo.Context) error {
	r := myKarte.NewRegistrationKey()
	if err := c.Bind(r); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// マイカルテ 本登録
	karte := myKarte.NewAPI()
	if err := karte.ConfirmationCode(r); err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(err.(*pkg.AppError).Code, err)
	}

	return c.NoContent(http.StatusOK)
}
