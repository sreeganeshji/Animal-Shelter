package rest

import (
	"back-end/models"
	"back-end/service"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type handler struct {
	Service *service.Service
}

// Init initializes the REST router
func Init(service service.Service) *mux.Router {

	handler := handler{Service: &service}

	router := mux.NewRouter()

	// Get Methods
	router.HandleFunc("/api/volunteer", handler.GetVolunteers).Methods(http.MethodGet)
	router.HandleFunc("/api/applicant", handler.GetApplicants).Methods(http.MethodGet)
	router.HandleFunc("/api/applicant/{email}", handler.GetApplicantByEmail).Methods(http.MethodGet)
	router.HandleFunc("/api/dog", handler.GetDogs).Methods(http.MethodGet)
	router.HandleFunc("/api/dog/{id}", handler.GetDog).Methods(http.MethodGet)
	router.HandleFunc("/api/application", handler.GetApplications).Methods(http.MethodGet)
	router.HandleFunc("/api/breed", handler.GetBreeds).Methods(http.MethodGet)
	router.HandleFunc("/api/dogbreed", handler.GetDogBreeds).Methods(http.MethodGet)
	router.HandleFunc("/api/expense", handler.GetExpenses).Methods(http.MethodGet)
	router.HandleFunc("/api/adoption", handler.GetAdoptions).Methods(http.MethodGet)
	router.HandleFunc("/api/volunteer", handler.GetVolunteers).Methods(http.MethodGet)
	router.HandleFunc("/api/expense-analysis", handler.ExpenseAnalysis).Methods(http.MethodGet)
	router.HandleFunc("/api/animal-control-report", handler.AnimalControlReport).Methods(http.MethodGet)
	router.HandleFunc("/api/monthly-adoption-report", handler.MonthlyAdoptionReport).Methods(http.MethodGet)
	router.HandleFunc("/api/animal-control-report-drilldown-surrendered", handler.AnimalControlReportDrillDownOne).Methods(http.MethodGet)
	router.HandleFunc("/api/animal-control-report-drilldown-adopted", handler.AnimalControlReportDrillDownTwo).Methods(http.MethodGet)
	router.HandleFunc("/api/logout", handler.Logout).Methods(http.MethodGet)

	// Post Methods
	router.HandleFunc("/api/volunteer", handler.CreateVolunteer).Methods(http.MethodPost)
	router.HandleFunc("/api/applicant", handler.CreateApplicant).Methods(http.MethodPost)
	router.HandleFunc("/api/dog", handler.CreateDog).Methods(http.MethodPost)
	router.HandleFunc("/api/application", handler.CreateApplication).Methods(http.MethodPost)
	router.HandleFunc("/api/breed", handler.CreateBreed).Methods(http.MethodPost)
	router.HandleFunc("/api/dogbreed", handler.CreateDogBreed).Methods(http.MethodPost)
	router.HandleFunc("/api/expense", handler.CreateExpense).Methods(http.MethodPost)
	router.HandleFunc("/api/adoption", handler.CreateAdoption).Methods(http.MethodPost)
	router.HandleFunc("/api/login", handler.Login).Methods(http.MethodPost, http.MethodOptions)

	// Put Methods
	router.HandleFunc("/api/application/{id}", handler.ChangeApplicationStatus).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/dog/{id}", handler.UpdateDog).Methods(http.MethodPut, http.MethodOptions)

	// Static files
	ui := uiHandler{staticPath: "build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(ui).Methods(http.MethodGet)

	return router
}

// CreateVolunteer accepts a JSON payload from which to create a new volunteer
func (h *handler) CreateVolunteer(w http.ResponseWriter, r *http.Request) {

	// Unmarshal JSON object into Volunteer struct
	decoder := json.NewDecoder(r.Body)
	var volunteer models.Volunteer
	err := decoder.Decode(&volunteer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Volunteer struct to create a new volunteer
	insertedVolunteer, err := h.Service.CreateVolunteer(volunteer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedVolunteer)
}

func RangeOfMonthOffset(o int) (string, string) {
	// get first day and last day of a month offset from the current month by o.
	t := time.Now()
	t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	//add months
	t = t.AddDate(0, o, 0)
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonthFirstDay := firstDay.AddDate(0, 1, 0)
	lastDay := nextMonthFirstDay.AddDate(0, 0, -1)

	return fmt.Sprintf("%d-%02d-%02d", firstDay.Year(), firstDay.Month(), firstDay.Day()), fmt.Sprintf("%d-%02d-%02d", lastDay.Year(), lastDay.Month(), lastDay.Day())
}

// Monthly analysis returns the summary for a given month
func (h *handler) MonthlyAdoptionReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	month := r.URL.Query().Get("month")
	var err error
	var monthlyAdoptionReport []models.MonthlyAdoptionReportItem
	var monthVal, _ = strconv.Atoi(month)
	//monthVal = -10

	startDate, endDate := RangeOfMonthOffset(monthVal)

	monthlyAdoptionReport, err = h.Service.MonthlyAdoptionReport(startDate, endDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// Return volunteers as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAdoptionReport)
}

// ExpenseAnalysis returns the expense analysis
func (h *handler) ExpenseAnalysis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	// Get volunteers with a first or last name like the provided string
	expenseAnalysis, err := h.Service.ExpenseAnalysis()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// Return volunteers as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenseAnalysis)
}

// GetVolunteer returns a list of all of the volunteers in the database
func (h *handler) GetVolunteers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	like := r.URL.Query().Get("like")
	var err error
	var volunteers []models.Volunteer

	if like == "" {
		// Get all volunteers
		volunteers, err = h.Service.GetVolunteers()
	} else {
		// Get volunteers with a first or last name like the provided string
		volunteers, err = h.Service.GetVolunteersLike(like)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// Return volunteers as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(volunteers)
}

func (h *handler) GetApplicantByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	email := mux.Vars(r)["email"]
	applicant, err := h.Service.GetApplicantByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicant)
}

