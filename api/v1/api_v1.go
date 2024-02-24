package api_v1

import (
	"encoding/json"
	"goproject/storage"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	Db storage.StorageInterface
}

func (e *EventHandler) HandleGetEventResultsRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	eventId := getEventId(request)
	if eventId == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	ch := make(chan storage.Event)
	go e.Db.GetEvent(ch, eventId)
	event := <-ch

	if event.UUID == "" {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	allVoters := getAllVoterNames(event)
	suitableDates := getDatesWithAllVoters(allVoters, event)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(
		EventResultDTO{
			Id:            event.UUID,
			Name:          event.Name,
			SuitableDates: suitableDates,
		},
	)
}

func (e *EventHandler) HandleVoteEventRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var voteDto NewVoteDTO
	var error = json.NewDecoder(request.Body).Decode(&voteDto)
	if error != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	eventId := getEventId(request)
	if eventId == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	ch := make(chan bool)
	go e.Db.AddVote(ch, eventId, voteDto.Name, voteDto.Votes)
	if !<-ch {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventChannel := make(chan storage.Event)
	go e.Db.GetEvent(eventChannel, eventId)
	writeEventToResponse(response, <-eventChannel)
}

func (e *EventHandler) HandleEventsRequest(response http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		createNewEvent(response, request, e.Db)
	} else if request.Method == "GET" {
		getAllvents(response, request, e.Db)
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (e *EventHandler) HandleGetEventDetailsRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	eventId := getEventId(request)
	if eventId == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	ch := make(chan storage.Event)
	go e.Db.GetEvent(ch, eventId)
	writeEventToResponse(response, <-ch)
}

func writeEventToResponse(response http.ResponseWriter, event storage.Event) {
	if event.UUID != "" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(
			EventDTO{
				Id:    event.UUID,
				Name:  event.Name,
				Dates: mapProposedDatesToStrings(event.ProposedDates),
				Votes: mapVotesToDTO(event.ProposedDates),
			},
		)
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}

func getAllvents(response http.ResponseWriter, request *http.Request, db storage.StorageInterface) {

	ch := make(chan []storage.Event)
	go db.GetAllEvents(ch)
	events := mapToSimpleEventDTO(<-ch)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(events)
}

func createNewEvent(response http.ResponseWriter, request *http.Request, db storage.StorageInterface) {
	var eventDto NewEventDTO
	var error = json.NewDecoder(request.Body).Decode(&eventDto)
	if error != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	ch := make(chan string)
	go db.AddEvent(ch, eventDto.Name, eventDto.Dates)
	newUuid := <-ch

	if newUuid != "" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(EventCreatedDTO{Id: newUuid})
	} else {
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func getDatesWithAllVoters(voters []string, event storage.Event) []ProposedDateVotesDTO {
	var datesWithMostVotes []ProposedDateVotesDTO
	for _, date := range event.ProposedDates {
		if len(date.Votes) == len(voters) {
			datesWithMostVotes = append(datesWithMostVotes,
				ProposedDateVotesDTO{Date: date.Date, People: voters})
		}
	}
	return datesWithMostVotes
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

func getAllVotes(event storage.Event) []storage.Vote {
	var allVotes []storage.Vote
	for _, date := range event.ProposedDates {
		allVotes = append(allVotes, date.Votes...)
	}
	return allVotes
}

func getDatesWithVotes(dates []storage.ProposedDate) []storage.ProposedDate {
	var datesWithVotes []storage.ProposedDate
	for _, date := range dates {
		for _, vote := range date.Votes {
			if vote.Name != "" {
				datesWithVotes = append(datesWithVotes, date)
				break
			}
		}
	}
	return datesWithVotes
}

func getPeopleWhoVoted(votes []storage.Vote) []string {
	var people []string
	for _, vote := range votes {
		people = append(people, vote.Name)
	}
	return people
}

func mapVotesToDTO(dates []storage.ProposedDate) []ProposedDateVotesDTO {
	var votes []ProposedDateVotesDTO
	for _, vote := range getDatesWithVotes(dates) {
		people := getPeopleWhoVoted(vote.Votes)
		votes = append(votes, ProposedDateVotesDTO{Date: vote.Date, People: people})
	}
	return votes
}

func getEventId(request *http.Request) string {
	vars := mux.Vars(request)
	return vars["eventId"]
}

func mapProposedDatesToStrings(dates []storage.ProposedDate) []string {
	var timeDates []string
	for _, date := range dates {
		timeDates = append(timeDates, date.Date)
	}
	return timeDates
}

func mapToSimpleEventDTO(events []storage.Event) []SimpleEventDTO {
	var eventDtos []SimpleEventDTO
	if events == nil {
		return eventDtos
	}

	for _, item := range events {
		eventDtos = append(eventDtos, SimpleEventDTO{
			Id:   item.UUID,
			Name: item.Name,
		})
	}
	return eventDtos
}
