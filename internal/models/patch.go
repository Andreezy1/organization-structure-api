package models

import "encoding/json"

type Patch[T any] struct {
	Value T
	Set   bool
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	p.Set = true
	return json.Unmarshal(data, &p.Value)
}
