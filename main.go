package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Order model
type Order struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	CustomerName string     `json:"customer_name"`
	OrderedAt    time.Time  `json:"ordered_at"`
	Items        []Item     `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Item model
type Item struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
	OrderID     uint   `json:"-"`
}

func initDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=assignment_2 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Order{}, &Item{})
}

func createOrder(c *gin.Context) {
	var request OrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := Order{
		OrderedAt:    request.OrderedAt,
		CustomerName: request.CustomerName,
		Items:        request.Items,
	}

	db.Create(&order)
	c.JSON(http.StatusCreated, convertToOrderResponse(order))
}

func getOrders(c *gin.Context) {
	var orders []Order
	db.Preload("Items").Find(&orders)

	var response []OrderResponse
	for _, order := range orders {
		response = append(response, convertToOrderResponse(order))
	}

	c.JSON(http.StatusOK, response)
}

func updateOrder(c *gin.Context) {
	orderID := c.Param("id")

	var request OrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingOrder Order
	if err := db.Preload("Items").Where("id = ?", orderID).First(&existingOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Update existing order with request data
	existingOrder.OrderedAt = request.OrderedAt
	existingOrder.CustomerName = request.CustomerName
	existingOrder.Items = request.Items

	db.Save(&existingOrder)
	c.JSON(http.StatusOK, convertToOrderResponse(existingOrder))
}

func deleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	var existingOrder Order
	if err := db.Where("id = ?", orderID).First(&existingOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	db.Delete(&existingOrder)
	c.JSON(http.StatusNoContent, gin.H{"message": "Success delete"})
}

func convertToOrderResponse(order Order) OrderResponse {
	return OrderResponse{
		ID:           order.ID,
		OrderedAt:    order.OrderedAt,
		CustomerName: order.CustomerName,
		Items:        order.Items,
	}
}

type OrderRequest struct {
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items"`
}

type OrderResponse struct {
	ID           uint      `json:"id"`
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items,omitempty"`
}

func main() {
	initDB()

	r := gin.Default()

	r.POST("/orders", createOrder)
	r.GET("/orders", getOrders)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)

	port := 8080
	serverAddress := fmt.Sprintf(":%d", port)
	if err := r.Run(serverAddress); err != nil {
		log.Fatal(err)
	}
}
