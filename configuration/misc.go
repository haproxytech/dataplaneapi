package configuration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func decodeBootstrapKey(key string) (map[string]string, error) {
	raw, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", key, err)
	}
	var decodedKey map[string]string
	err = json.Unmarshal(raw, &decodedKey)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", key, err)
	}
	return decodedKey, nil
}
