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
