package game

import (
	"encoding/json"
	"os"
)

type Encoder struct {
	file    *os.File
	encoder *json.Encoder
}

func NewEncoder(filePath string) (Encoder, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0600)

	if err != nil {
		return Encoder{}, err
	}

	encoder := json.NewEncoder(file)

	return Encoder{
		file:    file,
		encoder: encoder,
	}, nil
}

func (e *Encoder) Encode(source any) error {
	return e.encoder.Encode(source)
}

func (e *Encoder) Close() error {
	return e.file.Close()
}
