package response

import (
	"go-mma/shared/common/errs"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	ErrorType    errs.ErrorType `json:"error_type"`
	ErrorMessage string         `json:"error_message"`
	Context      map[string]any `json:"context,omitempty"`
}

// Centralized error handling
func HandleError(c *gin.Context, err error) {
	// Convert non-AppError to AppError with type ErrOperationFailed
	appErr, ok := err.(*errs.AppError)
	if !ok {
		appErr = errs.NewAppError(errs.ErrOperationFailed, "internal error occurred", err, nil)
		err = appErr
	}

	// Get the appropriate HTTP status code
	statusCode := errs.GetHTTPStatus(err)

	// Return structured response with error type and message
	c.JSON(statusCode, ErrorResponse{
		ErrorType:    appErr.Type,
		ErrorMessage: appErr.Message,
		Context:      appErr.Context,
	})
}
