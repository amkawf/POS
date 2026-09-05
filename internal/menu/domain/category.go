package domain

import "github.com/google/uuid"

type Category struct {
	ID        uuid.UUID
	MenuID    uuid.UUID
	Name      string
	SortOrder int32
	Status    string
}
