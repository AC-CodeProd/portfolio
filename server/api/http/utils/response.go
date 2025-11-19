package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portfolio/domain"
	"portfolio/shared"
	"time"
)

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to encode response: %v\n", err)
		WriteErrorResponse(w, domain.NewInternalError("Failed to encode response", err))
		return
	}
	if _, err := w.Write(dataBytes); err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, errors ...error) {
	if len(errors) == 0 {
		errors = []error{domain.NewInternalError("No error provided", nil)}
	}

	httpStatus := http.StatusInternalServerError

	var apiErrors []*shared.APIError

	for _, err := range errors {
		if domainErr, ok := domain.AsDomainError(err); ok {
			apiError := DomainErrorToAPIError(domainErr)
			apiErrors = append(apiErrors, apiError)
			if domainErr.HTTPStatus() > httpStatus || httpStatus == http.StatusInternalServerError {
				httpStatus = domainErr.HTTPStatus()
			}
		} else {
			domainErr := domain.NewInternalError(err.Error(), err)
			apiError := DomainErrorToAPIError(domainErr)
			apiErrors = append(apiErrors, apiError)
		}
	}

	if len(apiErrors) > 1 && httpStatus == http.StatusInternalServerError {
		httpStatus = http.StatusBadRequest
	}

	response := shared.APIResponse{
		Errors: apiErrors,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("failed to encode response: %v\n", err)
	}
}

func StringsToErrors(errors []string) []error {
	errs := make([]error, len(errors))
	for i, errStr := range errors {
		errs[i] = domain.NewValidationError("validation_error", errStr, nil)
	}
	return errs
}

func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	JSONResponse(w, statusCode, data)
}

func WriteDeleteSuccessResponse(w http.ResponseWriter, message string, requestID string) {
	response := shared.APIResponse{
		Message: message,
		Meta: shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": requestID,
		},
	}
	JSONResponse(w, http.StatusOK, response)
}
