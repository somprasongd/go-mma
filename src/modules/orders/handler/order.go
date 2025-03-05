package handler

import (
	"go-mma/modules/orders/dtos"
	"go-mma/modules/orders/features"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	CancelOrder(c *gin.Context)
}

type orderHandler struct {
}

func NewOrderHandler() OrderHandler {
	return &orderHandler{}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	payload := dtos.CreateOrderRequest{}
	if err := c.BindJSON(&payload); err != nil {
		response.HandleError(c, errs.NewInvalidRequestError(err.Error()))
		return
	}

	result, err := mediator.Send[*features.CreateOrderCommand, *features.CreateOrderResult](
		c.Request.Context(),
		&features.CreateOrderCommand{CreateOrderRequest: &payload},
	)

	if err != nil {
		response.HandleError(c, err)
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, result.CreateOrderResponse)
}

func (h *orderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	_, err = mediator.Send[*features.CancelOrderCommand, *mediator.NoResponse](
		c.Request.Context(),
		&features.CancelOrderCommand{ID: orderID},
	)

	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
