package actions

import (
	"goproject/storage"
	"slices"
)

type Actions struct {
	Db storage.IStorage
}

func (a *Actions) CreateNewVote(eventId string, person string, dates []string) bool {
	ch := make(chan bool)
	go a.Db.AddVote(ch, eventId, person, dates)
	return <-ch
}

func (a *Actions) CreateNewEvent(name string, proposedDates []string) string {
	ch := make(chan string)
	go a.Db.AddEvent(ch, name, proposedDates)
	return <-ch
}

func (a *Actions) GetEvent(eventId string) storage.Event {
	ch := make(chan storage.Event)
	go a.Db.GetEvent(ch, eventId)
	return <-ch
}

func (a *Actions) GetAllEvents() []storage.Event {
	ch := make(chan []storage.Event)
	go a.Db.GetAllEvents(ch)
	return <-ch
}

func (a *Actions) GetEventWithSuitableVotes(eventId string) (storage.Event, []SuitableDates) {
	ch := make(chan storage.Event)
	go a.Db.GetEvent(ch, eventId)
	event := <-ch

	voterNames := getAllVoterNames(event)
	suitableDates := getDatesWithAllVoters(voterNames, event)

	return event, suitableDates
}

func getAllVoterNames(event storage.Event) []string {

	var voterNames []string
	for _, vote := range getAllVotes(event) {
		if vote.Name != "" && !slices.Contains(voterNames, vote.Name) {
			voterNames = append(voterNames, vote.Name)
		}
	}
	return voterNames
}

func getDatesWithAllVoters(voters []string, event storage.Event) []SuitableDates {
	var datesWithMostVotes []SuitableDates
	for _, date := range event.ProposedDates {
		if len(date.Votes) == len(voters) {
			datesWithMostVotes = append(datesWithMostVotes,
				SuitableDates{Date: date.Date, People: voters})
		}
	}
	return datesWithMostVotes
}

func getAllVotes(event storage.Event) []storage.Vote {
	var allVotes []storage.Vote
	for _, date := range event.ProposedDates {
		allVotes = append(allVotes, date.Votes...)
	}
	return allVotes
}
