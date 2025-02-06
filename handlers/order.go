package handlers

import (
	"go-mma/dtos"
	"go-mma/services"
	"go-mma/util/errs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	CancelOrder(c *gin.Context)
}

type orderHandler struct {
	orderServ services.OrderService
}

func NewOrderHandler(orderServ services.OrderService) OrderHandler {
	return &orderHandler{orderServ: orderServ}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	payload := dtos.CreateOrderRequest{}
	if err := c.BindJSON(&payload); err != nil {
		handleError(c, errs.NewJSONParseError(err.Error()))
		return
	}

	id, err := h.orderServ.CreateOrder(c.Request.Context(), &payload)
	if err != nil {
		handleError(c, err)
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *orderHandler) CancelOrder(c *gin.Context) {
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
