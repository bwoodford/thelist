package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)



var handlers Handlers
var TestTime func() time.Time = func() time.Time {
        return time.Date(2025, 0, 0, 0, 0, 0, 0, time.UTC)
}

func init() {
        gin.SetMode(gin.TestMode)
        handlers = Handlers{ 
                currTime: TestTime,
        }
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
                {ID: 1, Title: "Bouldering", Description: "Go to Vertical Endeavors and try bouldering.", IsActive: true, CreatedDate: now},
                {ID: 2, Title: "LeBurger", Description: "Have a meal at LeBurger.", IsActive: true, CreatedDate: now},
                {ID: 3, Title: "Symphony", Description: "See a show at the Minneapolis Symphony.", IsActive: true, CreatedDate: now},
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

        expected := Item{
                ID: itemID,
                Title: "Updated Item",
                Description: "This item has been updated",
                CreatedDate: now,
                ModifiedDate: TestTime(),
                IsActive: false,
        }

        reqBody := map[string]any{
                "title":       expected.Title,
                "description": expected.Description,
                "isActive": expected.IsActive,
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
                t.Errorf("Expected status code: %d, but got: %d", http.StatusOK, w.Code)
        }

        var response Item
        err = json.Unmarshal(w.Body.Bytes(), &response)
        if err != nil {
                t.Fatalf("Failed to unmarshal response body: %v", err)
        }

        if !areItemsEqual([]Item{ expected }, []Item{ response }) {
                t.Errorf("Expected response body: %#v, but got: %#v", expected, response)
        }
}

// TestPatchItem test the PATCH /items/:id endpoint.
func TestPatchItem(t *testing.T) {
        r := createRouter(handlers)

        itemID := 1

        expected := Item{
                ID: itemID,
                Title: "Updated Item",
                Description: "This item has been updated",
                CreatedDate: now,
                ModifiedDate: TestTime(),
                IsActive: false,
        }

        reqBody := map[string]any{
                "title":       expected.Title,
                "description": expected.Description,
                "isActive": expected.IsActive,
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
                t.Errorf("Expected status code: %d, but got: %d", http.StatusOK, w.Code)
        }

        var response Item
        err = json.Unmarshal(w.Body.Bytes(), &response)
        if err != nil {
                t.Fatalf("Failed to unmarshal response body: %v", err)
        }

        if !areItemsEqual([]Item{ expected }, []Item{ response }) {
                t.Errorf("Expected response body: %v, but got: %v", expected, response)
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
        if !areMapsEqual(expectedResponse, response) {
                t.Errorf("Expected response body: %v, but got: %v", expectedResponse, response)
        }
}

// Helper function to compare two slices of Item
func areItemsEqual(a, b []Item) bool {
        if len(a) != len(b) {
                return false
        }
        for i := range a {
                if a[i].ID != b[i].ID || 
                a[i].Title != b[i].Title || 
                a[i].Description != b[i].Description || 
                a[i].IsActive != b[i].IsActive || 
                a[i].CreatedDate != b[i].CreatedDate || 
                a[i].ModifiedDate != b[i].ModifiedDate || 
                a[i].CompletedDate != b[i].CompletedDate {
                        return false
                }
        }
        return true
}

// Helper function to compare two maps
func areMapsEqual(a, b map[string]string) bool {
        if len(a) != len(b) {
                return false
        }
        for key := range a {
                if a[key] != b[key] {
                        return false
                }
        }
        return true
}
