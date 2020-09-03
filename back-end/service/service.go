package service

import (
	"back-end/database"
	"back-end/models"
)

type Service struct {
	Postgres database.Postgres
}

// CreateVolunteer creates a new volunteer
func (s *Service) CreateVolunteer(volunteer models.Volunteer) (models.Volunteer, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	// Save volunteer in database
	return s.Postgres.SaveVolunteer(volunteer)
}

// GetVolunteers returns a list of all of the volunteers
func (s *Service) GetVolunteers() ([]models.Volunteer, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	// Get volunteers from database
	return s.Postgres.GetVolunteers()
}

// ExpenseAnalysis returns the expense analysis
func (s *Service) ExpenseAnalysis() ([]models.ExpenseAnalysisItem, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.ExpenseAnalysis()
}

// GetVolunteersLike
func (s *Service) GetVolunteersLike(like string) ([]models.Volunteer, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	// Get volunteers from database
	return s.Postgres.GetVolunteersLike(like)
}

// GetVolunteerByEmail returns volunteer with email address, if volunteer exists
func (s *Service) GetVolunteerByEmail(email string) (models.Volunteer, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	// Get volunteers from database
	return s.Postgres.GetVolunteerByEmail(email)
}

// GetApplicants
func (s *Service) GetApplicants() ([]models.Applicant, error) {
	return s.Postgres.GetApplicants()
}

// GetApplicantByEmail
func (s *Service) GetApplicantByEmail(email string) (models.Applicant, error) {
	return s.Postgres.GetApplicantByEmail(email)
}

// GetApplicantByID
func (s *Service) GetApplicantByID(applicantID int) (models.Applicant, error) {
	return s.Postgres.GetApplicantByID(applicantID)
}

// GetApprovedApplicants
func (s *Service) GetApprovedApplicants() ([]models.Applicant, error) {
	return s.Postgres.GetApprovedApplicants()
}

// GetApprovedApplicantsLike
func (s *Service) GetApprovedApplicantsLike(lastNameFragment string) ([]models.Applicant, error) {
	return s.Postgres.GetApprovedApplicantsLike(lastNameFragment)
}

// GetDogs
func (s *Service) GetDogs() ([]models.Dog, error) {
	return s.Postgres.GetDogs()
}

// GetDogs
func (s *Service) GetCurrentDogs() ([]models.Dog, error) {
	return s.Postgres.GetCurrentDogs()
}

// GetDog
func (s *Service) GetDog(id int) (models.Dog, error) {
	return s.Postgres.GetDog(id)
}

// UpdateDog
func (s *Service) UpdateDog(dog models.Dog) error {
	return s.Postgres.UpdateDog(dog)
}

// GetApprovedApplications
func (s *Service) GetApprovedApplications() ([]models.Application, error) {
	return s.Postgres.GetApprovedApplications()
}

// GetPendingApplications
func (s *Service) GetPendingApplications() ([]models.Application, error) {
	return s.Postgres.GetPendingApplications()
}

// GetLatestApprovedApplication
func (s *Service) GetLatestApprovedApplication(applicantID int, coApplicantFirstName string, coApplicantLastName string) ([]models.Application, error) {
	return s.Postgres.GetLatestApprovedApplication(applicantID, coApplicantFirstName, coApplicantLastName)
}

// ChangeApplicationStatus
func (s *Service) ChangeApplicationStatus(id string, approve bool) error {
	return s.Postgres.ChangeApplicationStatus(id, approve)
}

// GetBreeds
func (s *Service) GetBreeds() ([]models.Breed, error) {
	return s.Postgres.GetBreeds()
}

// GetBreeds
func (s *Service) GetBreedID(breed string) (int, error) {
	return s.Postgres.GetBreedID(breed)
}

// GetDogBreeds
func (s *Service) GetDogBreeds() ([]models.DogBreed, error) {
	return s.Postgres.GetDogBreeds()
}

// GetBreedsByDogID
func (s *Service) GetBreedsByDogID(dogID int) ([]string, error) {
	return s.Postgres.GetBreedsByDogID(dogID)
}

// DeleteDogBreedByDogID
func (s *Service) DeleteDogBreedByDogID(dogID int) error {
	return s.Postgres.DeleteDogBreedByDogID(dogID)
}

// GetExpenses
func (s *Service) GetExpenses() ([]models.Expense, error) {
	return s.Postgres.GetExpenses()
}

// GetExpensesByDogID
func (s *Service) GetExpensesByDogID(id int) ([]models.Expense, error) {
	return s.Postgres.GetExpensesByDogID(id)
}

// GetAdoptions
func (s *Service) GetAdoptions() ([]models.Adoption, error) {
	return s.Postgres.GetAdoptions()
}

// CreateApplicant
func (s *Service) CreateApplicant(applicant models.Applicant) (models.Applicant, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveApplicant(applicant)
}

// CreateDog
func (s *Service) CreateDog(dog models.Dog) (models.Dog, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveDog(dog)
}

// CreateApplication
func (s *Service) CreateApplication(application models.Application) (models.Application, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveApplication(application)
}

// CreateBreed
func (s *Service) CreateBreed(breed models.Breed) (models.Breed, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveBreed(breed)
}

// CreateDogBreed
func (s *Service) CreateDogBreed(dogBreed models.DogBreed) (models.DogBreed, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveDogBreed(dogBreed)
}

// CreateExpense
func (s *Service) CreateExpense(expense models.Expense) (models.Expense, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveExpense(expense)
}

// CreateAdoption
func (s *Service) CreateAdoption(adoption models.Adoption) (models.Adoption, error) {

	//
	// Perform any business logic or data manipulation/validation here
	//

	return s.Postgres.SaveAdoption(adoption)
}

// GetAnimalControlReport returns the animal control report
func (s *Service) GetAnimalControlReport() ([]models.AnimalControlReport, error) {
	//
	// Perform any business logic or data manipulation/validation here
	//

	// Get volunteers from database
	return s.Postgres.GetAnimalControlReport()
}

// Monthly adoption report
func (s *Service) MonthlyAdoptionReport(startDate string, endDate string) ([]models.MonthlyAdoptionReportItem, error) {

	return s.Postgres.MonthlyAdoptionReport(startDate, endDate)
}

// GetAnimalControlReportDrillDownOne returns the animal control report
func (s *Service) GetAnimalControlReportDrillDownOne(startDate string, endDate string) ([]models.AnimalControlReportDrillDown, error) {
	return s.Postgres.GetAnimalControlReportDrillDownOne(startDate, endDate)
}

// GetAnimalControlReportDrillDownTwo returns the animal control report
func (s *Service) GetAnimalControlReportDrillDownTwo(startDate string, endDate string) ([]models.AnimalControlReportDrillDown, error) {
	return s.Postgres.GetAnimalControlReportDrillDownTwo(startDate, endDate)
}
