package game

import (
	"encoding/json"
	"os"
)

type PlayerScore struct {
	Player string `json:"player"`
	Score  uint32 `json:"score"`
}

func DecodeInto[T any](filePath string, source T) (T, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return source, err
	}
	err = json.Unmarshal(dat, &source)
	if err != nil {
		return source, err
	}

	return source, nil
}
