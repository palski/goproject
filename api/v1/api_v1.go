package api_v1

import (
	"encoding/json"
	"goproject/actions"
	"goproject/storage"
	"net/http"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	Actions actions.IActions
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

	event, suitableDates := e.Actions.GetEventWithSuitableVotes(eventId)
	if event.UUID == "" {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(
		EventResultDTO{
			Id:            event.UUID,
			Name:          event.Name,
			SuitableDates: mapToDto(suitableDates),
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

	success := e.Actions.CreateNewVote(eventId, voteDto.Name, voteDto.Votes)
	if !success {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	event := e.Actions.GetEvent(eventId)

	writeEventToResponse(response, event)
}

func (e *EventHandler) HandleEventsRequest(response http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		createNewEvent(response, request, e.Actions)
	} else if request.Method == "GET" {
		getAllvents(response, request, e.Actions)
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

	event := e.Actions.GetEvent(eventId)

	writeEventToResponse(response, event)
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

func getAllvents(response http.ResponseWriter, request *http.Request, actions actions.IActions) {

	events := mapToSimpleEventDTO(actions.GetAllEvents())

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(events)
}

func createNewEvent(response http.ResponseWriter, request *http.Request, action actions.IActions) {
	var eventDto NewEventDTO
	var error = json.NewDecoder(request.Body).Decode(&eventDto)
	if error != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	newUuid := action.CreateNewEvent(eventDto.Name, eventDto.Dates)

	if newUuid != "" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(EventCreatedDTO{Id: newUuid})
	} else {
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func mapToDto(suitableDates []actions.SuitableDates) []ProposedDateVotesDTO {
	var datesWithMostVotes []ProposedDateVotesDTO
	for _, suitableDate := range suitableDates {
		datesWithMostVotes = append(datesWithMostVotes,
			ProposedDateVotesDTO{Date: suitableDate.Date, People: suitableDate.People})
	}
	return datesWithMostVotes
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
