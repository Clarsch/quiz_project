package dbdto

import (
	"encoding/json"
	"fmt"
)

type CategoryIncomingDTO struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c CategoryIncomingDTO) String() string {
	// Convert the Category struct to a JSON string for better visibility
	cJSON, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting Question to JSON: %v", err)
	}
	return string(cJSON)
}
