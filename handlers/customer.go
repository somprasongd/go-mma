package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	// Implement the logic to create an customer
	type CreateCustomerRequest struct {
		Name        string `json:"name"`
		CreditLimit int    `json:"credit_limit"`
	}
	payload := CreateCustomerRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save new customer to the database
	fmt.Println("Creating new customer:", payload)

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Customer created"})
}
