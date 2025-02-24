package handler

import (
	"go-mma/modules/customers/dtos"
	"go-mma/modules/customers/features"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler interface {
	CreateCustomer(c *gin.Context)
}

type customerHandler struct {
}

func NewCustomerHandler() CustomerHandler {
	return &customerHandler{}
}

func (h *customerHandler) CreateCustomer(c *gin.Context) {
	payload := dtos.CreateCustomerRequest{}
	if err := c.BindJSON(&payload); err != nil {
		response.HandleError(c, errs.NewJSONParseError(err.Error()))
		return
	}

	resp, err := mediator.Send[*features.CreateCustomerCommand, *dtos.CreateCustomerResponse](
		c.Request.Context(),
		&features.CreateCustomerCommand{CreateCustomerRequest: &payload},
	)

	if err != nil {
		response.HandleError(c, err)
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, resp)
}
