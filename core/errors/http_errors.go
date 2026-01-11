package errors

import "net/http"

func ToHTTPCode(errCode string) int {
	switch errCode {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeDuplicate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
