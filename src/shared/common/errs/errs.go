package errs

import (
	"fmt"
)

// AppError represents a structured error used in the service layer
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	// Err     error     `json:"-"` // Underlying error, not exposed in JSON response
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// // Unwrap allows for errors.Is and errors.As compatibility
// func (e *AppError) Unwrap() error {
// 	return e.Err
// }

// NewAppError creates a new AppError
func NewAppError(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		// Err:     err,
	}
}

// Helper functions for each error type
func NewJSONParseError(message string) *AppError {
	return NewAppError(ErrJSONParse, message)
}

func NewValidationError(message string) *AppError {
	return NewAppError(ErrValidation, message)
}

func NewAuthenticationError(message string) *AppError {
	return NewAppError(ErrAuthentication, message)
}

func NewAuthorizationError(message string) *AppError {
	return NewAppError(ErrAuthorization, message)
}

func NewResourceNotFoundError(message string) *AppError {
	return NewAppError(ErrResourceNotFound, message)
}

func NewDuplicateEntryError(message string) *AppError {
	return NewAppError(ErrDuplicateEntry, message)
}

func NewDataIntegrityError(message string) *AppError {
	return NewAppError(ErrDataIntegrity, message)
}

func NewDatabaseFailureError(message string) *AppError {
	return NewAppError(ErrDatabaseFailure, message)
}

func NewServiceDependencyError(message string) *AppError {
	return NewAppError(ErrServiceDependency, message)
}

func NewBusinessLogicError(message string) *AppError {
	return NewAppError(ErrBusinessLogic, message)
}

func NewOperationFailedError(message string) *AppError {
	return NewAppError(ErrOperationFailed, message)
}
