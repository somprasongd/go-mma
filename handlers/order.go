package handlers

import (
	"go-mma/dtos"
	"go-mma/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderServ *services.OrderService
}

func NewOrderHandler(orderServ *services.OrderService) *OrderHandler {
	return &OrderHandler{orderServ: orderServ}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	payload := dtos.CreateOrderRequest{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.orderServ.CreateOrder(c.Request.Context(), &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.orderServ.CancelOrder(c.Request.Context(), orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
