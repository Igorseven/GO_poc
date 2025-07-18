package middleware

import (
	constant "PocGo/internal/domain/constants"
	notify "PocGo/internal/domain/notification"
	logger "log"
	httpclient "net/http"
	clock "time"
)

func Logging(next httpclient.Handler) httpclient.Handler {
	return httpclient.HandlerFunc(func(responseWriter httpclient.ResponseWriter, request *httpclient.Request) {
		start := clock.Now()
		wrapped := NewDtoResponse(responseWriter)
		next.ServeHTTP(responseWriter, request)
		logger.Printf(
			notify.LogMiddleware,
			clock.Now().Format(constant.FormatDate),
			request.Method,
			request.RequestURI,
			wrapped.Status,
			clock.Since(start),
		)
	})
}
