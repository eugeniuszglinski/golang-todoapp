package core_http_types

import (
	"encoding/json"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

// Nullable in core_http_types is a type used only for transport layer purposes
type Nullable[T any] struct {
	domain.Nullable[T]
}

// UnmarshalJSON used to implement Unmarshaler interface with a custom method for the Nullable[T any]
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {

	// if UnmarshalJSON was called, this means that some value in JSON was received
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{Value: n.Value, Set: n.Set}
}
