package watcher

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Metadata struct {
	Names []string `json:"names"`
}

func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Metadata) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not slice bytes")
	}
	var metadata Metadata
	if err := json.Unmarshal(bytes, &metadata); err != nil {
		return err
	}
	*m = metadata
	return nil
}
