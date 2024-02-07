package henrymedscodechallenge

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Creates router and maps paths to functions
func GetRouter(database *Database, logger *log.Logger) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Request received on /")
		json.NewEncoder(w).Encode("This is working")
	}).Methods("GET")

	router.HandleFunc("/provider/{id}/availability", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ID := mux.Vars(r)["id"]
		availability, err := database.SaveAvailability(ID, json.NewDecoder(r.Body))
		if err != nil {
			logger.Println("Error saving Availability", err)
		} else {
			json.NewEncoder(w).Encode(availability)
		}

	}).Methods("POST")

	router.HandleFunc("/provider/availability", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		availabilities := database.GetAvailability()
		json.NewEncoder(w).Encode(availabilities)

	}).Methods("GET")

	router.HandleFunc("/client/{id}/appointment", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ID := mux.Vars(r)["id"]
		appointment, err := database.SaveAppointment(ID, json.NewDecoder(r.Body))
		if err != nil || appointment == nil {
			logger.Println("Error saving Appointment", err)
			w.WriteHeader(http.StatusInternalServerError) // Probably better to return a bad request here
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(appointment)
		}
	}).Methods("POST")

	router.HandleFunc("/client/{id}/appointment/{appointmentId}/confirm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		appointmentID := vars["appointmentId"]

		appointment, err := database.ConfirmAppointment(appointmentID)
		if err != nil || appointment == nil {
			logger.Println("Error confirming Appointment", err)
			w.WriteHeader(http.StatusInternalServerError) // Probably better to return a bad request here
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(appointment)
		}

	}).Methods("PUT")

	return router
}
