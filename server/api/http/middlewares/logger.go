package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"portfolio/api/http/utils"
	"portfolio/logger"
	"portfolio/shared"
	"regexp"
	"strings"
	"time"
)

func maskSensitiveData(jsonStr string) string {
	sensitiveFields := []string{
		"password", "Password", "PASSWORD",
		"current_password", "currentPassword", "current_Password",
		"new_password", "newPassword", "new_Password",
		"confirm_password", "confirmPassword", "confirm_Password",
		"old_password", "oldPassword", "old_Password",
		"passwd", "pwd", "secret", "token", "api_key", "apiKey",
		"private_key", "privateKey", "access_token", "accessToken",
		"refresh_token", "refreshToken", "auth_token", "authToken",
	}

	result := jsonStr
	for _, field := range sensitiveFields {
		pattern := `("` + field + `"\s*:\s*")([^"]*)(")`
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, "${1}xxxxxxxxx${3}")

		pattern2 := `(` + field + `=)([^&\s]*)`
		re2 := regexp.MustCompile(pattern2)
		result = re2.ReplaceAllString(result, "${1}xxxxxxxxx")
	}

	return result
}

func formatHeaders(headers http.Header) string {
	sensitiveHeaders := map[string]bool{
		"authorization":  true,
		"cookie":         true,
		"x-api-key":      true,
		"x-auth-token":   true,
		"x-access-token": true,
	}

	var headerPairs []string
	for name, values := range headers {
		lowerName := strings.ToLower(name)
		value := strings.Join(values, ",")

		if sensitiveHeaders[lowerName] {
			if value != "" {
				value = "xxxxxxxxx"
			}
		}

		if len(value) > 100 {
			value = value[:100] + "..."
		}

		headerPairs = append(headerPairs, name+": "+value)
	}

	return fmt.Sprintf("[%s]", strings.Join(headerPairs, " | "))
}

func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := utils.GenerateRequestID()
			ctx := context.WithValue(r.Context(), shared.REQUEST_ID_KEY, requestID)
			r = r.WithContext(ctx)

			var requestBody string
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					requestBody = string(bodyBytes)
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			wrapper := &responseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
			}

			next.ServeHTTP(wrapper, r)

			duration := time.Since(start)

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
				scheme = forwardedProto
			}

			headersForLog := formatHeaders(r.Header)

			requestBodyForLog := ""
			if requestBody != "" {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(requestBody), &jsonData); err == nil {
					if compactBytes, err := json.Marshal(jsonData); err == nil {
						requestBodyForLog = string(compactBytes)
					} else {
						requestBodyForLog = requestBody
					}
				} else {
					requestBodyForLog = requestBody
				}
				requestBodyForLog = maskSensitiveData(requestBodyForLog)
				if len(requestBodyForLog) > 1000 {
					requestBodyForLog = requestBodyForLog[:1000] + "..."
				}
			}

			responseBodyForLog := ""
			responseBody := wrapper.body.String()
			if responseBody != "" && strings.Contains(w.Header().Get("Content-Type"), "application/json") {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(responseBody), &jsonData); err == nil {
					if compactBytes, err := json.Marshal(jsonData); err == nil {
						responseBodyForLog = string(compactBytes)
					} else {
						responseBodyForLog = responseBody
					}
				} else {
					responseBodyForLog = responseBody
				}
				responseBodyForLog = maskSensitiveData(responseBodyForLog)
				if len(responseBodyForLog) > 1000 {
					responseBodyForLog = responseBodyForLog[:1000] + "..."
				}
			}
			// godump.Dump(r)
			loggerPrefix := "HTTP"
			if requestBodyForLog != "" && responseBodyForLog != "" {
				logger.Info("[%s] request_id=%s method=%s uri=%s status=%d duration=%v remote_addr=%s headers=%s request_body=%s response_body=%s",
					loggerPrefix, requestID, r.Method, fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI), wrapper.statusCode, duration, r.RemoteAddr, headersForLog, requestBodyForLog, responseBodyForLog)
			} else if requestBodyForLog != "" {
				logger.Info("[%s] request_id=%s method=%s uri=%s status=%d duration=%v remote_addr=%s headers=%s request_body=%s",
					loggerPrefix, requestID, r.Method, fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI), wrapper.statusCode, duration, r.RemoteAddr, headersForLog, requestBodyForLog)
			} else if responseBodyForLog != "" {
				logger.Info("[%s] request_id=%s method=%s uri=%s status=%d duration=%v remote_addr=%s headers=%s response_body=%s",
					loggerPrefix, requestID, r.Method, fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI), wrapper.statusCode, duration, r.RemoteAddr, headersForLog, responseBodyForLog)
			} else {
				logger.Info("[%s] request_id=%s method=%s uri=%s status=%d duration=%v remote_addr=%s headers=%s",
					loggerPrefix, requestID, r.Method, fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI), wrapper.statusCode, duration, r.RemoteAddr, headersForLog)
			}
		})
	}
}

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *responseWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWrapper) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}
