package model

import "github.com/google/uuid"

type Country struct {
	Id uuid.UUID
	Name string
	Code string
}