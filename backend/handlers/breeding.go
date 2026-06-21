package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"guppy-breeding/database"
	"guppy-breeding/models"

	"github.com/gin-gonic/gin"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func validateBreedingRecord(record *models.BreedingRecord) string {
	if record.TankID <= 0 {
		return "请选择鱼缸"
	}
	if record.Strain == "" {
		return "请填写品系名称"
	}
	if record.PairDate == "" {
		return "请填写配对日期"
	}
	if !dateRegex.MatchString(record.PairDate) {
		return "配对日期格式不正确，请使用 YYYY-MM-DD 格式"
	}
	if record.ExpectedBirthDate != "" && !dateRegex.MatchString(record.ExpectedBirthDate) {
		return "预计产仔日期格式不正确，请使用 YYYY-MM-DD 格式"
	}

	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM tanks WHERE id = ?", record.TankID).Scan(&count)
	if err != nil || count == 0 {
		return "所选鱼缸不存在"
	}

	return ""
}

func CreateBreedingRecord(c *gin.Context) {
	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确：" + err.Error()})
		return
	}

	if msg := validateBreedingRecord(&record); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败：" + err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确：" + err.Error()})
		return
	}

	if msg := validateBreedingRecord(&record); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	_, err = database.DB.Exec(
		"UPDATE breeding_records SET tank_id = ?, strain = ?, pair_date = ?, expected_birth_date = ?, notes = ?, status = ? WHERE id = ?",
		record.TankID, record.Strain, record.PairDate, record.ExpectedBirthDate, record.Notes, record.Status, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败：" + err.Error()})
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
