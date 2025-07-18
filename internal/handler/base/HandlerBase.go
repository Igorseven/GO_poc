package base

import (
	notify "PocGo/internal/domain/notification"
	setJson "encoding/json"
	httpclient "net/http"
)

const (
	ContentTypeKey   = "Content-Type"
	ContentTypeValue = "application/json"
	AcceptKey        = "Accept"
)

// ValidateHTTPMethod validates if the request method matches the expected method
// Returns an error if the method is not allowed
func ValidateHTTPMethod(responseWriter httpclient.ResponseWriter, request *httpclient.Request, method string) error {
	if request.Method != method {
		httpclient.Error(responseWriter, notify.InvalidMethod, httpclient.StatusMethodNotAllowed)
		return notify.CreateNotification(notify.InvalidMethod)
	}

	return nil
}

// ValidationGetMethod validates if the request method is GET
// This is kept for backward compatibility
func ValidationGetMethod(responseWriter httpclient.ResponseWriter, request *httpclient.Request) error {
	return ValidateHTTPMethod(responseWriter, request, httpclient.MethodGet)
}

// GetFromQuery extracts a query parameter from the request URL
func GetFromQuery(request *httpclient.Request, key string) string {
	return request.URL.Query().Get(key)
}

// SendErrorResponse sends an error response with the appropriate status code
func SendErrorResponse(responseWriter httpclient.ResponseWriter, err error, statusCode int) {
	httpclient.Error(responseWriter, err.Error(), statusCode)
}

// SetResponseHeaders sets common headers for HTTP responses
func SetResponseHeaders(responseWriter httpclient.ResponseWriter) {
	responseWriter.Header().Set(ContentTypeKey, ContentTypeValue)
	responseWriter.Header().Set(AcceptKey, ContentTypeValue)
}

// SendJsonResponse sets headers and sends a JSON response with status 200 OK
func SendJsonResponse[T any](responseWriter httpclient.ResponseWriter, data T) error {
	SetResponseHeaders(responseWriter)
	responseWriter.WriteHeader(httpclient.StatusOK)
	return setJson.NewEncoder(responseWriter).Encode(data)
}

// SendJsonResponseWithStatus sets headers and sends a JSON response with the specified status code
func SendJsonResponseWithStatus[T any](responseWriter httpclient.ResponseWriter, data T, statusCode int) error {
	SetResponseHeaders(responseWriter)
	responseWriter.WriteHeader(statusCode)
	return setJson.NewEncoder(responseWriter).Encode(data)
}

// SetHeaders is kept for backward compatibility
// Use SetResponseHeaders instead
func SetHeaders(responseWriter httpclient.ResponseWriter) {
	SetResponseHeaders(responseWriter)
	responseWriter.WriteHeader(httpclient.StatusOK)
}
