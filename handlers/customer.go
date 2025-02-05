package handlers

import (
	"go-mma/dtos"
	"go-mma/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.custServ.CreateCustomer(c.Request.Context(), &payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
