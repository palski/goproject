package storage

type StorageInterface interface {
	InitializeDatabase()
	AddEvent(ch chan string, name string, proposedDates []string)
	GetAllEvents(ch chan []Event)
	GetEvent(ch chan Event, id string)
	AddVote(ch chan bool, eventId string, voterName string, date []string)
}
