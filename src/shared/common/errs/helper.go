package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

// GetErrorType extracts the error type from an error
func GetErrorType(err error) ErrorType {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return ErrOperationFailed // Default error type if not recognized
}

// Map error type to HTTP status code
func GetHTTPStatus(err error) int {
	switch GetErrorType(err) {
	case ErrJSONParse:
		// JSON parse errors result in 400 (Bad Request)
		return http.StatusBadRequest
	case ErrValidation:
		// 422 Unprocessable Entity: Validation error (e.g., invalid data format)
		return http.StatusUnprocessableEntity // 422
	case ErrAuthentication:
		// 401 Unauthorized: Authentication required, wrong credentials
		return http.StatusUnauthorized // 401
	case ErrAuthorization:
		// 403 Forbidden: The user doesn't have permission
		return http.StatusForbidden // 403
	case ErrResourceNotFound:
		// 404 Not Found: The requested resource doesn't exist
		return http.StatusNotFound // 404
	case ErrDuplicateEntry:
		// 409 Conflict: Resource already exists
		return http.StatusConflict // 409
	case ErrDataIntegrity, ErrDatabaseFailure:
		// 500 Internal Server Error: Database issues, constraint violations
		return http.StatusInternalServerError // 500
	case ErrServiceDependency:
		// 503 Service Unavailable: External service failure
		return http.StatusServiceUnavailable // 503
	case ErrBusinessLogic:
		// 400 Bad Request: Business logic issues (e.g., invalid state)
		return http.StatusBadRequest // 400
	case ErrOperationFailed:
		// 500 Internal Server Error: General failure case
		return http.StatusInternalServerError // 500
	default:
		// Default: Unknown errors, fallback to internal server error
		return http.StatusInternalServerError // 500
	}
}

// HandleDBError maps PostgreSQL errors to custom application errors
func HandleDBError(err error) error {
	fmt.Println(errors.Unwrap(err))
	if pgErr, ok := err.(*pq.Error); ok {
		switch pgErr.Code {
		case "23505": // Unique constraint violation
			return NewAppError(ErrDuplicateEntry, "duplicate entry detected: "+pgErr.Message)
		case "23503": // Foreign key violation
			return NewAppError(ErrDataIntegrity, "foreign key constraint violation: "+pgErr.Message)
		case "23502": // Not null violation
			return NewAppError(ErrDataIntegrity, "not null constraint violation: "+pgErr.Message)
		default:
			return NewAppError(ErrDatabaseFailure, "database error: "+pgErr.Message)
		}
	}
	// Fallback for unknown DB errors
	return NewAppError(ErrDatabaseFailure, err.Error())
}

func IsErrValidation(err error) bool {
	if err == nil {
		return false
	}
	return ErrValidation == GetErrorType(err)
}

func IsErrAuthentication(err error) bool {
	if err == nil {
		return false
	}
	return ErrAuthentication == GetErrorType(err)
}

func IsErrAuthorization(err error) bool {
	if err == nil {
		return false
	}
	return ErrAuthorization == GetErrorType(err)
}

func IsErrResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	return ErrResourceNotFound == GetErrorType(err)
}

func IsErrDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	return ErrDuplicateEntry == GetErrorType(err)
}

func IsErrDataIntegrity(err error) bool {
	if err == nil {
		return false
	}
	return ErrDataIntegrity == GetErrorType(err)
}

func IsErrDatabaseFailure(err error) bool {
	if err == nil {
		return false
	}
	return ErrDatabaseFailure == GetErrorType(err)
}
func IsErrServiceDependency(err error) bool {
	if err == nil {
		return false
	}
	return ErrServiceDependency == GetErrorType(err)
}

func IsErrBusinessLogic(err error) bool {
	if err == nil {
		return false
	}
	return ErrBusinessLogic == GetErrorType(err)
}

func IsErrOperationFailed(err error) bool {
	if err == nil {
		return false
	}
	return ErrOperationFailed == GetErrorType(err)
}
