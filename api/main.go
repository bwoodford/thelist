package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID int `json:"id"`
	Title string `json:"title"`
        Description string `json:"description" `
	CreatedDate time.Time `json:"createdDate"`
	ModifiedDate time.Time `json:"modifiedDate"`
        CompletedDate time.Time `json:"completedDate"`
	IsActive bool `json:"isActive"`
}

type PatchItem struct {
	Title *string `json:"title"`
        Description *string `json:"description" `
        CompletedDate *time.Time `json:"completedDate"`
	IsActive *bool `json:"isActive"`
}

var now time.Time = time.Now()

var items = []Item {
	{ID: 1, Title: "Bouldering", Description: "Go to Vertical Endeavors and try bouldering.", CreatedDate: now},
	{ID: 2, Title: "LeBurger", Description: "Have a meal at LeBurger.", CreatedDate: now},
	{ID: 3, Title: "Symphony", Description: "See a show at the Minneapolis Symphony.", CreatedDate: now},
}

func createRouter() *gin.Engine {
        r := gin.Default()

        r.GET("/items", getItems)
        r.POST("/items", createItem)
        r.PUT("/items/:id", updateItem)
        r.PATCH("/items/:id", patchItem)
        r.DELETE("/items/:id", deleteItem)

        return r
}

func main() {
        r := createRouter()
        r.Run()
}

func getItems(c *gin.Context) {
        c.JSON(http.StatusOK, items)
}

func createItem(c *gin.Context) {
        var newItem Item
        if err := c.ShouldBindJSON(&newItem); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        newItem.ID = len(items) + 1
        newItem.CreatedDate = time.Now()
        items = append(items, newItem)
        c.JSON(http.StatusCreated, newItem)
}

func updateItem(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
                return
        }

        var updatedItem Item
        if err := c.ShouldBindJSON(&updatedItem); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        for i, item := range items {
                if item.ID == id {
                        items[i].Title = updatedItem.Title
                        items[i].Description = updatedItem.Description
                        items[i].ModifiedDate = time.Now()
                        items[i].IsActive = updatedItem.IsActive
                        c.JSON(http.StatusOK, items[i])
                        return
                }
        }
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

func patchItem(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
                return
        }

        
        var patchData PatchItem
        if err := c.ShouldBindJSON(&patchData); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
                return
        }

        itemIndex := -1

        for i, item := range items {
                if item.ID == id {
                        itemIndex = i
                        break
                }
        }

        if itemIndex == -1 {
                c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
                return
        }

        item := &items[itemIndex]

        if patchData.Title != nil {
                item.Title = *patchData.Title
        }

        if patchData.Description != nil {
                item.Description = *patchData.Description
        }

        item.ModifiedDate = time.Now()

        if patchData.CompletedDate != nil {
                item.CompletedDate = *patchData.CompletedDate
        }

        if patchData.IsActive != nil {
                item.IsActive = *patchData.IsActive
        }

        c.JSON(http.StatusOK, item)
}

func deleteItem(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
                return
        }

        for i, item := range items {
                if item.ID == id {
                        items = append(items[:i], items[i+1:]...)
                        c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
                        return
                }
        }
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}
