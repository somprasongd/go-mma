package response

import (
	"go-mma/shared/common/errs"

	"github.com/gin-gonic/gin"
)

// Centralized error handling
func HandleError(c *gin.Context, err error) {
	// Convert non-AppError to AppError with type ErrOperationFailed
	appErr, ok := err.(*errs.AppError)
	if !ok {
		appErr = &errs.AppError{
			Type:    errs.ErrOperationFailed,
			Message: err.Error(),
		}
		err = appErr
	}

	// Get the appropriate HTTP status code
	statusCode := errs.GetHTTPStatus(err)

	// Return structured response with error type and message
	c.JSON(statusCode, gin.H{
		"type":    appErr.Type,
		"message": appErr.Message,
	})
}
