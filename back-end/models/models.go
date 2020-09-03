package models

import (
	"database/sql"
)

type Dog struct {
	ID                          int            `json:"id"`
	Name                        string         `json:"name"`
	Breed                       []string       `json:"breed"`
	AlterationStatus            bool           `json:"alterationStatus"`
	Description                 string         `json:"description"`
	DateOfBirth                 string         `json:"dateOfBirth"`
	Sex                         string         `json:"sex"`
	MicrochipID                 sql.NullString `json:"microchipId"`
	SurrenderDate               string         `json:"surrenderDate"`
	SurrenderReason             string         `json:"surrenderReason"`
	SurrenderWasByAnimalControl bool           `json:"surrenderWasByAnimalControl"`
	VolunteerID                 int            `json:"volunteerId"`
	Expenses                    []Expense      `json:"expenses"`
}

type Volunteer struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Cell      string `json:"cell"`
	Password  string `json:"password"`
	StartDate string `json:"startDate"`
	IsTrusted bool   `json:"isTrusted"`
}

type MonthlyAdoptionReportItem struct {
	Breed          string          `json:"breed"`
	SurrenderCount sql.NullFloat64 `json:surrenderCount`
	AdoptionCount  sql.NullFloat64 `json:adoptionCount`
	Expenses       sql.NullFloat64 `json:expenses`
	AdoptionFees   sql.NullFloat64 `json:adoptionfees`
	Profit         sql.NullFloat64 `json:profit`
}

type ExpenseAnalysisItem struct {
	Vendor        string `json:"vendor"`
	TotalSpending int    `json:"totalSpending"`
}

type Applicant struct {
	ID                   int    `json:"id,omitempty"`
	Zip                  string `json:"zip,omitempty"`
	State                string `json:"state,omitempty"`
	City                 string `json:"city,omitempty"`
	Street               string `json:"street,omitempty"`
	PhoneNumber          string `json:"phoneNumber,omitempty"`
	Email                string `json:"email,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	CoApplicantLastName  string `json:"coApplicantLastName,omitempty"`
	CoApplicantFirstName string `json:"coApplicantFirstName,omitempty"`
}

type Application struct {
	ID                   string    `json:"id"`
	Date                 string    `json:"date"`
	State                string    `json:"state"`
	CoApplicantLastName  string    `json:"coApplicantLastName"`
	CoApplicantFirstName string    `json:"coApplicantFirstName"`
	ApplicantIdFk        int       `json:"applicantIdFk"`
	Applicant            Applicant `json:"applicant,omitempty"`
}

type Breed struct {
	ID    int    `json:"id"`
	Breed string `json:"breed"`
}

type DogBreed struct {
	DogIdFk   int `json:"dogIdFk"`
	BreedIdFk int `json:"breedIdFk"`
}
type Expense struct {
	Date          string `json:"date"`
	Vendor        string `json:"vendor"`
	Description   string `json:"description"`
	AmountInCents int    `json:"amountInCents"`
	DogIdFk       int    `json:"dogIdFk"`
}
type Adoption struct {
	DateAdopted         string `json:"dateAdopted"`
	ApplicationNumberFk int    `json:"applicationNumberFk"`
	DogIdFk             int    `json:"dogIdFk"`
}
type AnimalControlReport struct {
	Month              string  `json:"month"`
	Year               int     `json:"year"`
	DogsTotalCount     int     `json:"dogsTotalCount"`
	DogsSixtyDaysCount int     `json:"dogsSixtyDaysCount"`
	Expenses           float32 `json:"expenses"`
}
type AnimalControlReportDrillDown struct {
	DogID            int    `json:"dogId"`
	Sex              string `json:"sex"`
	AlterationStatus bool   `json:"alterationStatus"`
	MicrochipID      string `json:"microchipId"`
	SurrenderDate    string `json:"surrenderDate"`
	AdoptionDate     string `json:"adoptionDate"`
	Breed            string `json:"breed"`
}