func (h *handler) GetApplicants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	like := r.URL.Query().Get("like")
	status := r.URL.Query().Get("status")
	var applicants []models.Applicant
	var err error
	if status == "approved" {
		applicants, err = h.Service.GetApprovedApplicants()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting approved applicants from database: ", err)))
			return
		}
	} else if like != "" {
		applicants, err = h.Service.GetApprovedApplicantsLike(like)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
			return
		}
	} else {
		// get all applicants
		applicants, err = h.Service.GetApplicants()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicants)
}

func (h *handler) GetDogs(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	var dogs []models.Dog

	// query query for "current" param
	currentOnly, err := strconv.ParseBool(r.URL.Query().Get("current"))
	if err != nil || !currentOnly {
		// return ALL dogs
		dogs, err = h.Service.GetDogs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
			return
		}
	} else {
		// only return dogs currently in the shelter
		dogs, err = h.Service.GetCurrentDogs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
			return
		}
	}

	for i, dog := range dogs {
		dog.Breed, err = h.Service.GetBreedsByDogID(dog.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error retrieving record from database: ", err)))
			return
		}
		dogs[i] = dog
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dogs)
}

func (h *handler) GetDog(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	// query query for "id" param
	dogID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error parse dog ID from URL: ", err)))
		return
	}

	dog, err := h.Service.GetDog(dogID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	dog.Breed, err = h.Service.GetBreedsByDogID(dog.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error retrieving record from database: ", err)))
		return
	}

	dog.Expenses, err = h.Service.GetExpensesByDogID(dog.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error retrieving record from database: ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dog)
}

func (h *handler) GetApplications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	var applications []models.Application
	applicationStatus := r.URL.Query().Get("status")
	var err error
	if applicationStatus == "approved" {
		applications, err = h.Service.GetApprovedApplications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting applications from database: ", err)))
			return
		}
		for i, application := range applications {
			applicant, err := h.Service.GetApplicantByID(application.ApplicantIdFk)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprint("Error retrieving applicant from database: ", err)))
				return
			}
			applications[i].Applicant = applicant
		}
	} else if applicationStatus == "pending" {
		applications, err = h.Service.GetPendingApplications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error getting applications from database: ", err)))
			return
		}
		for i, application := range applications {
			applicant, err := h.Service.GetApplicantByID(application.ApplicantIdFk)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprint("Error retrieving applicant from database: ", err)))
				return
			}
			applications[i].Applicant = applicant
		}
	} else {
		applicationID, err := strconv.Atoi(r.URL.Query().Get("applicantId"))
		coApplicantFirstName := r.URL.Query().Get("coApplicantFirstName")
		coApplicantLastName := r.URL.Query().Get("coApplicantLastName")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error parsing applicantId from URL: ", err)))
			return
		}

		applications, err := h.Service.GetLatestApprovedApplication(applicationID, coApplicantFirstName, coApplicantLastName)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error retrieving lasted approved applicant application from database: ", err)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(applications)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}

