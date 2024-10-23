package rest

import (
	"strings"

	"github.com/avast/retry-go"
)

const (
	ErrContextTimeoutSubstr string = "context deadline exceeded"
	ErrIoTimeoutSubstr      string = "i/o timeout"
	ErrNoRouteToHostSubstr  string = "connect: no route to host"
	ErrHTTPEOFSubstr               = ": EOF"

	connPriceRefusedError string = "connect: connection refused"
	connResetError        string = "read: connection reset by peer"
	connWriteResetError   string = "write: connection reset by peer"
	errRespStatus         string = "err resp status:"
	errTCPTimeout         string = "dialing to the given TCP address timed out"
	errTLStimeoutSubstr   string = "net/http: TLS handshake timeout"
	errConnForceClosed    string = "client connection force closed via ClientConn.Close"

	// attemptsRestCount Число попыток.
	attemptsRestCount = 4
)

// GetRestClient Попытки получения данных из REST сервиса.
func GetRestClient(retryableFunc retry.RetryableFunc) error {
	return retry.Do(
		retryableFunc,
		retry.RetryIf(checkTempError),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(attemptsRestCount),
	)
}

// SendRestClient Попытки отправки данных из REST сервиса.
func SendRestClient(retryableFunc retry.RetryableFunc) error {
	return retry.Do(
		retryableFunc,
		retry.RetryIf(checkTempError),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(attemptsRestCount),
	)
}

// checkTempError Определяет временные ошибки.
func checkTempError(err error) bool {
	return strings.Contains(err.Error(), ErrContextTimeoutSubstr) || strings.Contains(err.Error(), ErrIoTimeoutSubstr) ||
		strings.Contains(err.Error(), connPriceRefusedError) || strings.Contains(err.Error(), connResetError) ||
		strings.Contains(err.Error(), connWriteResetError) || strings.Contains(err.Error(), errTLStimeoutSubstr) ||
		strings.Contains(err.Error(), ErrHTTPEOFSubstr) || strings.Contains(err.Error(), ErrNoRouteToHostSubstr) ||
		strings.Contains(err.Error(), errRespStatus) || strings.Contains(err.Error(), errTCPTimeout) ||
		strings.Contains(err.Error(), errConnForceClosed)
}
