package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func BindAndValidateRequest(c *gin.Context, request interface{}) error {
	err := c.ShouldBindJSON(request)
	if err != nil {
		return NewErrorWithMessage(ErrInvalidRequest, err.Error())
	}
	err = ValidateStruct(request)
	if err != nil {
		return err
	}
	return nil
}

func ValidateStruct(data interface{}) error {
	if err := validate.Struct(data); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return NewErrorWithMessage(ErrInvalidRequest, err.Error())
		}

		// For create operations, enforce all validations
		for _, err := range err.(validator.ValidationErrors) {
			// Check error type and return appropriate error
			if err.Tag() == "required" {
				return NewErrorWithMessage(ErrValidationError, fmt.Sprintf("%s is required", err.Field()))
			}
			if err.Tag() == "nefield" {
				return NewErrorWithMessage(ErrValidationError, fmt.Sprintf("%s and %s cannot be the same", err.Field(), err.Param()))
			}

			return NewErrorWithMessage(ErrValidationError, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return nil
}

func BindResponse(c *gin.Context, data interface{}, err error) error {
	if err != nil {
		return NewErrorWithMessage(ErrInvalidRequest, err.Error())
	}
	err = ValidateStruct(data)
	if err != nil {
		return err
	}
	return nil
}

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func ResponseError(c *gin.Context, err error) {
	if err, ok := err.(*Error); ok {
		errorDetails := GetErrorDetails(err)
		c.JSON(errorDetails.StatusCode, ErrorResponse{
			Code:    err.Code,
			Message: err.Message,
		})
	} else {
		response := GetErrorDetails(NewError(ErrInternalServerError))
		c.JSON(response.StatusCode, ErrorResponse{
			Code:    ErrInternalServerError,
			Message: response.Message,
		})
	}
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
