package storage

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlLiteDb struct {
}

func (d *SqlLiteDb) InitializeDatabase() {
	db := openDatabaseConnection()
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&ProposedDate{})
	db.AutoMigrate(&Vote{})
}

func (d *SqlLiteDb) AddEvent(ch chan string, name string, proposedDates []string) {
	defer close(ch)

	db := openDatabaseConnection()

	newUuid := uuid.New().String()
	event := Event{UUID: newUuid, Name: name}
	dates := createProposedDates(newUuid, proposedDates)
	success := addEventIntoDatabase(db, event, dates)

	if success {
		ch <- newUuid
	} else {
		ch <- ""
	}
}

func (d *SqlLiteDb) GetAllEvents(ch chan []Event) {

	defer close(ch)
	db := openDatabaseConnection()

	var events []Event
	db.Select("UUID", "Name").Find(&events)
	ch <- events
}

func (d *SqlLiteDb) GetEvent(ch chan Event, id string) {
	defer close(ch)

	db := openDatabaseConnection()

	var event Event
	db.Preload("ProposedDates.Votes").First(&event, "UUID = ?", id)
	ch <- event
}

func (d *SqlLiteDb) AddVote(ch chan bool, eventId string, voterName string, date []string) {
	defer close(ch)

	db := openDatabaseConnection()

	tx := db.Begin()

	// remove previous votes of the person
	result := tx.Delete(&Vote{}, "event_uuid = ? and name = ?", eventId, voterName)
	if result.Error != nil {
		tx.Rollback()
		ch <- false
		return
	}

	if !addNewVotes(tx, eventId, voterName, date) {
		tx.Rollback()
		ch <- false
		return
	}

	tx.Commit()
	ch <- true
}

func createProposedDates(eventUuid string, dates []string) []ProposedDate {
	var proposedDates []ProposedDate
	for _, date := range dates {
		proposedDates = append(proposedDates, ProposedDate{EventUUID: eventUuid, Date: date})
	}
	return proposedDates
}

func addNewVotes(tx *gorm.DB, eventId string, voterName string, date []string) bool {

	for _, d := range date {
		success := addVoteForProposedDate(tx, eventId, voterName, d)
		if !success {
			return false
		}
	}
	return true
}

func addVoteForProposedDate(
	tx *gorm.DB,
	eventId string,
	voterName string,
	date string) bool {

	var proposedDate ProposedDate
	result := tx.Where("event_uuid = ? and Date = ?", eventId, date).First(&proposedDate)
	if result.Error != nil {
		return false
	}
	result = tx.Create(&Vote{ProposedDateId: proposedDate.ID, EventUUID: eventId, Name: voterName})

	return result.Error == nil
}

func openDatabaseConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("event_database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func addProposedDatesToDatabase(
	tx *gorm.DB,
	event Event,
	proposedDates []ProposedDate) bool {

	for _, date := range proposedDates {
		result := tx.Create(&date)
		if result.Error != nil {
			return false
		}
	}

	return true
}

func addEventIntoDatabase(
	db *gorm.DB,
	event Event,
	proposedDates []ProposedDate) bool {

	tx := db.Begin()
	if tx.Error != nil {
		return false
	}

	ret := tx.Create(&event)
	if ret.Error != nil {
		tx.Rollback()
		return false
	}

	if !addProposedDatesToDatabase(tx, event, proposedDates) {
		tx.Rollback()
		return false
	}

	ret = tx.Commit()
	return ret.Error == nil
}
