//Package pkg は共通的な処理を持つパッケージです。名前がイマイチなのは認識しています。いつか細分化したいと思っています。
package pkg

import (
	"goexample/configs"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/xerrors"
)

//AppError はアプリで発生するエラーの構造体です。
type AppError struct {
	Err     error  `json:"-"`
	Code    int    `json:"-"`
	Message string `json:"message"`
	ErrCode string `json:"code"`
	frame   xerrors.Frame
}

//Error はAppErrorのErrを文字列として表示します。
func (e *AppError) Error() string {
	return fmt.Sprintf("%+v", e.Err)
}

//Unwrap はXErrorのwrap interfaceを実装したものです。
func (e *AppError) Unwrap() error {
	return e.Err
}

//Format はXErrorのwrap interfaceを実装したものです。
func (e *AppError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

//FormatError はXErrorのwrap interfaceを実装したものです。
func (e *AppError) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	e.frame.Format(p)
	return e.Err
}

//ToResponse はエラー内容をJSON形式で返す際に整形します。
func (e *AppError) ToResponse() map[string]interface{} {
	m := map[string]interface{}{
		"errors": map[string]interface{}{
			"code":          e.ErrCode,
			"error_message": e.Message,
		},
	}
	return m
}

//NewAppError は新しいAppErrorを作成します。
func NewAppError() *AppError {
	return &AppError{}
}

//IsAppError は i(error)がAppErrorであるか判定します。
func IsAppError(i interface{}) bool {
	if i == nil {
		return false
	}
	a := reflect.ValueOf(&AppError{}).Type()
	if reflect.TypeOf(i) == a {
		return true
	}
	return false
}

//NewParameterError はparam errorを作成します。
func NewParameterError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("param error: %w", err),
		Code:    http.StatusBadRequest,
		Message: "不正なパラメータです。",
		ErrCode: "400004",
	}
}

//NewQueryError はquery errorを作成します。
func NewQueryError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("query error: %w", err),
		Code:    http.StatusInternalServerError,
		Message: "サーバ内でエラーが発生しました。",
		ErrCode: "500001",
	}
}

//NewMarshalError はmarshal errorを作成します。
func NewMarshalError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("marshal error: %w", err),
		Code:    http.StatusBadRequest,
		Message: "不正なパラメータです。",
		ErrCode: "500001",
	}
}

//NewUnMarshalError はunmarshal errorを作成します。
func NewUnMarshalError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("unmarshal error: %w", err),
		Code:    http.StatusBadRequest,
		Message: "不正なパラメータです。",
		ErrCode: "500001",
	}
}

//NewResourceNotFoundError はresource not found errorを作成します。
func NewResourceNotFoundError() AppError {
	return AppError{
		Err:     xerrors.New("resource not found error"),
		Code:    http.StatusNotFound,
		Message: "レコードが存在しません。",
		ErrCode: "404001",
	}
}

//NewConflictError はconflict errorを作成します。
func NewConflictError() AppError {
	return AppError{
		Err:     xerrors.New("conflict error"),
		Code:    http.StatusConflict,
		Message: "レコードはすでに更新されています。",
		ErrCode: "409002",
	}
}

//NewInvalidTokenError はinvalid token errorを作成します。
func NewInvalidTokenError() AppError {
	return AppError{
		Err:     xerrors.New("invalid token error"),
		Code:    http.StatusUnauthorized,
		Message: "トークンが誤っています。",
		ErrCode: "409002", //TODO コードを採番
	}
}

//NewTokenExpiredError はtoken expired errorを作成します。
func NewTokenExpiredError() AppError {
	return AppError{
		Err:     xerrors.New("expired error"),
		Code:    http.StatusUnauthorized,
		Message: "トークンの有効期限が切れています。",
		ErrCode: "409002", //TODO コードを採番
	}
}

//NewValidateError は validate errorを作成します。
func NewValidateError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("validate Error: %w", err),
		Code:    http.StatusBadRequest,
		Message: "不正なパラメータです。",
		ErrCode: "400004",
	}
}

//NewInternalServerError はinternal server errorを作成します。
func NewInternalServerError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("internal server error: %w", err),
		Code:    http.StatusInternalServerError,
		Message: "サーバ内でエラーが発生しました。",
		ErrCode: "500001",
	}
}

//NewDuplicateError はduplicate errorを作成します。
func NewDuplicateError() AppError {
	return AppError{
		Err:     xerrors.New("duplicate error"),
		Code:    http.StatusConflict,
		Message: "レコードはすでに更新されています。",
		ErrCode: "409002",
	}
}

//NewMyKarteUnauthorizedError はマイカルテにログインしていないエラーを作成します。
func NewMyKarteUnauthorizedError() AppError {
	return AppError{
		Err:     xerrors.New("Unauthorized error."),
		Code:    http.StatusUnauthorized,
		Message: "マイカルテにログインしていません。",
		ErrCode: "409002",
	}
}

//NewMyKarteAuthorizationError はマイカルテへのログインエラーを作成します。
func NewMyKarteAuthorizationError(err error) AppError {
	return AppError{
		Err:     xerrors.New(fmt.Sprintf("Authorization failed. Error,%s:", err.Error())),
		Code:    http.StatusUnauthorized,
		Message: "マイカルテのログインに失敗しました。",
		ErrCode: "409002",
	}
}

//NewMyKarteRequestError はマイカルテへのリクエストエラーを作成します。
func NewMyKarteRequestError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("request Error: %w", err),
		Code:    http.StatusBadRequest,
		Message: "マイカルテへのリクエストに失敗しました。",
		ErrCode: "400004",
	}
}

//NewMedicineRequestError は薬剤サーバへのリクエストエラーを作成します。
func NewMedicineRequestError(err error) AppError {
	return AppError{
		Err:     xerrors.Errorf("request Error: %w", err),
		Code:    http.StatusBadRequest,
		Message: "薬剤サーバへのリクエストに失敗しました。",
		ErrCode: "400004",
	}
}

//NewNoResourceIDError はリソースIDが指定されていないエラーを作成します。
func NewNoResourceIDError() AppError {
	return AppError{
		Err:     xerrors.Errorf("No resource ID."),
		Code:    http.StatusBadRequest,
		Message: "リソースIDが指定されていません。",
		ErrCode: "400004",
	}
}

//IsDuplicateError は errがduplicate errorであるか判定します。
func IsDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == configs.MYSQL_ER_DUP_ENTRY {
			return true
		}
	}
	return false
}
