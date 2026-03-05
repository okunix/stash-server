package ports

import "encoding/json"

type ValidationError map[string]string

func (e ValidationError) Error() string {
	jsonBytes, _ := json.Marshal(e)
	return string(jsonBytes)
}
