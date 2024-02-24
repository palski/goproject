package storage

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID            uint
	UUID          string `gorm:"primaryKey"`
	Name          string
	ProposedDates []ProposedDate
}

type ProposedDate struct {
	gorm.Model
	ID        uint
	EventUUID string
	Date      string
	Votes     []Vote
}

type Vote struct {
	gorm.Model
	Name           string
	ProposedDateId uint
	EventUUID      string
}
