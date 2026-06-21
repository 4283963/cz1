package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"guppy-breeding/database"
	"guppy-breeding/models"

	"github.com/gin-gonic/gin"
)

func GetAlertConfig(c *gin.Context) {
	var config models.AlertConfig
	err := database.DB.QueryRow(
		"SELECT id, temp_min, temp_max, consecutive_count, notify_enabled, COALESCE(webhook_url, ''), updated_at FROM alert_config ORDER BY id DESC LIMIT 1",
	).Scan(&config.ID, &config.TempMin, &config.TempMax, &config.ConsecutiveCount, &config.NotifyEnabled, &config.WebhookURL, &config.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取配置失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func UpdateAlertConfig(c *gin.Context) {
	var config models.AlertConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确：" + err.Error()})
		return
	}

	if config.TempMin >= config.TempMax {
		c.JSON(http.StatusBadRequest, gin.H{"error": "温度下限必须小于温度上限"})
		return
	}
	if config.ConsecutiveCount < 1 || config.ConsecutiveCount > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "连续次数必须在 1 到 10 之间"})
		return
	}
	if config.NotifyEnabled == 1 && config.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开启通知后必须填写 Webhook 地址"})
		return
	}

	_, err := database.DB.Exec(`
		UPDATE alert_config SET temp_min = ?, temp_max = ?, consecutive_count = ?, notify_enabled = ?, webhook_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = (SELECT id FROM alert_config ORDER BY id DESC LIMIT 1)
	`, config.TempMin, config.TempMax, config.ConsecutiveCount, config.NotifyEnabled, config.WebhookURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存配置失败：" + err.Error()})
		return
	}

	GetAlertConfig(c)
}

func getCurrentAlertConfig() (*models.AlertConfig, error) {
	var config models.AlertConfig
	err := database.DB.QueryRow(
		"SELECT id, temp_min, temp_max, consecutive_count, notify_enabled, COALESCE(webhook_url, ''), updated_at FROM alert_config ORDER BY id DESC LIMIT 1",
	).Scan(&config.ID, &config.TempMin, &config.TempMax, &config.ConsecutiveCount, &config.NotifyEnabled, &config.WebhookURL, &config.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &config, nil
}

func getLastNWaterRecords(tankID int, n int) ([]models.WaterRecord, error) {
	rows, err := database.DB.Query(
		"SELECT id, tank_id, temperature, pH, recorded_at FROM water_records WHERE tank_id = ? ORDER BY id DESC LIMIT ?",
		tankID, n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.WaterRecord
	for rows.Next() {
		var r models.WaterRecord
		err := rows.Scan(&r.ID, &r.TankID, &r.Temperature, &r.PH, &r.RecordedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	return records, nil
}

func checkTankAlert(tankID int, config *models.AlertConfig) *models.TankAlertStatus {
	records, err := getLastNWaterRecords(tankID, config.ConsecutiveCount)
	if err != nil || len(records) < config.ConsecutiveCount {
		lastTemps := make([]float64, 0)
		for _, r := range records {
			lastTemps = append([]float64{r.Temperature}, lastTemps...)
		}
		return &models.TankAlertStatus{
			TankID:     tankID,
			IsAlerting: false,
			LastTemps:  lastTemps,
		}
	}

	lastTemps := make([]float64, len(records))
	for i, r := range records {
		lastTemps[len(records)-1-i] = r.Temperature
	}

	allOverMax := true
	allUnderMin := true

	for _, temp := range lastTemps {
		if temp <= config.TempMax {
			allOverMax = false
		}
		if temp >= config.TempMin {
			allUnderMin = false
		}
	}

	if allOverMax {
		return &models.TankAlertStatus{
			TankID:     tankID,
			IsAlerting: true,
			AlertType:  "high",
			LastTemps:  lastTemps,
			Message:    "温度过高！",
		}
	}

	if allUnderMin {
		return &models.TankAlertStatus{
			TankID:     tankID,
			IsAlerting: true,
			AlertType:  "low",
			LastTemps:  lastTemps,
			Message:    "温度过低！",
		}
	}

	return &models.TankAlertStatus{
		TankID:     tankID,
		IsAlerting: false,
		LastTemps:  lastTemps,
	}
}

func GetAllTankAlerts(c *gin.Context) {
	config, err := getCurrentAlertConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取配置失败：" + err.Error()})
		return
	}

	tankRows, err := database.DB.Query("SELECT id, name FROM tanks ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取鱼缸列表失败：" + err.Error()})
		return
	}
	defer tankRows.Close()

	type tankWithStatus struct {
		ID     int                    `json:"tank_id"`
		Name   string                 `json:"tank_name"`
		Status models.TankAlertStatus `json:"status"`
	}

	var result []tankWithStatus
	for tankRows.Next() {
		var id int
		var name string
		if err := tankRows.Scan(&id, &name); err != nil {
			continue
		}
		status := checkTankAlert(id, config)
		result = append(result, tankWithStatus{
			ID:     id,
			Name:   name,
			Status: *status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"config": config,
		"alerts": result,
	})
}

func sendWebhook(webhookURL string, message string) error {
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func TriggerAlertNotify(c *gin.Context) {
	config, err := getCurrentAlertConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取配置失败：" + err.Error()})
		return
	}

	if config.NotifyEnabled != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "通知开关未开启，请先在设置里打开通知开关"})
		return
	}

	if config.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未配置 Webhook 地址"})
		return
	}

	tankRows, err := database.DB.Query("SELECT id, name FROM tanks ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取鱼缸列表失败：" + err.Error()})
		return
	}
	defer tankRows.Close()

	var alertMessages []string
	for tankRows.Next() {
		var id int
		var name string
		if err := tankRows.Scan(&id, &name); err != nil {
			continue
		}
		status := checkTankAlert(id, config)
		if status.IsAlerting {
			alertMessages = append(alertMessages,
				"【水温报警】"+name+"："+status.Message+
					" 最近温度："+tempsToString(status.LastTemps),
			)

			database.DB.Exec(
				"INSERT INTO alert_logs (tank_id, alert_type, message, notified) VALUES (?, ?, ?, 1)",
				id, status.AlertType, name+"："+status.Message,
			)
		}
	}

	if len(alertMessages) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "当前没有需要通知的报警", "sent": false})
		return
	}

	fullMessage := ""
	for _, msg := range alertMessages {
		fullMessage += msg + "\n"
	}

	if err := sendWebhook(config.WebhookURL, fullMessage); err != nil {
		log.Printf("发送 webhook 失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送通知失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "通知已发送",
		"sent":    true,
		"content": alertMessages,
	})
}

func tempsToString(temps []float64) string {
	result := ""
	for i, t := range temps {
		if i > 0 {
			result += "、"
		}
		result += floatToString(t) + "°C"
	}
	return result
}

func floatToString(f float64) string {
	s := strconv.FormatFloat(f, 'f', 1, 64)
	return s
}
