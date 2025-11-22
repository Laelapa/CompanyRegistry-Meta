package logging

import (
	"net/http"

	"github.com/Laelapa/CompanyRegistry/util/netutils"
	"go.uber.org/zap"
)

func (l *Logger) ReqLog(
	loggingFunc func(msg string, fields ...zap.Field),
	r *http.Request,
	msg string,
	additionalFields ...zap.Field,
) {
	fields := l.ReqFields(r)
	fields = append(fields, additionalFields...)
	loggingFunc(msg, fields...)
}

func (l *Logger) ReqInfo(msg string, r *http.Request, additionalFields ...zap.Field) {
	fields := l.ReqFields(r)
	fields = append(fields, additionalFields...)
	l.Info(msg, fields...)
}

func (l *Logger) ReqWarn(msg string, r *http.Request, additionalFields ...zap.Field) {
	fields := l.ReqFields(r)
	fields = append(fields, additionalFields...)
	l.Warn(msg, fields...)
}

func (l *Logger) ReqError(msg string, r *http.Request, err error, adadditionalFields ...zap.Field) {
	fields := l.ReqFields(r)
	fields = append(fields, zap.Error(err))
	fields = append(fields, adadditionalFields...)
	l.Error(msg, fields...)
}

func (l *Logger) ReqFields(r *http.Request) []zap.Field {
	if r == nil {
		return []zap.Field{
			zap.String(FieldError, "request is nil"),
		}
	}

	return []zap.Field{
		zap.String(FieldRemoteAddr, l.FiletLogValue(netutils.GetClientIP(r))),
		zap.String(FieldMethod, r.Method),
		zap.String(FieldPath, l.FiletLogValue(r.URL.Path)),
		zap.String(FieldReferer, l.FiletLogValue(r.Referer())),
	}
}
