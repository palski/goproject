package actions

import "goproject/storage"

type IActions interface {
	GetEventWithSuitableVotes(eventId string) (storage.Event, []SuitableDates)
	CreateNewVote(eventId string, person string, dates []string) bool
	GetEvent(eventId string) storage.Event
	GetAllEvents() []storage.Event
	CreateNewEvent(name string, proposedDates []string) string
}

type SuitableDates struct {
	Date   string
	People []string
}
