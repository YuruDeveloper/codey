package appError

import (
	"errors"
	"fmt"
)

type ErrorType int 

const (
	JsonMarshalError = ErrorType(iota)
	JsonUnMarshalError
	FailFindHomeDir
	FailReadFile
	FailMakeFolder
	FailMakeFile
	FailLoadJsonData
	FailRefreshToken
	UnexpectedCode
	FailMakeRequest
	FailGetResponse
	HttpNotOK
	FailDecodeHttpBody
	FailOpenBrowser
	UnexpectedOS
	FailUpdateToken
	UnexpectedAuthIndex
)

var _ error = (*Error)(nil)

type Error struct {
	errorType ErrorType
	err error
}

func NewError(errorType ErrorType ,err error) error{
	return &Error {
		errorType: errorType,
		err: err,
	}
}

func NewValidError(errorType ErrorType,message string) error {
	return &Error{
		errorType: errorType,
		err: errors.New(message),
	}
}

func (instance *Error) String() string {
	switch instance.errorType {
		case JsonMarshalError:
			return "JsonMarshalError"
		case JsonUnMarshalError:
			return "JsonUnMarshalError"
		case FailFindHomeDir:
			return  "FailFindHomeDir"
		case FailReadFile:
			return "FailReadFile"
		case FailMakeFolder:
			return "FailMakeFolder"
		case FailMakeFile:
			return "FailMakeFile"
		case FailLoadJsonData:
			return "FailLoadJsonData"
		case FailRefreshToken:
			return "FailRefreshToken"
		case UnexpectedCode:
			return "UnexpectedCode"
		case FailMakeRequest:
			return "FailMakeRequest"
		case FailGetResponse:
			return "FailGetResponse"
		case HttpNotOK:
			return "HttpNotOK"
		case FailDecodeHttpBody:
			return "FailDecodeHttpBody"
		case FailOpenBrowser:
			return "FailOpenBrowser"
		case UnexpectedOS:
			return "UnexpectedOS"
		case FailUpdateToken:
			return "FailUpdateToken"
		case UnexpectedAuthIndex:
			return "UnexpectedAuthIndex"
		default:
			return "Unknown"
	}
}

func (instance *Error) Error() string {
	return fmt.Sprintf("%s : %s",instance.String(),instance.err.Error())
}

func (instance *Error) Unwrap() error {
	return instance.err
}

func (instance *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.errorType == instance.errorType
}