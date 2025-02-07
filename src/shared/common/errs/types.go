package errs

// ErrorType defines a type for different error categories in the service layer
type ErrorType string

const (
	ErrJSONParse         ErrorType = "json_parse_error"         // JSON Parsing Error
	ErrValidation        ErrorType = "validation_error"         // Invalid input (e.g., missing fields, format issues)
	ErrAuthentication    ErrorType = "authentication_error"     // Wrong credentials, not logged in
	ErrAuthorization     ErrorType = "authorization_error"      // No permission to access resource
	ErrResourceNotFound  ErrorType = "resource_not_found"       // Entity does not exist
	ErrDuplicateEntry    ErrorType = "duplicate_entry"          // Conflict, already exists
	ErrDataIntegrity     ErrorType = "data_integrity_error"     // Foreign key, constraint violations
	ErrDatabaseFailure   ErrorType = "database_failure"         // Generic DB error
	ErrServiceDependency ErrorType = "service_dependency_error" // External service unavailable
	ErrBusinessLogic     ErrorType = "business_logic_error"     // Business rule violation
	ErrOperationFailed   ErrorType = "operation_failed"         // General failure case
)