func (h *handler) ChangeApplicationStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}
	id := mux.Vars(r)["id"]
	approve, err := strconv.ParseBool(r.URL.Query().Get("approve"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error parsing \"approve\" variable from URL: ", err)))
		return
	}

	err = h.Service.ChangeApplicationStatus(id, approve)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error updating record in database: ", err)))
		return
	}

	// Return 200 by default
}

func (h *handler) GetBreeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	breeds, err := h.Service.GetBreeds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breeds)
}

func (h *handler) GetDogBreeds(w http.ResponseWriter, r *http.Request) {

	dogBreeds, err := h.Service.GetDogBreeds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dogBreeds)
}

func (h *handler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.Service.GetExpenses()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (h *handler) GetAdoptions(w http.ResponseWriter, r *http.Request) {
	adoptions, err := h.Service.GetAdoptions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adoptions)
}

// CreateApplicant accepts a JSON payload from which to create a new Applicant
func (h *handler) CreateApplicant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var applicant models.Applicant
	err := decoder.Decode(&applicant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Volunteer struct to create a new volunteer
	insertedApplicant, err := h.Service.CreateApplicant(applicant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedApplicant)
}

// CreateDog accepts a JSON payload from which to create a new dog
func (h *handler) CreateDog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var dog models.Dog
	err := decoder.Decode(&dog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Dog struct to create a new dog
	insertedDog, err := h.Service.CreateDog(dog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	for _, b := range dog.Breed {
		breedID, err := h.Service.GetBreedID(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
			return
		}
		dogBreed := models.DogBreed{DogIdFk: insertedDog.ID, BreedIdFk: breedID}
		_, err = h.Service.CreateDogBreed(dogBreed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
			return
		}
		insertedDog.Breed = append(insertedDog.Breed, b)
	}

	insertedDog.Expenses = make([]models.Expense, 0)

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedDog)
}

// UpdateDog accepts a JSON payload from which to create a new dog
func (h *handler) UpdateDog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	// Get dogID
	dogID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error parsing dogID: ", err)))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var dog models.Dog
	err = decoder.Decode(&dog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}
	dog.ID = dogID

	// Use Dog struct to update an existing dog
	err = h.Service.UpdateDog(dog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Get the existing breeds for specified dog
	dogBreeds, err := h.Service.GetBreedsByDogID(dogID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// If multiple existing dogbreeds, do nothing
	if len(dogBreeds) == 1 {
		// If single existing dogbreed other than Mixed or Unknown, do nothing
		if dogBreeds[0] == "Mixed" || dogBreeds[0] == "Unknown" {
			// If single existing dogbreed is Mixed or Unknown and submitted breed is same, no nothing
			if dogBreeds[0] != dog.Breed[0] {
				// else delete existing dogbreed and create new for each breed submitted
				err = h.Service.DeleteDogBreedByDogID(dogID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(fmt.Sprint("Error updating old breed: ", err)))
					return
				}
				for _, b := range dog.Breed {
					breedID, err := h.Service.GetBreedID(b)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
						return
					}
					dogBreed := models.DogBreed{DogIdFk: dogID, BreedIdFk: breedID}
					_, err = h.Service.CreateDogBreed(dogBreed)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
						return
					}
				}
			}
		}
	}

	// Return a 204 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// CreateApplication accepts a JSON payload from which to create a new dog
func (h *handler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var application models.Application
	err := decoder.Decode(&application)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	date := time.Now().String()
	application.Date = string(date[0:10])
	application.State = "pending approval"

	// Use Volunteer struct to create a new volunteer
	insertedApplication, err := h.Service.CreateApplication(application)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedApplication)
}

