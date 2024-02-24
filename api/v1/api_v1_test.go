package api_v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"goproject/storage"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
)

type MockDb struct{}

func (db *MockDb) InitializeDatabase() {
}
func (db *MockDb) AddEvent(ch chan string, name string, proposedDates []string) {
}
func (db *MockDb) GetAllEvents(ch chan []storage.Event) {
}

func (db *MockDb) GetEvent(ch chan storage.Event, id string) {

	proposedDates := make([]storage.ProposedDate, 2)
	proposedDates[0] = storage.ProposedDate{Date: "2021-01-01",
		Votes: []storage.Vote{{Name: "Ville", ProposedDateId: 1, EventUUID: "1"},
			{Name: "Kalle", ProposedDateId: 1, EventUUID: "1"}}}
	proposedDates[1] = storage.ProposedDate{Date: "2021-01-02",
		Votes: []storage.Vote{
			{Name: "Ville", ProposedDateId: 2, EventUUID: "1"},
			{Name: "Kalle", ProposedDateId: 2, EventUUID: "1"},
			{Name: "Joonatan", ProposedDateId: 2, EventUUID: "1"}}}

	ch <- storage.Event{UUID: "1", Name: "Test Event", ProposedDates: proposedDates}
}

func (db *MockDb) AddVote(ch chan bool, eventId string, voterName string, date []string) {
}

func TestHandleGetResults(t *testing.T) {
	handler := &EventHandler{Db: &MockDb{}}

	t.Run("Get date with all voters", func(t *testing.T) {

		mux.NewRouter()
		request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/events/1/results", nil), map[string]string{"eventId": "1"})
		response := httptest.NewRecorder()

		handler.HandleGetEventResultsRequest(response, request)

		var result EventResultDTO
		json.Unmarshal(response.Body.Bytes(), &result)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "1", result.Id)
		assert.Equal(t, 1, len(result.SuitableDates))
		assert.Equal(t, "2021-01-02", result.SuitableDates[0].Date)
		assert.Equal(t, 3, len(result.SuitableDates[0].People))
	})
}
