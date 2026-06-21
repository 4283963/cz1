package handlers

import (
	"net/http"
	"strconv"

	"guppy-breeding/database"
	"guppy-breeding/models"

	"github.com/gin-gonic/gin"
)

func CreateBreedingRecord(c *gin.Context) {
	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if record.Status == "" {
		record.Status = "breeding"
	}

	result, err := database.DB.Exec(
		"INSERT INTO breeding_records (tank_id, strain, pair_date, expected_birth_date, notes, status) VALUES (?, ?, ?, ?, ?, ?)",
		record.TankID, record.Strain, record.PairDate, record.ExpectedBirthDate, record.Notes, record.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	record.ID = int(id)

	var createdAt string
	database.DB.QueryRow("SELECT created_at FROM breeding_records WHERE id = ?", id).Scan(&createdAt)
	record.CreatedAt = createdAt

	c.JSON(http.StatusCreated, record)
}

func ListBreedingRecords(c *gin.Context) {
	tankID := c.Query("tank_id")
	status := c.Query("status")

	query := "SELECT id, tank_id, strain, pair_date, COALESCE(expected_birth_date, ''), COALESCE(notes, ''), status, created_at FROM breeding_records WHERE 1=1"
	var args []interface{}

	if tankID != "" {
		query += " AND tank_id = ?"
		args = append(args, tankID)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY pair_date DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var records []models.BreedingRecord
	for rows.Next() {
		var r models.BreedingRecord
		err := rows.Scan(&r.ID, &r.TankID, &r.Strain, &r.PairDate, &r.ExpectedBirthDate, &r.Notes, &r.Status, &r.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		records = append(records, r)
	}

	c.JSON(http.StatusOK, records)
}

func GetBreedingRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var record models.BreedingRecord
	err = database.DB.QueryRow(
		"SELECT id, tank_id, strain, pair_date, COALESCE(expected_birth_date, ''), COALESCE(notes, ''), status, created_at FROM breeding_records WHERE id = ?",
		id,
	).Scan(&record.ID, &record.TankID, &record.Strain, &record.PairDate, &record.ExpectedBirthDate, &record.Notes, &record.Status, &record.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func UpdateBreedingRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = database.DB.Exec(
		"UPDATE breeding_records SET tank_id = ?, strain = ?, pair_date = ?, expected_birth_date = ?, notes = ?, status = ? WHERE id = ?",
		record.TankID, record.Strain, record.PairDate, record.ExpectedBirthDate, record.Notes, record.Status, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	record.ID = id
	c.JSON(http.StatusOK, record)
}

func DeleteBreedingRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM breeding_records WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
