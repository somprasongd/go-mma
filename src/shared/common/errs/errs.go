package errs

import (
	"fmt"
	"runtime"
	"strings"
)

// AppError represents a structured error used in the service layer
type AppError struct {
	Type       ErrorType      `json:"type"`        // Error category
	Message    string         `json:"message"`     // Friendly message for clients
	Cause      error          `json:"-"`           // Root cause (internal)
	Context    map[string]any `json:"context"`     // Metadata for debugging
	StackTrace string         `json:"stack_trace"` // Call stack (for logs)
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// // Unwrap allows for errors.Is and errors.As compatibility
// func (e *AppError) Unwrap() error {
// 	return e.Cause
// }

// NewAppError creates a new AppError
func NewAppError(typ ErrorType, message string, cause error, ctx map[string]any) *AppError {
	return &AppError{
		Type:       typ,
		Message:    message,
		Cause:      cause,
		Context:    ctx,
		StackTrace: captureStackTrace(3),
	}
}

// Capture call stack (useful for multi-layer tracing)
func captureStackTrace(skip int) string {
	var sb strings.Builder
	for i := skip; i < skip+10; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		sb.WriteString(fmt.Sprintf("%s:%d\n", file, line))
	}
	return sb.String()
}

// Helper functions for each error type
func NewInvalidRequestError(message string) *AppError {
	return NewAppError(ErrInvalidRequest, message, nil, nil)
}

func NewValidationError(message string) *AppError {
	return NewAppError(ErrValidation, message, nil, nil)
}

func NewAuthenticationError(message string) *AppError {
	return NewAppError(ErrAuthentication, message, nil, nil)
}

func NewAuthorizationError(message string) *AppError {
	return NewAppError(ErrAuthorization, message, nil, nil)
}

func NewResourceNotFoundError(message string) *AppError {
	return NewAppError(ErrResourceNotFound, message, nil, nil)
}

func NewDuplicateEntryError(message string) *AppError {
	return NewAppError(ErrDuplicateEntry, message, nil, nil)
}

func NewDataIntegrityError(message string) *AppError {
	return NewAppError(ErrDataIntegrity, message, nil, nil)
}

func NewDatabaseFailureError(message string) *AppError {
	return NewAppError(ErrDatabaseFailure, message, nil, nil)
}

func NewServiceDependencyError(message string) *AppError {
	return NewAppError(ErrServiceDependency, message, nil, nil)
}

func NewBusinessLogicError(message string) *AppError {
	return NewAppError(ErrBusinessLogic, message, nil, nil)
}

func NewOperationFailedError(message string) *AppError {
	return NewAppError(ErrOperationFailed, message, nil, nil)
}
