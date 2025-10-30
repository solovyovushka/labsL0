package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTP_GetOrderSuccess(t *testing.T) {
	cache := NewCache()
	testUID := "test_http_success"
	testOrder := Order{
		OrderUID: testUID,
		TrackNumber: "TEST_TRACK",
	}
	cache.Set(testUID, testOrder)
	
	router := SetupRoutes(cache)
	req, _ := http.NewRequest("GET", "/order/"+testUID, nil)
	rr := httptest.NewRecorder()
	
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")
	
	var responseOrder Order
	err := json.Unmarshal(rr.Body.Bytes(), &responseOrder)
	assert.NoError(t, err, "Should successfully unmarshal response")
	assert.Equal(t, testUID, responseOrder.OrderUID, "Response UID should match")
}

func TestHTTP_GetOrderNotFound(t *testing.T) {
	cache := NewCache()
	
	router := SetupRoutes(cache)
	req, _ := http.NewRequest("GET", "/order/non_existent", nil)
	rr := httptest.NewRecorder()
	
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusNotFound, rr.Code, "Expected status code 404")
	assert.Equal(t, "Order not found\n", rr.Body.String(), "Expected specific error message")
}