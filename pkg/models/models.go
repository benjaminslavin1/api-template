package models

import (
	"errors"
	"time"
)

var (
	ErrNoTestError = errors.New("example error")
)

type TestStruct struct {
	Created     time.Time `db:"created"`
	LastUpdated time.Time `db:"last_updated"`
}
