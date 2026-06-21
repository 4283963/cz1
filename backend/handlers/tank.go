package handlers

import (
	"net/http"

	"guppy-breeding/database"
	"guppy-breeding/models"

	"github.com/gin-gonic/gin"
)

func CreateTank(c *gin.Context) {
	var tank models.Tank
	if err := c.ShouldBindJSON(&tank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("INSERT INTO tanks (name, location) VALUES (?, ?)", tank.Name, tank.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	tank.ID = int(id)
	c.JSON(http.StatusCreated, tank)
}

func ListTanks(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, COALESCE(location, ''), created_at FROM tanks ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tanks []models.Tank
	for rows.Next() {
		var t models.Tank
		err := rows.Scan(&t.ID, &t.Name, &t.Location, &t.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tanks = append(tanks, t)
	}

	c.JSON(http.StatusOK, tanks)
}
