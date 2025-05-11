package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description" `
	CreatedDate   time.Time  `json:"createdDate"`
	ModifiedDate  *time.Time `json:"modifiedDate"`
	CompletedDate *time.Time `json:"completedDate"`
	IsActive      bool       `json:"isActive"`
}

func (i *Item) Equals(o *Item) bool {

	timePointersEqual := func(t1 *time.Time, t2 *time.Time) bool {
		if t1 == nil && t2 == nil {
			return true
		}
		if t1 != nil && t2 == nil {
			return false
		}
		if t1 == nil && t2 != nil {
			return false
		}
		if !(*t1).Equal(*t2) {
			return false
		}
		return true
	}

	if i.ID != o.ID {
		return false
	}

	if i.Title != o.Title {
		return false
	}

	if i.Description != o.Description {
		return false
	}

	if i.IsActive != o.IsActive {
		return false
	}

	if !i.CreatedDate.Equal(o.CreatedDate) {
		return false
	}

	if !timePointersEqual(i.ModifiedDate, o.ModifiedDate) {
		return false
	}

	if !timePointersEqual(i.CompletedDate, o.CompletedDate) {
		return false
	}

	return true
}

type PutItem struct {
	Title         string     `json:"title"`
	Description   string     `json:"description" `
	CompletedDate *time.Time `json:"completedDate"`
	IsActive      bool       `json:"isActive"`
}

type PatchItem struct {
	Title         *string    `json:"title"`
	Description   *string    `json:"description" `
	CompletedDate *time.Time `json:"completedDate"`
	IsActive      *bool      `json:"isActive"`
}

func (m *PatchItem) HasChanges() bool {
	if m.Title != nil {
		return true
	}

	if m.Description != nil {
		return true
	}

	if m.CompletedDate != nil {
		return true
	}

	if m.IsActive != nil {
		return true
	}

	return false
}

type Handlers struct {
	db *sql.DB
}

func createRouter(handlers Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/items", handlers.getItems)
	r.POST("/items", handlers.createItem)
	r.PUT("/items/:id", handlers.putItem)
	r.PATCH("/items/:id", handlers.patchItem)
	r.DELETE("/items/:id", handlers.deleteItem)

	return r
}

func main() {

	dbPath, isSet := os.LookupEnv("DATABASE_PATH")

	if isSet == false {
		panic("DATABASE_PATH environment variable is not set.")
	}

	db, err := InitDB(dbPath)
	if err != nil {
		panic(err)
	}

	h := Handlers{
		db: db,
	}

	r := createRouter(h)
	r.Run()
}

func (h *Handlers) getItems(c *gin.Context) {
	items, err := GetItems(h.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handlers) createItem(c *gin.Context) {
	var newItem Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := InsertItem(h.db, &newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	item, err := GetItem(h.db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, item)
}

func (h *Handlers) putItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updatedItem PutItem
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := GetItem(h.db, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No item was found"})
		return
	}

	item.CompletedDate = updatedItem.CompletedDate
	item.Title = updatedItem.Title
	item.Description = updatedItem.Description
	item.IsActive = updatedItem.IsActive

	now := CurrentTime()

	item.ModifiedDate = &now

	if err = UpdateItem(h.db, id, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handlers) patchItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var patchItem PatchItem
	if err := c.ShouldBindJSON(&patchItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if !patchItem.HasChanges() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes detected"})
		return
	}

	item, err := GetItem(h.db, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No item was found"})
		return
	}

	if patchItem.Title != nil {
		item.Title = *patchItem.Title
	}

	if patchItem.Description != nil {
		item.Description = *patchItem.Description
	}

	if patchItem.IsActive != nil {
		item.IsActive = *patchItem.IsActive
	}

	if patchItem.CompletedDate != nil {
		item.CompletedDate = patchItem.CompletedDate
	}

	now := CurrentTime()

	item.ModifiedDate = &now

	if err = UpdateItem(h.db, id, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handlers) deleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	if err = DeleteItem(h.db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
