package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/logger"
	"portfolio/shared"
	"strings"
	"time"
)

type ResponseWrapper struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func (rw *ResponseWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *ResponseWrapper) Write(b []byte) (int, error) {
	return rw.buf.Write(b)
}

func ResponseMiddleware(logger *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &ResponseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rw, r)
			response := shared.APIResponse{}
			if rw.statusCode >= 400 {
				value, valid := parseAPIError(rw.buf.Bytes())

				if completeResponse, isComplete := value.(struct {
					Errors []*shared.APIError `json:"errors"`
				}); isComplete {
					for k, v := range rw.Header() {
						w.Header().Set(k, strings.Join(v, ","))
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(rw.statusCode)
					err := json.NewEncoder(w).Encode(completeResponse)
					if err != nil {
						logger.Printf("failed to encode response: %v", err)
					}
					return
				}

				var errors []*shared.APIError
				if valid {
					switch v := value.(type) {
					case *shared.APIError:
						errors = []*shared.APIError{v}
					case []*shared.APIError:
						errors = v
					}
				} else {
					errors = []*shared.APIError{
						{
							Status: rw.statusCode,
							Title:  http.StatusText(rw.statusCode),
							Detail: strings.TrimSpace(rw.buf.String()),
							Meta:   shared.Meta{"timestamp": time.Now().Format(time.RFC3339)},
						},
					}
				}
				response.Errors = errors
			} else {
				var data interface{}
				if err := json.Unmarshal(rw.buf.Bytes(), &data); err == nil {
					response.Data = data
				}
			}

			for k, v := range rw.Header() {
				w.Header().Set(k, strings.Join(v, ","))
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(rw.statusCode)
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				logger.Printf("failed to encode response: %v", err)
			}
		})
	}
}

func parseAPIError(b []byte) (any, bool) {
	var jsonAPIResponse struct {
		Errors []*shared.APIError `json:"errors"`
	}
	if err := json.Unmarshal(b, &jsonAPIResponse); err == nil && len(jsonAPIResponse.Errors) > 0 {
		return jsonAPIResponse, true
	}

	var domainErr domain.DomainError
	if err := json.Unmarshal(b, &domainErr); err == nil && domainErr.Code != "" {
		return utils.DomainErrorToAPIError(&domainErr), true
	}

	var domainErrs []*domain.DomainError
	if err := json.Unmarshal(b, &domainErrs); err == nil && len(domainErrs) > 0 {
		return utils.DomainErrorsToAPIErrors(domainErrs), true
	}

	var single shared.APIError
	if err := json.Unmarshal(b, &single); err == nil && single.Status != 0 {
		return &single, true
	}

	var multiple []*shared.APIError
	if err := json.Unmarshal(b, &multiple); err == nil && len(multiple) > 0 {
		return multiple, true
	}

	return string(b), false
}
