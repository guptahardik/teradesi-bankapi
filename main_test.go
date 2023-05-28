package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/accounts/:id", GetAccount)
	router.POST("/accounts/:id/deposit", Deposit)
	router.POST("/accounts/:id/withdraw", Withdraw)
	router.POST("/accounts/:id/create", CreateAccount)
	router.POST("/accounts/:id/add-money", AddMoney)
	return router
}

func TestGetAccount(t *testing.T) {
	accounts = make(map[string]Account)
	accounts["1"] = Account{
		ID:      "1",
		Balance: 500,
	}

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/accounts/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	var account Account
	_ = json.Unmarshal(rr.Body.Bytes(), &account)

	if account.ID != "1" {
		t.Errorf("Expected account ID '1', got '%v'", account.ID)
	}
	if account.Balance != 500 {
		t.Errorf("Expected account balance 500, got %v", account.Balance)
	}
}

func TestDeposit(t *testing.T) {
	accounts = make(map[string]Account)
	accounts["1"] = Account{
		ID:      "1",
		Balance: 500,
	}

	router := setupRouter()

	amount := 100.0
	jsonData, _ := json.Marshal(amount)
	req, _ := http.NewRequest("POST", "/accounts/1/deposit", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	var account Account
	_ = json.Unmarshal(rr.Body.Bytes(), &account)

	if account.Balance != 600 {
		t.Errorf("Expected account balance 600, got %v", account.Balance)
	}
}

func TestWithdraw(t *testing.T) {
	accounts = make(map[string]Account)
	accounts["1"] = Account{
		ID:      "1",
		Balance: 500,
	}

	router := setupRouter()

	amount := 200.0
	jsonData, _ := json.Marshal(amount)
	req, _ := http.NewRequest("POST", "/accounts/1/withdraw", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	var account Account
	_ = json.Unmarshal(rr.Body.Bytes(), &account)

	if account.Balance != 300 {
		t.Errorf("Expected account balance 300, got %v", account.Balance)
	}
}

func TestWithdrawInvalidAmount(t *testing.T) {
	accounts = make(map[string]Account)
	accounts["1"] = Account{
		ID:      "1",
		Balance: 500,
	}

	router := setupRouter()

	amount := 600.0
	jsonData, _ := json.Marshal(amount)
	req, _ := http.NewRequest("POST", "/accounts/1/withdraw", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rr.Code)
	}
}

func TestGetAccountNotFound(t *testing.T) {
	accounts = make(map[string]Account)

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/accounts/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %v", rr.Code)
	}
}

func TestDepositInvalidAmount(t *testing.T) {
	accounts = make(map[string]Account)
	accounts["1"] = Account{
		ID:      "1",
		Balance: 500,
	}

	router := setupRouter()

	amount := 15000.0
	jsonData, _ := json.Marshal(amount)
	req, _ := http.NewRequest("POST", "/accounts/1/deposit", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rr.Code)
	}
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	accounts = make(map[string]Account)
	accounts["1"] = Account{
		ID:      "1",
		Balance: 500,
	}

	router := setupRouter()

	amount := 600.0
	jsonData, _ := json.Marshal(amount)
	req, _ := http.NewRequest("POST", "/accounts/1/withdraw", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rr.Code)
	}
}
func TestAddMoney(t *testing.T) {
	router := setupRouter()

	// Create a new account
	reqCreate, _ := http.NewRequest("POST", "/accounts/12345/create", nil)
	rrCreate := httptest.NewRecorder()
	router.ServeHTTP(rrCreate, reqCreate)

	if rrCreate.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %v", rrCreate.Code)
		return
	}

	// Add money to the account
	addMoneyData := struct {
		Amount float64 `json:"amount"`
	}{
		Amount: 500.0,
	}

	addMoneyJSON, _ := json.Marshal(addMoneyData)

	reqAddMoney, _ := http.NewRequest("POST", "/accounts/12345/add-money", bytes.NewReader(addMoneyJSON))
	rrAddMoney := httptest.NewRecorder()
	router.ServeHTTP(rrAddMoney, reqAddMoney)

	if rrAddMoney.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rrAddMoney.Code)
		return
	}

	// Retrieve the account details
	reqGetAccount, _ := http.NewRequest("GET", "/accounts/12345", nil)
	rrGetAccount := httptest.NewRecorder()
	router.ServeHTTP(rrGetAccount, reqGetAccount)

	if rrGetAccount.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rrGetAccount.Code)
		return
	}

	var account Account
	if err := json.NewDecoder(rrGetAccount.Body).Decode(&account); err != nil {
		t.Error("Failed to decode response body")
		return
	}

	expectedBalance := 500.0
	if account.Balance != expectedBalance {
		t.Errorf("Expected account balance %v, got %v", expectedBalance, account.Balance)
	}
}

func TestCreateAccount(t *testing.T) {
	accounts = make(map[string]Account)

	router := setupRouter()

	req, _ := http.NewRequest("POST", "/accounts/12345/create", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %v", rr.Code)
	}

	var account Account
	_ = json.Unmarshal(rr.Body.Bytes(), &account)

	if account.ID != "12345" {
		t.Errorf("Expected account ID '12345', got '%v'", account.ID)
	}
	if account.Balance != 0 {
		t.Errorf("Expected account balance 0, got %v", account.Balance)
	}
}
