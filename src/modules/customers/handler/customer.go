package handler

import (
	"go-mma/modules/customers/service"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/response"
	"net/http"

	"github.com/gin-gonic/gin"

	customerContracts "go-mma/shared/contracts/customer_contracts"
)

type CustomerHandler interface {
	CreateCustomer(c *gin.Context)
}

type customerHandler struct {
	custServ service.CustomerService
}

func NewCustomerHandler(custServ service.CustomerService) CustomerHandler {
	return &customerHandler{custServ: custServ}
}

func (h *customerHandler) CreateCustomer(c *gin.Context) {
	payload := customerContracts.CreateCustomerRequest{}
	if err := c.BindJSON(&payload); err != nil {
		response.HandleError(c, errs.NewJSONParseError(err.Error()))
		return
	}

	id, err := h.custServ.CreateCustomer(c.Request.Context(), &payload)

	if err != nil {
		response.HandleError(c, err)
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
