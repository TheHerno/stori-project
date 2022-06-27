package errors

import (
	"net/http"
	"stori-service/src/libs/i18n"
)

var (
	//ErrInternalServer indicates an internal server error, this is used a default error when it's a general error
	ErrInternalServer = NewMyError(http.StatusInternalServerError, i18n.Message{MessageID: "ERRORS.INTERNAL_SERVER"})

	//ErrIDNotNumeric indicates a numeric ID wasn't send as numeric
	ErrIDNotNumeric = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.ID_NOT_NUMERIC"})

	//ErrConnectionProvider indicates comunication error with provider
	ErrConnectionProvider = NewMyError(http.StatusServiceUnavailable, i18n.Message{MessageID: "ERRORS.CONNECTION_PROVIDER"})

	//ErrPageSizeTooHigh indicates that pageSize param is higher than the valid number
	ErrPageSizeTooHigh = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.PAGE_SIZE_TOO_LARGE"})

	//ErrPageTooHigh indicates that page param is higher than the valid number
	ErrPageTooHigh = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.PAGE_TOO_LARGE"})

	//ErrURLNotFound indicates a requested URL doesn't exist
	ErrURLNotFound = NewMyError(http.StatusNotFound, i18n.Message{MessageID: "ERRORS.URL_NOT_FOUND"})

	//ErrNotFound indicates an entity not found error
	ErrNotFound = NewMyError(http.StatusNotFound, i18n.Message{MessageID: "ERRORS.NOT_FOUND"})

	//ErrInvalidFileLine indicates a line in the file is invalid
	ErrInvalidFileLine = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.INVALID_FILE_LINE"})
)

//Private errors
var (
	errUnsupportedFieldValue = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.UNSUPPORTED_FIELD_VALUE"})
	errFieldValidation       = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.FIELD_VALIDATION"})
)
