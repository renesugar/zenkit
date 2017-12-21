package logging

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/goadesign/goa"
	"github.com/sirupsen/logrus"
	"github.com/goadesign/goa/middleware"
)


// loggingResponseWriter wraps an http.ResponseWriter and writes only raw
// response data (as text) to the context logger.
type loggingResponseWriter struct {
	http.ResponseWriter
	ctx context.Context
}

// Write will write error responses to the logger if debug is enabled
func (lrw *loggingResponseWriter) Write(buf []byte) (int, error) {
	entry := ContextLogger(lrw.ctx)
	if entry != nil {
		logger := entry.Logger
		if logger.Level == logrus.DebugLevel {
			resp := goa.ContextResponse(lrw.ctx)
			if code := resp.ErrorCode; code != "" {
				reqID := middleware.ContextRequestID(lrw.ctx)
				errorResponse := goa.ErrorResponse{}
				err := json.Unmarshal(buf, &errorResponse)
				if err == nil {
					logger.WithFields(logrus.Fields{
						"req_id":  reqID,
						"error_id":  errorResponse.ID,
						"code": errorResponse.Code,
						"status": errorResponse.Status,
						"detail": errorResponse.Detail,
					}).Debug("returned an error")
				} else {
					logger.WithError(err).WithFields(logrus.Fields{
						"req_id":  reqID,
					}).Error("Unable to unmarshall buffer into ErrorResponse")
				}
			}
		}
	}
	return lrw.ResponseWriter.Write(buf)
}

// LogErrorResponse creates an error response logger middleware.
// Only logs the error responses if debug logging is enabled.
// Modeled on middleware.LogResponse
func LogErrorResponse() goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			// chain a new logging writer to the current response writer.
			resp := goa.ContextResponse(ctx)
			resp.SwitchWriter(
				&loggingResponseWriter{
					ResponseWriter: resp.SwitchWriter(nil),
					ctx:            ctx,
				})

			// next
			return h(ctx, rw, req)
		}
	}
}

