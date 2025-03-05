package errs

// ErrorType defines a type for different error categories in the service layer
type ErrorType string

const (
	// Input and Request Issues
	ErrInvalidRequest ErrorType = "invalid_request"  // Malformed JSON, missing/invalid fields
	ErrValidation     ErrorType = "validation_error" // Field-level validation (business rules)

	// Authentication and Authorization
	ErrAuthentication ErrorType = "authentication_error" // Invalid credentials
	ErrAuthorization  ErrorType = "authorization_error"  // Permission denied

	// Resource State Issues
	ErrResourceNotFound ErrorType = "resource_not_found"   // Missing entity
	ErrDuplicateEntry   ErrorType = "duplicate_entry"      // Conflict (e.g., unique constraint)
	ErrBusinessLogic    ErrorType = "business_logic_error" // Violations of business rules

	// Infrastructure and General Failures
	ErrDataIntegrity     ErrorType = "data_integrity_error" // Foreign key, constraint violations
	ErrDatabaseFailure   ErrorType = "database_failure"     // DB-level errors
	ErrServiceDependency ErrorType = "service_dependency"   // External service unavailability
	ErrOperationFailed   ErrorType = "operation_failed"     // Generic internal failures
)
