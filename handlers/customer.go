package handlers

import (
	"go-mma/dtos"
	"go-mma/services"
	"go-mma/util/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler interface {
	CreateCustomer(c *gin.Context)
}

type customerHandler struct {
	custServ services.CustomerService
}

func NewCustomerHandler(custServ services.CustomerService) CustomerHandler {
	return &customerHandler{custServ: custServ}
}

func (h *customerHandler) CreateCustomer(c *gin.Context) {
	payload := dtos.CreateCustomerRequest{}
	if err := c.BindJSON(&payload); err != nil {
		handleError(c, errs.NewJSONParseError(err.Error()))
		return
	}

	id, err := h.custServ.CreateCustomer(c.Request.Context(), &payload)

	if err != nil {
		handleError(c, err)
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
