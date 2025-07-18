package middleware

import httpclient "net/http"

type DtoResponse struct {
	httpclient.ResponseWriter
	Status int
}

func NewDtoResponse(responseWriter httpclient.ResponseWriter) *DtoResponse {
	return &DtoResponse{
		ResponseWriter: responseWriter,
		Status:         httpclient.StatusOK,
	}
}

func (dtoResponse *DtoResponse) WriteHeader(code int) {
	dtoResponse.Status = code
	dtoResponse.ResponseWriter.WriteHeader(code)
}
