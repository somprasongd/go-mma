package handlers

import (
	"context"
	"go-mma/data/sqldb"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	dbCtx sqldb.DBContext
}

func NewCustomerHandler(dbCtx sqldb.DBContext) *CustomerHandler {
	return &CustomerHandler{dbCtx: dbCtx}
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

	// validate payload
	if len(payload.Name) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "name is required"})
		return
	}

	if payload.CreditLimit <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "credit_limit must be greater than 0"})
		return
	}

	// save new customer to the database
	sql := `INSERT INTO customers (name, credit_limit) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var id int
	if err := h.dbCtx.DB().QueryRowContext(ctx, sql, payload.Name, payload.CreditLimit).Scan(&id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
