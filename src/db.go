package henrymedscodechallenge

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type Database struct {
	db *memdb.MemDB
}

func CreateDB() *Database {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"availability": {
				Name: "availability",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ProviderID"},
					},
				},
			},
			"appointment": {
				Name: "appointment",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"provider_id": {
						Name:    "provider_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ProviderID"},
					},
					"client_id": {
						Name:    "client_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ClientID"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return &Database{
		db: db,
	}
}

func (database *Database) SaveAvailability(ID string, decoder *json.Decoder) (*Availability, error) {
	availability := &Availability{
		ProviderID: ID,
	}
	err := decoder.Decode(availability)
	if err != nil {
		return nil, err
	}
	txn := database.db.Txn(true)
	if err := txn.Insert("availability", availability); err != nil {
		panic(err)
	}
	txn.Commit()
	return availability, nil
}

func (database *Database) GetAvailability() []*Availability {
	availabilities := []*Availability{}
	txn := database.db.Txn(false)

	it, err := txn.Get("availability", "id")
	if err != nil {
		panic(err)
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		a := obj.(*Availability)

		appointments := database.GetAppointmentsForProvider(a.ProviderID)
		appointmentAvailabilities := splitAvailabilities(a, 15, appointments)

		availabilities = append(appointmentAvailabilities, a)
	}
	return availabilities
}

func (database *Database) ProviderIsAvailableForAppointment(appointment *Appointment) bool {
	available := false
	availabilities := database.GetAvailabilityForProvider(appointment.ProviderID)
	for _, a := range availabilities {
		if canBookAppointment(a.Start, time.Now(), appointment) {
			available = true
		}
	}
	return available

}

func (database *Database) GetAvailabilityForProvider(providerID string) []*Availability {
	txn := database.db.Txn(false)

	obj, err := txn.First("availability", "id", providerID)
	if err != nil {
		panic(err)
	}
	a := obj.(*Availability)

	appointments := database.GetAppointmentsForProvider(providerID)
	appointmentAvailabilities := splitAvailabilities(a, 15, appointments)

	return appointmentAvailabilities
}

func splitAvailabilities(a *Availability, interval int, appointments []*Appointment) []*Availability {
	availabilities := []*Availability{}

	currentStart := a.Start
	appointmentLength := time.Duration(interval * int(time.Minute))

	for currentStart.Before(a.End) {
		currentEnd := currentStart.Add(appointmentLength)
		if !dateOverlapsExistingAppointment(currentStart, appointments) {
			newAvailability := &Availability{
				Start:      currentStart,
				End:        currentEnd,
				ProviderID: a.ProviderID,
			}

			availabilities = append(availabilities, newAvailability)
		}
		currentStart = currentEnd
	}
	return availabilities
}

// If the proposed start date is on or after the appointment start and before the appointment end
func datesOverlapAppointment(startTime time.Time, a *Appointment) bool {
	return (startTime.Equal(a.Start) || startTime.After(a.Start)) && startTime.Before(a.End)
}

func dateOverlapsExistingAppointment(startTime time.Time, appointments []*Appointment) bool {
	for _, a := range appointments {
		if datesOverlapAppointment(startTime, a) {
			return true
		}
	}
	return false
}

// datesOverlapAppointment AND
// the appointment has not expired OR it has expired and is pending
func canBookAppointment(startTime, now time.Time, a *Appointment) bool {
	return datesOverlapAppointment(startTime, a) && (now.Before(a.Expires) || (now.After(a.Expires) && a.Status != Pending))
}

func (database *Database) GetAppointmentsForProvider(providerID string) []*Appointment {
	txn := database.db.Txn(false)

	it, err := txn.Get("appointment", "provider_id", providerID)
	if err != nil {
		panic(err)
	}
	appointments := []*Appointment{}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		a := obj.(*Appointment)

		appointments = append(appointments, a)
	}

	return appointments

}

func (database *Database) GetAppointmentByIDForClient(appointmentID string) *Appointment {
	txn := database.db.Txn(false)

	obj, err := txn.First("appointment", "id", appointmentID)
	if err != nil {
		panic(err)
	}
	if obj == nil {
		return nil
	}
	a := obj.(*Appointment)

	return a

}

func (database *Database) validateAppointment(a *Appointment) error {
	var err error
	// Appointments cannot being less than 24 hours from now
	earliestStartDate := time.Now().Add(time.Hour * 24)
	if a.Start.Before(earliestStartDate) {
		err = errors.New("appointments can only be booked 24 hours or more in advance")
	}
	// The provider must have availability for the requested appointment
	if !database.ProviderIsAvailableForAppointment(a) {
		err = errors.New("the provider is no longer available for this appointments")
	}
	return err
}

func (database *Database) SaveAppointment(ID string, decoder *json.Decoder) (*Appointment, error) {
	appointment := &Appointment{
		ID:       uuid.New().String(),
		ClientID: ID,
		Status:   Pending,
		Expires:  time.Now().Add(time.Minute * 30),
	}
	err := decoder.Decode(appointment)
	if err != nil {
		return nil, err
	}

	err = database.validateAppointment(appointment)
	if err != nil {
		return nil, err
	}

	txn := database.db.Txn(true)
	if err := txn.Insert("appointment", appointment); err != nil {
		panic(err)
	}
	txn.Commit()
	return appointment, nil
}

func (database *Database) ConfirmAppointment(appointmentID string) (*Appointment, error) {
	appointment := database.GetAppointmentByIDForClient(appointmentID)
	if appointment == nil {
		return nil, fmt.Errorf("no appointment was found with ID %s", appointmentID)
	}
	appointment.Status = Confirmed
	txn := database.db.Txn(true)
	if err := txn.Insert("appointment", appointment); err != nil {
		panic(err)
	}
	txn.Commit()
	return appointment, nil
}
