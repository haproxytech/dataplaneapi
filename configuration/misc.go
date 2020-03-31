package configuration

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func decodeBootstrapKey(key string) ([]string, error) {
	base64Data := strings.Split(key, ".")
	if len(base64Data) != 8 {
		return nil, fmt.Errorf("bottstrap key in unrecognized format")
	}
	result := make([]string, len(base64Data))
	for i, val := range base64Data {
		raw, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, fmt.Errorf("%s - %w", val, err)
		}
		result[i] = string(raw)
	}
	return result, nil
}
