package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Item struct {
	Chrt_id      int    `json:"chrt_id"`
	Track_number string `json:"track_number"`
	Price        int    `json:"price"`
	Rid          string `json:"rid"`
	Name         string `json:"name"`
	Sale         int    `json:"sale"`
	Size         string `json:"size"`
	Total_price  int    `json:"total_price"`
	Nm_id        int    `json:"nm_id"`
	Brand        string `json:"brand"`
	Status       int    `json:"status"`
}

func (i Item) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Item) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion failed")
	}

	return json.Unmarshal(b, &i)
}
