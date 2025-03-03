package utils

import (
	"net/http"
)

type ErrorCode string

type ErrorDetails struct {
	Message    string
	StatusCode int
}

type Error struct {
	Code    ErrorCode
	Message string
}

const (
	ErrInvalidRequest  ErrorCode = "INVALID_REQUEST"
	ErrValidationError ErrorCode = "VALIDATION_ERROR"

	ErrInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"

	ErrUserNotFound        ErrorCode = "USER_NOT_FOUND"
	ErrTransactionNotFound ErrorCode = "TRANSACTION_NOT_FOUND"
	ErrTransactionSameUser ErrorCode = "TRANSACTION_SAME_USER"
	ErrWalletNotFound      ErrorCode = "WALLET_NOT_FOUND"
	ErrWalletInactive      ErrorCode = "WALLET_INACTIVE"
	ErrInsufficientBalance ErrorCode = "INSUFFICIENT_BALANCE"
	ErrUserAlreadyExists   ErrorCode = "USER_ALREADY_EXISTS"
)

var errorMessages = map[ErrorCode]ErrorDetails{
	ErrInvalidRequest: {
		Message:    "Invalid Request",
		StatusCode: http.StatusBadRequest,
	},
	ErrUserNotFound: {
		Message:    "User Not Found",
		StatusCode: http.StatusNotFound,
	},
	ErrTransactionNotFound: {
		Message:    "Transaction Not Found",
		StatusCode: http.StatusNotFound,
	},
	ErrTransactionSameUser: {
		Message:    "Sender And Receiver Cannot Be The Same User",
		StatusCode: http.StatusBadRequest,
	},
	ErrWalletNotFound: {
		Message:    "Wallet Not Found",
		StatusCode: http.StatusNotFound,
	},
	ErrWalletInactive: {
		Message:    "Wallet Is Inactive, Balance Cannot Be Updated",
		StatusCode: http.StatusBadRequest,
	},
	ErrInsufficientBalance: {
		Message:    "Insufficient Balance",
		StatusCode: http.StatusBadRequest,
	},
	ErrUserAlreadyExists: {
		Message:    "User Already Exists",
		StatusCode: http.StatusBadRequest,
	},
	ErrValidationError: {
		Message:    "Validation Error",
		StatusCode: http.StatusBadRequest,
	},
	ErrInternalServerError: {
		Message:    "Internal Server Error",
		StatusCode: http.StatusInternalServerError,
	},
}

func (e *Error) Error() string {
	return e.Message
}

func IsError(err error, code ErrorCode) bool {
	return err != nil && err.(*Error).Code == code
}

func NewError(code ErrorCode) *Error {
	return &Error{Code: code, Message: errorMessages[code].Message}
}

func NewErrorWithMessage(code ErrorCode, message string) *Error {
	return &Error{Code: code, Message: message}
}

func GetErrorDetails(err *Error) ErrorDetails {
	errorDetails, ok := errorMessages[err.Code]
	if !ok {
		return ErrorDetails{
			Message:    "Unknown Error",
			StatusCode: http.StatusInternalServerError,
		}
	}
	errorDetails.Message = err.Message
	return errorDetails
}
