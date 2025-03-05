package errs

import (
	"errors"
	"net/http"
	"regexp"

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
	case ErrInvalidRequest:
		// 400 Bad Request: JSON parsing errors, invalid input format
		return http.StatusBadRequest // 400
	case ErrValidation:
		// 422 Unprocessable Entity: Valid request but fails business rules
		return http.StatusUnprocessableEntity // 422
	case ErrAuthentication:
		// 401 Unauthorized: Authentication failure
		return http.StatusUnauthorized // 401
	case ErrAuthorization:
		// 403 Forbidden: No permission to access resource
		return http.StatusForbidden // 403
	case ErrResourceNotFound:
		// 404 Not Found: Missing or deleted entity
		return http.StatusNotFound // 404
	case ErrDuplicateEntry:
		// 409 Conflict: Resource already exists (e.g., unique violation)
		return http.StatusConflict // 409
	case ErrBusinessLogic:
		// 422 Unprocessable Entity: Failed logical validation (e.g., insufficient funds)
		return http.StatusUnprocessableEntity // 422
	case ErrDataIntegrity, ErrDatabaseFailure:
		// 500 Internal Server Error: Database-related issues
		return http.StatusInternalServerError // 500
	case ErrServiceDependency:
		// 503 Service Unavailable: External service failure
		return http.StatusServiceUnavailable // 503
	case ErrOperationFailed:
		// 500 Internal Server Error: Catch-all for other failures
		return http.StatusInternalServerError // 500
	default:
		// 500 Internal Server Error: Unknown errors
		return http.StatusInternalServerError // 500
	}
}

// HandleDBError maps PostgreSQL errors to custom application errors
func HandleDBError(err error) error {
	if pgErr, ok := err.(*pq.Error); ok {
		switch pgErr.Code {
		case "23505": // Unique violation
			// Attempt to parse the field name and value from the error message
			re := regexp.MustCompile(`key \((.*?)\)=\((.*?)\) already exists`)
			matches := re.FindStringSubmatch(pgErr.Message)
			if len(matches) > 0 {
				return NewAppError(ErrDuplicateEntry, "duplicate entry violation", err, map[string]interface{}{
					"field": matches[1],
					"value": matches[2],
				})
			}
			// Fallback if parsing fails
			return NewAppError(ErrDuplicateEntry, "duplicate entry violation", err, nil)
		case "23503": // Foreign key violation
			// Foreign key violation message format: Key (column_name)=(value) is still referenced from table "other_table"
			re := regexp.MustCompile(`Key \((\S+)\)=\((.*?)\)`)
			matches := re.FindStringSubmatch(pgErr.Message)
			if len(matches) > 0 {
				// Extract column and value from the error message
				return NewAppError(ErrDataIntegrity, "foreign key violation", err, map[string]interface{}{
					"field": matches[1], // Extracted column name
					"value": matches[2], // Extracted value
				})
			}
			// If parsing fails, fallback to a generic error
			return NewAppError(ErrDataIntegrity, "foreign key violation", err, nil)

		case "23502": // Not-null violation
			// Not-null violation message format: null value in column "column_name" violates not-null constraint
			re := regexp.MustCompile(`null value in column "(\S+)" violates not-null constraint`)
			matches := re.FindStringSubmatch(pgErr.Message)
			if len(matches) > 0 {
				// Extract column name from the error message
				return NewAppError(ErrDataIntegrity, "not-null constraint violation", err, map[string]interface{}{
					"field": matches[1], // Extracted column name
				})
			}
			// If parsing fails, fallback to a generic error
			return NewAppError(ErrDataIntegrity, "not-null constraint violation", err, nil)

		default:
			// Handle other cases like duplicate entry or general DB failure
			return NewAppError(ErrDatabaseFailure, "database error", err, nil)
		}
	}
	// Fallback for unknown DB errors
	return NewAppError(ErrDatabaseFailure, "database error", err, nil)
}

// func IsAppError(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	_, ok := err.(*AppError)
// 	return ok
// }

// func IsErrValidation(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrValidation == GetErrorType(err)
// }

// func IsErrAuthentication(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrAuthentication == GetErrorType(err)
// }

// func IsErrAuthorization(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrAuthorization == GetErrorType(err)
// }

// func IsErrResourceNotFound(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrResourceNotFound == GetErrorType(err)
// }

func IsErrDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	return ErrDuplicateEntry == GetErrorType(err)
}

// func IsErrDataIntegrity(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrDataIntegrity == GetErrorType(err)
// }

// func IsErrDatabaseFailure(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrDatabaseFailure == GetErrorType(err)
// }
// func IsErrServiceDependency(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrServiceDependency == GetErrorType(err)
// }

// func IsErrBusinessLogic(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrBusinessLogic == GetErrorType(err)
// }

// func IsErrOperationFailed(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return ErrOperationFailed == GetErrorType(err)
// }
