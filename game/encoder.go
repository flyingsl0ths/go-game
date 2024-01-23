package game

import (
	"encoding/json"
	"os"
)

type Encoder[T any] struct {
	file    *os.File
	encoder *json.Encoder
}

func NewEncoder[T any](filePath string) (Encoder[T], error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0600)

	if err != nil {
		return Encoder[T]{}, err
	}

	encoder := json.NewEncoder(file)

	return Encoder[T]{
		file:    file,
		encoder: encoder,
	}, nil
}

func (e *Encoder[T]) Encode(source T) error {
	return e.encoder.Encode(source)
}

func (e *Encoder[T]) Close() error {
	return e.file.Close()
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
