package handler

import (
	"go-mma/modules/customer/dtos"
	"go-mma/modules/customer/service"
	"go-mma/util/errs"
	"go-mma/util/response"
	"net/http"

	"github.com/gin-gonic/gin"
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
	payload := dtos.CreateCustomerRequest{}
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
