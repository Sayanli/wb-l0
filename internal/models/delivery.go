package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Delivery struct {
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Region  string `json:"region"`
	Address string `json:"address"`
}

func (d Delivery) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Delivery) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &d)
}
