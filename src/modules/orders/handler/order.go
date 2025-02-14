package handler

import (
	"go-mma/shared/common/errs"
	"go-mma/shared/common/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	orderContracts "go-mma/shared/contracts/order_contracts"
)

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	CancelOrder(c *gin.Context)
}

type orderHandler struct {
	orderServ orderContracts.OrderService
}

func NewOrderHandler(orderServ orderContracts.OrderService) OrderHandler {
	return &orderHandler{orderServ: orderServ}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	payload := orderContracts.CreateOrderRequest{}
	if err := c.BindJSON(&payload); err != nil {
		response.HandleError(c, errs.NewJSONParseError(err.Error()))
		return
	}

	id, err := h.orderServ.CreateOrder(c.Request.Context(), &payload)
	if err != nil {
		response.HandleError(c, err)
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
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
