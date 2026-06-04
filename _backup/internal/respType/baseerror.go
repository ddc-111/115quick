package respType

import "strconv"

const defaultCode = 1001   //通用错误码
const DBError = 1002       //数据库错误码
const ParamError = 1003    //参数错误码
const NotFoundError = 1004 //未找到错误码
const AuthError = 1005     //权限错误码
const OtherError = 1006    //其他错误码

type CodeError struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}
func NewDBError(msg string) error {
	return NewCodeError(DBError, msg)
}
func NewParamError(msg string) error {
	return NewCodeError(ParamError, msg)
}
func NewNotFoundError(msg string) error {
	return NewCodeError(NotFoundError, msg)
}
func NewAuthError(msg string) error {
	return NewCodeError(AuthError, msg)
}
func NewOtherError(msg string) error {
	return NewCodeError(OtherError, msg)
}

func NewDefaultErrorByErrorCode(code int) error {

	return NewCodeError(defaultCode, strconv.Itoa(code))
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
