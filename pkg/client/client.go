package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	ts "go-tr-syst/pkg/gen-tr-syst"
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
)

const serverAddr = "localhost:50051"

func createGRPCClient() (ts.TransactionServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("Не удалось подключиться к серверу gRPC: %v", err)
	}

	client := ts.NewTransactionServiceClient(conn)
	return client, conn, nil
}

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/transaction/:id", TransactionIdGet)
	router.POST("/transaction", TransactionPost)

	return router
}

func TransactionPost(c *gin.Context) {
	client, conn, err := createGRPCClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to gRPC server",
		})
		return
	}
	defer conn.Close()

	var requestBody struct {
		Amount               float64 `json:"amount,omitempty"`
		Currency             string  `json:"currency,omitempty"`
		SenderCardNumber     string  `json:"sender_card_number,omitempty"`
		RecipientPhoneNumber string  `json:"recipient_phone_number,omitempty"`
	}

	err = c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	createReq := &ts.CreateTransactionRequest{
		Amount:               fmt.Sprintf("%f",requestBody.Amount),
		Currency:             requestBody.Currency,
		SenderCardNumber:     requestBody.SenderCardNumber,
		RecipientPhoneNumber: requestBody.RecipientPhoneNumber,
	}
	createResp, err := client.CreateTransaction(context.Background(), createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при создании транзакции: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"id":      createResp.Id,
	})
}

func TransactionIdGet(c *gin.Context) {
	client, conn, err := createGRPCClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to gRPC server",
		})
		return
	}
	defer conn.Close()

	id := c.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	getReq := &ts.GetTransactionRequest{
		Id: strconv.FormatInt(idInt64, 10),
	}
	getResp, err := client.GetTransaction(context.Background(), getReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при получении транзакции: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"amount":               getResp.Amount,
		"currency":             getResp.Currency,
		"sender_card_number":   getResp.SenderCardNumber,
		"recipient_phone_number": getResp.RecipientPhoneNumber,
		"date_time":            getResp.DateTime,
	})
}



