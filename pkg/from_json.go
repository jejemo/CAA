package pkg

import (
	"encoding/json"
	"fmt"
)

func FromJSON(jsonStr string, v interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), v)
	if err != nil {
		return fmt.Errorf("error unmarshaling from JSON: %v", err)
	}
	return nil
}
