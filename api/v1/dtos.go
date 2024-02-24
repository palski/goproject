package api_v1

type NewEventDTO struct {
	Id    string
	Name  string
	Dates []string
}

type EventCreatedDTO struct {
	Id string `json:"id"`
}

type SimpleEventDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EventResultDTO struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	SuitableDates []ProposedDateVotesDTO `json:"suitableDates"`
}

type EventDTO struct {
	Id    string                 `json:"id"`
	Name  string                 `json:"name"`
	Dates []string               `json:"dates"`
	Votes []ProposedDateVotesDTO `json:"votes"`
}

type ProposedDateVotesDTO struct {
	Date   string   `json:"date"`
	People []string `json:"people"`
}

type NewVoteDTO struct {
	Name  string
	Votes []string
}
