package actions

import (
	"goproject/storage"
	"testing"
)

type MockStorage struct {
	event storage.Event
}

func (m *MockStorage) GetEvent(ch chan storage.Event, eventId string) {
	ch <- m.event
}

func (m *MockStorage) InitializeDatabase() {}

func (m *MockStorage) AddEvent(ch chan string, name string, proposedDates []string) {
}

func (m *MockStorage) GetAllEvents(ch chan []storage.Event) {
}

func (m *MockStorage) AddVote(ch chan bool, eventId string, voterName string, date []string) {
}

func GetSuitableDate(t *testing.T) {
	mockStorage := &MockStorage{
		event: storage.Event{
			UUID: "1",
			Name: "Villes Birthday",
			ProposedDates: []storage.ProposedDate{
				{Date: "2025-05-19", Votes: []storage.Vote{{Name: "Jussi"}, {Name: "Kalle"}}},
				{Date: "2025-05-20", Votes: []storage.Vote{{Name: "Jussi"}, {Name: "Kalle"}, {Name: "Anssi"}}},
			},
		},
	}
	actions := &Actions{Db: mockStorage}

	event, suitableDates := actions.GetEventWithSuitableVotes("1")

	if event.UUID != "1" {
		t.Errorf("Invalid event")
	}

	if len(suitableDates) != 1 {
		t.Errorf("invalid suitable date count")
	}

	if suitableDates[0].Date != "2025-05-20" {
		t.Errorf("Invliad suitable date")
	}

	if len(suitableDates[0].People) != 3 {
		t.Errorf("Invalid suitable date people count")
	}
}
