package models

type Tank struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
}

type WaterRecord struct {
	ID          int     `json:"id"`
	TankID      int     `json:"tank_id"`
	Temperature float64 `json:"temperature"`
	PH          float64 `json:"ph"`
	RecordedAt  string  `json:"recorded_at"`
}

type BreedingRecord struct {
	ID                int    `json:"id"`
	TankID            int    `json:"tank_id"`
	Strain            string `json:"strain"`
	PairDate          string `json:"pair_date"`
	ExpectedBirthDate string `json:"expected_birth_date"`
	Notes             string `json:"notes"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
}

type AlertConfig struct {
	ID               int     `json:"id"`
	TempMin          float64 `json:"temp_min"`
	TempMax          float64 `json:"temp_max"`
	ConsecutiveCount int     `json:"consecutive_count"`
	NotifyEnabled    int     `json:"notify_enabled"`
	WebhookURL       string  `json:"webhook_url"`
	UpdatedAt        string  `json:"updated_at"`
}

type TankAlertStatus struct {
	TankID     int       `json:"tank_id"`
	IsAlerting bool      `json:"is_alerting"`
	AlertType  string    `json:"alert_type"`
	LastTemps  []float64 `json:"last_temps"`
	Message    string    `json:"message"`
}

type AlertLog struct {
	ID        int    `json:"id"`
	TankID    int    `json:"tank_id"`
	AlertType string `json:"alert_type"`
	Message   string `json:"message"`
	Notified  int    `json:"notified"`
	CreatedAt string `json:"created_at"`
}
