package utils

import (
	"portfolio/domain"
	"portfolio/shared"
	"time"
)

func DomainErrorToAPIError(domainErr *domain.DomainError) *shared.APIError {
	meta := shared.Meta{
		"timestamp": time.Now().Format(time.RFC3339),
		"code":      domainErr.Code,
	}

	if domainErr.Field != "" {
		meta["field"] = domainErr.Field
	}

	if domainErr.Details != nil {
		for k, v := range domainErr.Details {
			meta[k] = v
		}
	}

	if domainErr.Cause != nil {
		meta["cause"] = domainErr.Cause.Error()
	}

	return &shared.APIError{
		Status: domainErr.HTTPStatus(),
		Title:  domainErr.Title(),
		Detail: domainErr.Message,
		Meta:   meta,
	}
}

func DomainErrorsToAPIErrors(domainErrs []*domain.DomainError) []*shared.APIError {
	apiErrors := make([]*shared.APIError, len(domainErrs))
	for i, domainErr := range domainErrs {
		apiErrors[i] = DomainErrorToAPIError(domainErr)
	}
	return apiErrors
}

func ErrorToAPIError(err error, defaultStatus int) *shared.APIError {
	if domainErr, ok := domain.AsDomainError(err); ok {
		return DomainErrorToAPIError(domainErr)
	}

	return &shared.APIError{
		Status: defaultStatus,
		Title:  "Internal Server Error",
		Detail: err.Error(),
		Meta: shared.Meta{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}
