package pkg

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)

	if err != nil {
		return "", fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return string(bytes), nil
}
