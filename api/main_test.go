package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var handlers Handlers
var createTime time.Time = time.Now().UTC()

func createTestData(db *sql.DB) {
	items := []Item{
		{ID: 1, Title: "Bouldering", Description: "Go to Vertical Endeavors and try bouldering.", IsActive: true, CreatedDate: createTime},
		{ID: 2, Title: "LeBurger", Description: "Have a meal at LeBurger.", IsActive: true, CreatedDate: createTime},
		{ID: 3, Title: "Symphony", Description: "See a show at the Minneapolis Symphony.", IsActive: true, CreatedDate: createTime},
	}

	for _, item := range items {
		stmt := `
		INSERT INTO items (
			title,
			description,
			created_date,
			modified_date,
			completed_date,
			is_active
		) VALUES (?, ?, ?, ?, ?, ?)`

		_, err := db.Exec(stmt,
			item.Title,
			item.Description,
			item.CreatedDate,
			item.ModifiedDate,
			item.CompletedDate,
			item.IsActive)

		if err != nil {
			panic(err.Error())
		}
	}
}

func init() {
	gin.SetMode(gin.TestMode)
	db, err := InitDB(":memory:")
	if err != nil {
		panic(err.Error())
	}
	handlers = Handlers{
		db: db,
	}
	createTestData(db)
}

// TestGetItems tests the GET /items endpoint.
func TestGetItems(t *testing.T) {
	r := createRouter(handlers)

	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code: %d, but got: %d", http.StatusOK, w.Code)
	}

	expectedResponse := []Item{
		{ID: 1, Title: "Bouldering", Description: "Go to Vertical Endeavors and try bouldering.", IsActive: true, CreatedDate: createTime},
		{ID: 2, Title: "LeBurger", Description: "Have a meal at LeBurger.", IsActive: true, CreatedDate: createTime},
		{ID: 3, Title: "Symphony", Description: "See a show at the Minneapolis Symphony.", IsActive: true, CreatedDate: createTime},
	}

	var response []Item
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !areItemsEqual(expectedResponse, response) {
		t.Errorf("Expected response body: %v, but got: %v", expectedResponse, response)
	}
}

// TestCreateItem tests the POST /items endpoint.
func TestCreateItem(t *testing.T) {
	r := createRouter(handlers)

	newItem := Item{Title: "New Item", Description: "This is a new item."}
	jsonBody, err := json.Marshal(newItem)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code: %d, but got: %d", http.StatusCreated, w.Code)
	}

	var response Item
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if response.Title != newItem.Title || response.Description != newItem.Description {
		t.Errorf("Expected response body: %v, but got: %v", newItem, response)
	}
}

// TestUpdateItem tests the PUT /items/:id endpoint.
func TestUpdateItem(t *testing.T) {
	r := createRouter(handlers)

	itemID := 1

	modifyTime := time.Now().UTC()

	CurrentTime = func() time.Time {
		return modifyTime
	}

	expected := Item{
		ID:           itemID,
		Title:        "Updated Item",
		Description:  "This item has been updated",
		CreatedDate:  createTime,
		ModifiedDate: &modifyTime,
		IsActive:     false,
	}

	reqBody := map[string]any{
		"title":       expected.Title,
		"description": expected.Description,
		"isActive":    expected.IsActive,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	url := "/items/" + strconv.Itoa(itemID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Unexpected result: %v\n", w.Result())
		t.Fatalf("Expected status code: %d, but got: %d", http.StatusOK, w.Code)
	}

	var response Item
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !areItemsEqual([]Item{expected}, []Item{response}) {
		t.Errorf("Expected response body: %#v, but got: %#v", expected, response)
	}
}

// TestPatchItem test the PATCH /items/:id endpoint.
func TestPatchItem(t *testing.T) {
	r := createRouter(handlers)

	itemID := 1

	modifyTime := time.Now().UTC()

	CurrentTime = func() time.Time {
		return modifyTime
	}

	expected := Item{
		ID:           itemID,
		Title:        "Updated Item",
		Description:  "This item has been updated",
		CreatedDate:  createTime,
		ModifiedDate: &modifyTime,
		IsActive:     false,
	}

	reqBody := map[string]any{
		"title":       expected.Title,
		"description": expected.Description,
		"isActive":    expected.IsActive,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	url := "/items/" + strconv.Itoa(itemID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code: %d, but got: %d", http.StatusOK, w.Code)
	}

	var response Item
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !areItemsEqual([]Item{expected}, []Item{response}) {
		t.Errorf("Expected response body: %#v, but got: %#v", expected, response)
	}
}

// TestDeleteItem tests the DELETE /items/:id endpoint.
func TestDeleteItem(t *testing.T) {
	r := createRouter(handlers)

	itemID := 2

	url := "/items/" + strconv.Itoa(itemID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code: %d, but got: %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	expectedResponse := map[string]string{
		"message": "Item deleted",
	}
	if !maps.Equal(expectedResponse, response) {
		t.Errorf("Expected response body: %v, but got: %v", expectedResponse, response)
	}
}

func areItemsEqual(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i, item := range a {
		if !item.Equals(&b[i]) {
			return false
		}
	}
	return true
}
