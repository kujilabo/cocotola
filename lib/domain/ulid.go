package domain

import (
	"github.com/oklog/ulid/v2"
)

func NewULID() string {
	return ulid.Make().String()
}