// CreateBreed
// CreateDog accepts a JSON payload from which to create a new dog
func (h *handler) CreateBreed(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var breed models.Breed
	err := decoder.Decode(&breed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Volunteer struct to create a new volunteer
	insertedBreed, err := h.Service.CreateBreed(breed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedBreed)
}

// CreateDogBreed
// CreateDog accepts a JSON payload from which to create a new dog
func (h *handler) CreateDogBreed(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dogBreed models.DogBreed
	err := decoder.Decode(&dogBreed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Volunteer struct to create a new volunteer
	insertedDogBreed, err := h.Service.CreateDogBreed(dogBreed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedDogBreed)
}

// CreateExpense
// CreateDog accepts a JSON payload from which to create a new dog
func (h *handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var expense models.Expense
	err := decoder.Decode(&expense)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Volunteer struct to create a new volunteer
	insertedExpense, err := h.Service.CreateExpense(expense)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedExpense)
}

// CreateAdoption
// CreateDog accepts a JSON payload from which to create a new dog
func (h *handler) CreateAdoption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var adoption models.Adoption
	err := decoder.Decode(&adoption)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error unmarshalling JSON: ", err)))
		return
	}

	// Use Volunteer struct to create a new volunteer
	insertedAdoption, err := h.Service.CreateAdoption(adoption)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error saving record to database: ", err)))
		return
	}

	// Return a 201 if successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedAdoption)
}

// Login stuff

// UserLogin type
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login accepts a Form(email,password) payload,
// and checks DB for matching volunteer
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	var userLogin UserLogin
	json.NewDecoder(r.Body).Decode(&userLogin)

	var email = userLogin.Email
	if len(email) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Please provide valid email address.")))
		return
	}
	var password = userLogin.Password
	if len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Please provide valid password.")))
		return
	}

	// Use Volunteer struct to create a new volunteer
	volunteer, err := h.Service.GetVolunteerByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	if volunteer.Password == password {
		cookie := createCookie(volunteer.IsTrusted)
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(volunteer)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprint("No such email/password combination found. Please try re-entering your information.")))
}

func createCookie(isTrusted bool) http.Cookie {
	// add sessionKey to cookie
	var unparsedAttributePairs bytes.Buffer
	unparsedAttributePairs.WriteString("isTrustedVolunteer=")
	unparsedAttributePairs.WriteString(strconv.FormatBool(isTrusted))
	cookie := http.Cookie{
		Name:       "isTrustedVolunteer",
		Value:      strconv.FormatBool(isTrusted),
		Path:       "/",
		Domain:     "localhost",
		Expires:    time.Now().AddDate(0, 0, 1),
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		Raw:        "",
		Unparsed:   []string{unparsedAttributePairs.String()}, // Raw text of unparsed attribute-value pairs
	}
	return cookie
}

func createEmptyCookie() http.Cookie {
	// add sessionKey to cookie
	var unparsedAttributePairs bytes.Buffer
	unparsedAttributePairs.WriteString("isTrustedVolunteer=\"\"")
	cookie := http.Cookie{
		Name:       "isTrustedVolunteer",
		Value:      "",
		Path:       "/",
		Domain:     "localhost",
		Expires:    time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		Raw:        "",
		Unparsed:   []string{unparsedAttributePairs.String()}, // Raw text of unparsed attribute-value pairs
	}
	return cookie
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	cookie := createEmptyCookie()
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", 307)
}

// Following pattern shown at https://github.com/gorilla/mux#static-files
type uiHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP serves static files
func (h uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// AnimalControlReport returns the animal control report
func (h *handler) AnimalControlReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	animalControlReport, err := h.Service.GetAnimalControlReport()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// Return volunteers as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animalControlReport)
}

// AnimalControlReportDrillDownOne returns the animal control report
func (h *handler) AnimalControlReportDrillDownOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	animalControlReportDrillDownOne, err := h.Service.GetAnimalControlReportDrillDownOne(startDate, endDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// Return volunteers as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animalControlReportDrillDownOne)
}

// AnimalControlReportDrillDownTwo returns the animal control report
func (h *handler) AnimalControlReportDrillDownTwo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	animalControlReportDrillDownTwo, err := h.Service.GetAnimalControlReportDrillDownTwo(startDate, endDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error getting records from database: ", err)))
		return
	}

	// Return volunteers as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animalControlReportDrillDownTwo)
}
