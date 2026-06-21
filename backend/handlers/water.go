package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"guppy-breeding/database"
	"guppy-breeding/models"

	"github.com/gin-gonic/gin"
)

func CreateWaterRecord(c *gin.Context) {
	var record models.WaterRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO water_records (tank_id, temperature, pH) VALUES (?, ?, ?)",
		record.TankID, record.Temperature, record.PH,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	record.ID = int(id)

	var recordedAt string
	database.DB.QueryRow("SELECT recorded_at FROM water_records WHERE id = ?", id).Scan(&recordedAt)
	record.RecordedAt = recordedAt

	c.JSON(http.StatusCreated, record)
}

func ListWaterRecords(c *gin.Context) {
	tankID := c.Query("tank_id")
	limit := c.DefaultQuery("limit", "50")

	var rows *sql.Rows
	var err error

	query := "SELECT id, tank_id, temperature, pH, recorded_at FROM water_records"
	if tankID != "" {
		query += " WHERE tank_id = ?"
		query += " ORDER BY recorded_at DESC LIMIT ?"
		rows, err = database.DB.Query(query, tankID, limit)
	} else {
		query += " ORDER BY recorded_at DESC LIMIT ?"
		rows, err = database.DB.Query(query, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var records []models.WaterRecord
	for rows.Next() {
		var r models.WaterRecord
		err := rows.Scan(&r.ID, &r.TankID, &r.Temperature, &r.PH, &r.RecordedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		records = append(records, r)
	}

	c.JSON(http.StatusOK, records)
}

func GetLatestWaterRecord(c *gin.Context) {
	tankIDStr := c.Param("tank_id")
	tankID, err := strconv.Atoi(tankIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tank_id"})
		return
	}

	var record models.WaterRecord
	err = database.DB.QueryRow(
		"SELECT id, tank_id, temperature, pH, recorded_at FROM water_records WHERE tank_id = ? ORDER BY id DESC LIMIT 1",
		tankID,
	).Scan(&record.ID, &record.TankID, &record.Temperature, &record.PH, &record.RecordedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}
