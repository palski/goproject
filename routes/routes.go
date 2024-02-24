package routes

import (
	api_v1 "goproject/api/v1"
	"goproject/storage"

	"github.com/gorilla/mux"
)

func SetupRoutes(db storage.StorageInterface) *mux.Router {

	router := mux.NewRouter()
	handlerV1 := &api_v1.EventHandler{Db: db}
	router.HandleFunc("/api/v1/events", handlerV1.HandleEventsRequest)
	router.HandleFunc("/api/v1/events/{eventId}", handlerV1.HandleGetEventDetailsRequest)
	router.HandleFunc("/api/v1/events/{eventId}/results", handlerV1.HandleGetEventResultsRequest)
	router.HandleFunc("/api/v1/events/{eventId}/vote", handlerV1.HandleVoteEventRequest)

	return router
}
