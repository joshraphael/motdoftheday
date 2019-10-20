package apierror

import (
	"errors"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type IApiError interface {
	Error() string
	Status() string
	Code() int
}

const (
	MethodHTTP string = "HTTP"
	MethodGRPC string = "GRPC"
)

type ApiError struct {
	IApiError
	validator *validator.Validate
	err       error  `validate:"required"`
	status    string `validate:"required"`
	method    string `validate:"required"`
}

func New(e error, status string, m string) IApiError {
	apiErr := &ApiError{
		validator: validator.New(),
		err:       e,
		status:    status,
		method:    m,
	}
	if err := apiErr.validate(); err != nil {
		msg := "error validating ApiError '" + status + "' choosing default INTERNAL error using HTTP method: " + err.Error()
		return &ApiError{
			validator: validator.New(),
			err:       errors.New(msg),
			status:    "INTERNAL",
			method:    "HTTP",
		}
	}
	return apiErr
}

func (ae ApiError) validate() error {
	if ae.validator == nil {
		return errors.New("no validator in ApiError")
	}
	err := ae.validator.Struct(ae)
	if err != nil {
		return errors.New("could not validate ApiError: " + err.Error())
	}
	validStatus := ae.statusList()
	if _, ok := validStatus[ae.status]; !ok {
		return errors.New("status does not exist: " + ae.status)
	}
	if _, ok := validStatus[ae.status][ae.method]; !ok {
		return errors.New("method does not exist for status '" + ae.status + "': " + ae.status)
	}
	return nil
}

func (ae ApiError) statusList() map[string]map[string]int {
	return map[string]map[string]int{
		"OK": map[string]int{
			"HTTP": http.StatusOK,
			"GRPC": 0,
		},
		"BAD_REQUEST": map[string]int{
			"HTTP": http.StatusBadRequest,
			"GRPC": 3,
		},
		"NOT_FOUND": map[string]int{
			"HTTP": http.StatusNotFound,
			"GRPC": 3,
		},
		"INTERNAL": map[string]int{
			"HTTP": http.StatusInternalServerError,
			"GRPC": 13,
		},
	}
}

func (ae ApiError) Error() string {
	return ae.err.Error()
}

func (ae ApiError) Status() string {
	return ae.status
}

func (ae ApiError) Code() int {
	validStatus := ae.statusList()
	return validStatus[ae.status][ae.method]
}
