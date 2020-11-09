package handler

import (
	"goexample/configs"
	"github.com/labstack/echo"
)

// Error はレスポンス用Error構造体です。(全てpkg.Errorで行いたいため削除予定)
type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

//TODO :create DBConnectionError

// NewError は新しいErrorを作成します。(全てpkg.Errorで行いたいため削除予定)
func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["body"] = v.Message
	default:
		e.Errors["body"] = v.Error()
	}
	return e
}

// NewMiddlewareValidatorError はMiddleWareによる新しいバリデーションによるErrorを作成します。(全てpkg.Errorで行いたいため削除予定)
// return {"error":{code:"000",error_error_message:"xxx"}} and output log
// return BadRequestError
func NewMiddlewareValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})

	mConf := configs.Conf.Message
	e.Errors["code"] = configs.CodeBadRequest
	e.Errors["error_message"] = mConf.E0001

	return e
}

// NewModelValidatorError Modelによる新しいバリデーションによるErrorを作成します。(全てpkg.Errorで行いたいため削除予定)
//return {"error":{code:"000",error_error_message:"xxx"}} and output log
// return BadRequestError
func NewModelValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})

	mConf := configs.Conf.Message
	e.Errors["code"] = configs.CodeBadRequest
	e.Errors["error_message"] = mConf.E0001

	return e
}
