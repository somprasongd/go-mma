package handlers

import (
	"go-mma/data/sqldb"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	dbCtx sqldb.DBContext
}

func NewOrderHandler(dbCtx sqldb.DBContext) *OrderHandler {
	return &OrderHandler{dbCtx: dbCtx}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Implement the logic to create an order
	var CreateOrderRequest struct {
		CustomerID string `json:"customer_id"`
		OrderTotal int    `json:"order_total"`
	}
	if err := c.ShouldBindJSON(&CreateOrderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get customer details from the database
	// customerID := CreateOrderRequest.CustomerID

	// check credit limit
	// if creditLimit < CreateOrderRequest.OrderTotal {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient credit limit"})
	// 	return
	// }

	// reserve credit limit for the customer
	// creditLimit -= CreateOrderRequest.OrderTotal
	// save the customer details to the database

	// save the order to the database
	log.Println("Creating new order:", CreateOrderRequest)

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	// Implement the logic to cancel an order
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// check if the order exists in the database
	// if !orderExists(orderID) {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	// 	return
	// }

	// get order details from the database
	// order := getOrder(orderID)
	// if order == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	// 	return
	// }

	// get cutomer details from the database
	// customer := getCustomer(order.CustomerID)
	// if customer == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
	// 	return
	// }

	// release credit limit for the customer
	// creditLimit += CreateOrderRequest.OrderTotal
	// save the customer details to the database

	// update the order status in the database
	log.Println("Cancelling order:", orderID)

	//
	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}
