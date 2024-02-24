package routes

import (
	"goproject/actions"
	api_v1 "goproject/api/v1"

	"github.com/gorilla/mux"
)

func SetupRoutes(actions actions.IActions) *mux.Router {

	router := mux.NewRouter()
	handlerV1 := &api_v1.EventHandler{Actions: actions}
	router.HandleFunc("/api/v1/events", handlerV1.HandleEventsRequest)
	router.HandleFunc("/api/v1/events/{eventId}", handlerV1.HandleGetEventDetailsRequest)
	router.HandleFunc("/api/v1/events/{eventId}/results", handlerV1.HandleGetEventResultsRequest)
	router.HandleFunc("/api/v1/events/{eventId}/vote", handlerV1.HandleVoteEventRequest)

	return router
}
