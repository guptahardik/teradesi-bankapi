package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Account struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

var accounts map[string]Account

func main() {
	accounts = make(map[string]Account)

	router := gin.Default()

	router.GET("/accounts/:id", GetAccount)
	router.POST("/accounts/:id/deposit", Deposit)
	router.POST("/accounts/:id/withdraw", Withdraw)
	router.POST("/accounts/:id/create", CreateAccount)
	router.POST("/accounts/:id/add-money", AddMoney)

	if err := router.Run(":7000"); err != nil {
		log.Fatal("Server error:", err)
	}
}

func GetAccount(c *gin.Context) {
	accountID := c.Param("id")
	account, ok := accounts[accountID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}

func Deposit(c *gin.Context) {
	accountID := c.Param("id")

	var depositAmount float64
	if err := c.ShouldBindJSON(&depositAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit amount"})
		return
	}

	account, ok := accounts[accountID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if depositAmount > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deposit amount exceeds the limit"})
		return
	}

	account.Balance += depositAmount
	accounts[accountID] = account

	c.JSON(http.StatusOK, account)
}

func Withdraw(c *gin.Context) {
	accountID := c.Param("id")

	var withdrawAmount float64
	if err := c.ShouldBindJSON(&withdrawAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdraw amount"})
		return
	}

	account, ok := accounts[accountID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if withdrawAmount > 0.9*account.Balance || account.Balance-withdrawAmount < 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	account.Balance -= withdrawAmount
	accounts[accountID] = account

	c.JSON(http.StatusOK, account)
}

func CreateAccount(c *gin.Context) {
	accountID := c.Param("id")

	_, ok := accounts[accountID]
	if ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Account already exists"})
		return
	}

	account := Account{
		ID:      accountID,
		Balance: 0,
	}

	accounts[accountID] = account

	c.JSON(http.StatusCreated, account)
}

func AddMoney(c *gin.Context) {
	accountID := c.Param("id")

	var depositRequest struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&depositRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit amount"})
		return
	}

	account, ok := accounts[accountID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if depositRequest.Amount > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deposit amount exceeds the limit"})
		return
	}

	account.Balance += depositRequest.Amount
	accounts[accountID] = account

	c.JSON(http.StatusOK, account)
}
