package database

import (
	"back-end/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Postgres holds the information for the database conneciton
type Postgres struct {
	db      *sql.DB
	timeout int
}

// Connect establishes a connection to Postgres
func Connect(host string, port int, user, password, dbname string, timeout int) (Postgres, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{db: db, timeout: timeout}, nil
}

// Close closes the Postgres connection
func (p *Postgres) Close() {
	p.db.Close()
}

// SaveVolunteer adds a new volunteer to the database
func (p *Postgres) SaveVolunteer(volunteer models.Volunteer) (models.Volunteer, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Volunteer (email, first_name, last_name, cell_phone_number, password, start_date, is_trusted_volunteer)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING volunteer_id, email, first_name, last_name, cell_phone_number, start_date, is_trusted_volunteer`
	row := p.db.QueryRowContext(context, sqlStatement, volunteer.Email, volunteer.FirstName, volunteer.LastName,
		volunteer.Cell, volunteer.Password, volunteer.StartDate, volunteer.IsTrusted)

	var insertedVolunteer models.Volunteer
	err := row.Scan(&insertedVolunteer.ID, &insertedVolunteer.Email, &insertedVolunteer.FirstName, &insertedVolunteer.LastName,
		&insertedVolunteer.Cell, &insertedVolunteer.StartDate, &insertedVolunteer.IsTrusted)
	if err != nil {
		return models.Volunteer{}, err
	}

	return insertedVolunteer, nil
}

func (p *Postgres) MonthlyAdoptionReport(startDate string, endDate string) ([]models.MonthlyAdoptionReportItem, error) {
	//get expenses,
	sqlStatement := fmt.Sprintf(`

select 
string_agg(distinct breed, '/' order by breed asc) as breedName,
sum(surrendered) as surrendercount,
sum(adopted) as adoptioncount,
sum(amount_in_cents) as expenses,
sum(adoptionfees) as adoptionfees,
sum(profit) as profit

from
(select *, 
(0.15*adoptionfees) as profit
from
(select *,
(case when(surrender_was_by_animal_control = true and adopted =1) then
 0.15*amount_in_cents
 
 when (surrender_was_by_animal_control = false and adopted =1) then
 1.15*amount_in_cents
end)
 as adoptionFees
 from
(select dog_id, surrender_was_by_animal_control ,
 (case when (surrender_date >= date '%s' AND surrender_date <= date '%s' ) then

 						   1
 							when (surrender_date <  date '%s' AND surrender_date >  date '%s') then
 							0
						   end) as surrendered,
 						   (case when (date_adopted >= date '%s' AND date_adopted <=  date '%s') then

 						   1
 							when (date_adopted <  date '%s' AND date_adopted > date '%s') then
 							0
						   end) as adopted,breed,
 						   
 						amount_in_cents
							from dog
 							left join (select date_adopted, dog_id_fk from adoption) as adoption on adoption.dog_id_fk = dog.dog_id
 							left join (select dog_id_fk, amount_in_cents from expense) as expense on expense.dog_id_fk = dog.dog_id
 							join dogbreed on dogbreed.dog_id_fk = dog.dog_id
  							join breed on dogbreed.breed_id_fk = breed.breed_id
  where (surrender_date >=  date '%s' AND surrender_date <=  date '%s') OR (date_adopted >=  date '%s' AND date_adopted <=  date '%s') 
) as nest1
 ) as nest2
	) as nest3		
			group by dog_id	
			order by breedname asc`, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate)

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	monthlyAdoptionItems := make([]models.MonthlyAdoptionReportItem, 0)

	//expenseAnalysisItems := make([]models.ExpenseAnalysisItem, 0)
	defer rows.Close()

	for rows.Next() {
		var monthlyAdoptionItem models.MonthlyAdoptionReportItem
		if err := rows.Scan(&monthlyAdoptionItem.Breed, &monthlyAdoptionItem.SurrenderCount, &monthlyAdoptionItem.AdoptionCount, &monthlyAdoptionItem.Expenses, &monthlyAdoptionItem.AdoptionFees, &monthlyAdoptionItem.Profit); err != nil {
			return nil, err
		}
		monthlyAdoptionItems = append(monthlyAdoptionItems, monthlyAdoptionItem)
	}

	//for rows.Next() {
	//	var expenseAnalysisItem models.ExpenseAnalysisItem
	//	if err := rows.Scan(&expenseAnalysisItem.TotalSpending, &expenseAnalysisItem.Vendor); err != nil {
	//		return nil, err
	//	}
	//	expenseAnalysisItems = append(expenseAnalysisItems, expenseAnalysisItem)
	//}

	//return expenseAnalysisItems, nil
	return monthlyAdoptionItems, nil
}

func (p *Postgres) ExpenseAnalysis() ([]models.ExpenseAnalysisItem, error) {

	sqlStatement := `
	SELECT SUM(amount_in_cents) AS total_amount, vendor 
	FROM expense GROUP BY vendor ORDER BY total_amount DESC`
	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	expenseAnalysisItems := make([]models.ExpenseAnalysisItem, 0)

	defer rows.Close()
	for rows.Next() {
		var expenseAnalysisItem models.ExpenseAnalysisItem
		if err := rows.Scan(&expenseAnalysisItem.TotalSpending, &expenseAnalysisItem.Vendor); err != nil {
			return nil, err
		}
		expenseAnalysisItems = append(expenseAnalysisItems, expenseAnalysisItem)
	}

	return expenseAnalysisItems, nil
}

// GetVolunteers gets all of the volunteers from the database
func (p *Postgres) GetVolunteers() ([]models.Volunteer, error) {

	sqlStatement := `
		SELECT volunteer_id, email, first_name, last_name, cell_phone_number,
		start_date, is_trusted_volunteer FROM Volunteer`
	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	volunteers := make([]models.Volunteer, 0)

	defer rows.Close()
	for rows.Next() {
		var volunteer models.Volunteer
		if err := rows.Scan(&volunteer.ID, &volunteer.Email, &volunteer.FirstName, &volunteer.LastName, &volunteer.Cell,
			&volunteer.StartDate, &volunteer.IsTrusted); err != nil {
			return nil, err
		}
		volunteers = append(volunteers, volunteer)
	}

	return volunteers, nil
}

// GetVolunteers gets all of the volunteers from the database
func (p *Postgres) GetVolunteersLike(like string) ([]models.Volunteer, error) {

	sqlStatement := `
		SELECT volunteer_id, email, first_name, last_name, cell_phone_number,
		start_date, is_trusted_volunteer FROM Volunteer
		WHERE LOWER(first_name) LIKE '%' || $1 || '%' OR LOWER(last_name) LIKE '%' || $1 || '%'
		OR CONCAT(LOWER(first_name), ' ',LOWER(last_name)) LIKE '%' || $1 || '%'`
	rows, err := p.db.Query(sqlStatement, strings.ToLower(like))
	if err != nil {
		return nil, err
	}

	volunteers := make([]models.Volunteer, 0)

	defer rows.Close()
	for rows.Next() {
		var volunteer models.Volunteer
		if err := rows.Scan(&volunteer.ID, &volunteer.Email, &volunteer.FirstName, &volunteer.LastName, &volunteer.Cell,
			&volunteer.StartDate, &volunteer.IsTrusted); err != nil {
			return nil, err
		}
		volunteers = append(volunteers, volunteer)
	}

	return volunteers, nil
}

// GetVolunteerByEmail gets volunteer with email specified from the database, if volunteer exists
func (p *Postgres) GetVolunteerByEmail(email string) (models.Volunteer, error) {

	sqlStatement := `
		SELECT volunteer_id, is_trusted_volunteer, password FROM Volunteer WHERE email=$1`
	rows, err := p.db.Query(sqlStatement, email)
	if err != nil {
		return models.Volunteer{}, err
	}

	volunteer := models.Volunteer{}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&volunteer.ID, &volunteer.IsTrusted, &volunteer.Password); err != nil {
			return models.Volunteer{}, err
		}
	}

	return volunteer, nil
}

// GetApplicants
func (p *Postgres) GetApplicants() ([]models.Applicant, error) {
	sqlStatement := `
	SELECT applicant_id, zip, state, city, street, phone_number, email, last_name, first_name
	FROM Applicant`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	applicants := make([]models.Applicant, 0)

	defer rows.Close()
	for rows.Next() {
		var applicant models.Applicant
		if err := rows.Scan(&applicant.ID,
			&applicant.Zip,
			&applicant.State,
			&applicant.City,
			&applicant.Street,
			&applicant.PhoneNumber,
			&applicant.Email,
			&applicant.LastName,
			&applicant.FirstName); err != nil {
			return nil, err
		}
		applicants = append(applicants, applicant)
	}
	return applicants, nil
}

// GetApprovedApplicants gets all of the applicants from approved
// applications from the database
func (p *Postgres) GetApprovedApplicants() ([]models.Applicant, error) {
	sqlStatement := `
		SELECT DISTINCT applicant_id, zip, Applicant.state, city, street, phone_number, email, last_name, first_name, co_applicant_last_name, co_applicant_first_name
		FROM Applicant 
		INNER JOIN Application ON applicant_id = applicant_id_fk
		WHERE Application.state = 'approved'
		AND application_number NOT IN (SELECT application_number_fk 
									   AS application_number 
									   FROM Adoption)`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	approvedApplicants := make([]models.Applicant, 0)

	defer rows.Close()
	for rows.Next() {
		var applicant models.Applicant
		if err := rows.Scan(&applicant.ID,
			&applicant.Zip,
			&applicant.State,
			&applicant.City,
			&applicant.Street,
			&applicant.PhoneNumber,
			&applicant.Email,
			&applicant.LastName,
			&applicant.FirstName,
			&applicant.CoApplicantLastName,
			&applicant.CoApplicantFirstName); err != nil {
			return nil, err
		}
		approvedApplicants = append(approvedApplicants, applicant)
	}

	return approvedApplicants, nil
}

// GetApprovedApplicantsLike gets all of the applicants from approved
// applications whose last name contains the fragment provided from the database
func (p *Postgres) GetApprovedApplicantsLike(like string) ([]models.Applicant, error) {

	sqlStatement := `
		SELECT DISTINCT on (applicant_id) applicant_id, zip, Applicant.state, city, street, phone_number, email, last_name, first_name, co_applicant_last_name, co_applicant_first_name
		FROM Applicant 
		INNER JOIN Application ON applicant_id = applicant_id_fk
		WHERE Application.state = 'approved'
		AND LOWER(last_name) LIKE '%' || $1 || '%'
		AND application_number NOT IN (SELECT application_number_fk
									   AS application_number
									   FROM Adoption)
		UNION DISTINCT
		SELECT DISTINCT applicant_id, zip, Applicant.state, city, street, phone_number, email, last_name, first_name, co_applicant_last_name, co_applicant_first_name
		FROM Applicant 
		INNER JOIN Application ON applicant_id = applicant_id_fk
		WHERE Application.state = 'approved'
		AND LOWER(co_applicant_last_name) LIKE '%' || $1 || '%'
		AND application_number NOT IN (SELECT application_number_fk
									   AS application_number
									   FROM Adoption)`

	rows, err := p.db.Query(sqlStatement, strings.ToLower(like))
	if err != nil {
		return nil, err
	}

	approvedApplicants := make([]models.Applicant, 0)

	defer rows.Close()
	for rows.Next() {
		var applicant models.Applicant
		if err := rows.Scan(&applicant.ID,
			&applicant.Zip,
			&applicant.State,
			&applicant.City,
			&applicant.Street,
			&applicant.PhoneNumber,
			&applicant.Email,
			&applicant.LastName,
			&applicant.FirstName,
			&applicant.CoApplicantLastName,
			&applicant.CoApplicantFirstName); err != nil {
			return nil, err
		}
		approvedApplicants = append(approvedApplicants, applicant)
	}

	return approvedApplicants, nil
}

// GetApplicantByEmail returns applicant by email
func (p *Postgres) GetApplicantByEmail(email string) (models.Applicant, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
	SELECT applicant_id, zip, state, city, street, phone_number, email, last_name, first_name
	FROM Applicant WHERE email=$1`

	rows, err := p.db.QueryContext(context, sqlStatement, email)
	if err != nil {
		return models.Applicant{}, err
	}

	var applicant models.Applicant

	defer rows.Close()
	for rows.Next() {

		if err := rows.Scan(&applicant.ID,
			&applicant.Zip,
			&applicant.State,
			&applicant.City,
			&applicant.Street,
			&applicant.PhoneNumber,
			&applicant.Email,
			&applicant.LastName,
			&applicant.FirstName); err != nil {
			return models.Applicant{}, err
		}
	}
	return applicant, nil
}

// GetApplicantByID returns applicant by ID
func (p *Postgres) GetApplicantByID(applicantID int) (models.Applicant, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
	SELECT email, zip, state, city, street, phone_number, email, last_name, first_name
	FROM Applicant WHERE applicant_id=$1`

	rows, err := p.db.QueryContext(context, sqlStatement, applicantID)
	if err != nil {
		return models.Applicant{}, err
	}

	var applicant models.Applicant

	defer rows.Close()
	for rows.Next() {

		if err := rows.Scan(&applicant.Email,
			&applicant.Zip,
			&applicant.State,
			&applicant.City,
			&applicant.Street,
			&applicant.PhoneNumber,
			&applicant.Email,
			&applicant.LastName,
			&applicant.FirstName); err != nil {
			return models.Applicant{}, err
		}
	}
	applicant.ID = applicantID
	return applicant, nil
}

// GetDogs
func (p *Postgres) GetDogs() ([]models.Dog, error) {
	sqlStatement := `
	SELECT dog_id, name, alteration_status, description, date_of_birth, sex, microchip_id, surrender_date, surrender_reason, surrender_was_by_animal_control, volunteer_id_fk
	FROM Dog ORDER BY surrender_date ASC`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Dogs := make([]models.Dog, 0)

	defer rows.Close()
	for rows.Next() {
		var Dog models.Dog
		if err := rows.Scan(
			&Dog.ID,
			&Dog.Name,
			&Dog.AlterationStatus,
			&Dog.Description,
			&Dog.DateOfBirth,
			&Dog.Sex,
			&Dog.MicrochipID,
			&Dog.SurrenderDate,
			&Dog.SurrenderReason,
			&Dog.SurrenderWasByAnimalControl,
			&Dog.VolunteerID); err != nil {
			return nil, err
		}
		Dogs = append(Dogs, Dog)
	}
	return Dogs, nil
}

// GetDogs
func (p *Postgres) GetCurrentDogs() ([]models.Dog, error) {
	sqlStatement := `
	SELECT dog_id, name, alteration_status, description, date_of_birth, sex, microchip_id, surrender_date, surrender_reason, surrender_was_by_animal_control, volunteer_id_fk
	FROM Dog 
	WHERE dog_id NOT IN (SELECT dog_id_fk AS dog_id FROM Adoption)`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Dogs := make([]models.Dog, 0)

	defer rows.Close()
	for rows.Next() {
		var Dog models.Dog
		if err := rows.Scan(
			&Dog.ID,
			&Dog.Name,
			&Dog.AlterationStatus,
			&Dog.Description,
			&Dog.DateOfBirth,
			&Dog.Sex,
			&Dog.MicrochipID,
			&Dog.SurrenderDate,
			&Dog.SurrenderReason,
			&Dog.SurrenderWasByAnimalControl,
			&Dog.VolunteerID); err != nil {
			return nil, err
		}
		Dogs = append(Dogs, Dog)
	}
	return Dogs, nil
}

// GetDog
func (p *Postgres) GetDog(id int) (models.Dog, error) {
	sqlStatement := `
	SELECT name, alteration_status, description ,date_of_birth, sex, microchip_id, surrender_date, surrender_reason, surrender_was_by_animal_control, volunteer_id_fk
	FROM Dog WHERE dog_id = $1`

	var Dog models.Dog

	rows, err := p.db.Query(sqlStatement, id)
	if err != nil {
		return Dog, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&Dog.Name,
			&Dog.AlterationStatus,
			&Dog.Description,
			&Dog.DateOfBirth,
			&Dog.Sex,
			&Dog.MicrochipID,
			&Dog.SurrenderDate,
			&Dog.SurrenderReason,
			&Dog.SurrenderWasByAnimalControl,
			&Dog.VolunteerID); err != nil {
			return Dog, err
		}
	}

	Dog.ID = id

	return Dog, nil
}

// GetApprovedApplications
func (p *Postgres) GetApprovedApplications() ([]models.Application, error) {

	sqlStatement := `
	SELECT DISTICNT application_number, date, state, co_applicant_last_name, co_applicant_first_name, applicant_id_fk
	FROM Application
	WHERE state = 'approved'
	ORDER BY date DESC`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Applications := make([]models.Application, 0)

	defer rows.Close()
	for rows.Next() {
		var Application models.Application
		if err := rows.Scan(
			&Application.ID,
			&Application.Date,
			&Application.State,
			&Application.CoApplicantLastName,
			&Application.CoApplicantFirstName,
			&Application.ApplicantIdFk); err != nil {
			return nil, err
		}
		Applications = append(Applications, Application)
	}
	return Applications, nil
}

// GetApprovedApplication
func (p *Postgres) GetLatestApprovedApplication(applicantID int, coApplicantFirstName string, coApplicantLastName string) ([]models.Application, error) {

	var sqlStatement string
	var application models.Application
	Applications := make([]models.Application, 0)
	var rows *sql.Rows
	var err error

	if len(coApplicantLastName) > 0 {
		sqlStatement = `
			SELECT application_number, date
			FROM Application
			WHERE state = 'approved'
			AND applicant_id_fk = $1 
			AND co_applicant_first_name = $2 
			AND co_applicant_last_name = $3
			AND application_number NOT IN (SELECT application_number_fk AS application_number FROM Adoption)
			ORDER BY date DESC
			LIMIT 1`
		rows, err = p.db.Query(sqlStatement, applicantID, coApplicantFirstName, coApplicantLastName)
	} else {
		sqlStatement = `
			SELECT application_number, date
			FROM Application
			WHERE state = 'approved'
			AND applicant_id_fk = $1 
			AND application_number NOT IN (SELECT application_number_fk AS application_number FROM Adoption)
			ORDER BY date DESC
			LIMIT 1`
		rows, err = p.db.Query(sqlStatement, applicantID)
	}

	if err != nil {
		return Applications, err
	}

	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&application.ID,
			&application.Date); err != nil {
			return Applications, err
		}

		application.State = "approved"
		application.ApplicantIdFk = applicantID
		application.CoApplicantFirstName = coApplicantFirstName
		application.CoApplicantLastName = coApplicantLastName
		Applications = append(Applications, application)
	}
	return Applications, nil
}

// GetPendingApplications
func (p *Postgres) GetPendingApplications() ([]models.Application, error) {

	sqlStatement := `
	SELECT application_number, date, state, co_applicant_last_name, co_applicant_first_name, applicant_id_fk
	FROM Application WHERE state = 'pending approval'`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Applications := make([]models.Application, 0)

	defer rows.Close()
	for rows.Next() {
		var Application models.Application
		if err := rows.Scan(
			&Application.ID,
			&Application.Date,
			&Application.State,
			&Application.CoApplicantLastName,
			&Application.CoApplicantFirstName,
			&Application.ApplicantIdFk); err != nil {
			return nil, err
		}
		Applications = append(Applications, Application)
	}
	return Applications, nil
}

// GetBreeds
func (p *Postgres) GetBreeds() ([]models.Breed, error) {
	sqlStatement := `
	SELECT breed_id, breed
	FROM Breed`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Breeds := make([]models.Breed, 0)

	defer rows.Close()
	for rows.Next() {
		var Breed models.Breed
		if err := rows.Scan(
			&Breed.ID,
			&Breed.Breed); err != nil {
			return nil, err
		}
		Breeds = append(Breeds, Breed)
	}
	return Breeds, nil
}

// GetDogBreeds
func (p *Postgres) GetDogBreeds() ([]models.DogBreed, error) {
	sqlStatement := `
	SELECT dog_id_fk, breed_id_fk
	FROM DogBreed`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	DogBreeds := make([]models.DogBreed, 0)

	defer rows.Close()
	for rows.Next() {
		var DogBreed models.DogBreed
		if err := rows.Scan(
			&DogBreed.DogIdFk,
			&DogBreed.BreedIdFk); err != nil {
			return nil, err
		}
		DogBreeds = append(DogBreeds, DogBreed)
	}
	return DogBreeds, nil
}

// GetBreedsByDogID
func (p *Postgres) GetBreedsByDogID(dogID int) ([]string, error) {
	sqlStatement := `
	SELECT breed
	FROM DogBreed INNER JOIN Breed ON breed_id = breed_id_fk 
	WHERE dog_id_fk = $1`

	rows, err := p.db.Query(sqlStatement, dogID)
	if err != nil {
		return nil, err
	}
	Breeds := make([]string, 0)

	defer rows.Close()
	for rows.Next() {
		var Breed models.Breed
		if err := rows.Scan(&Breed.Breed); err != nil {
			return nil, err
		}
		Breeds = append(Breeds, Breed.Breed)
	}
	return Breeds, nil
}

// DeleteDogBreedByDogID
func (p *Postgres) DeleteDogBreedByDogID(dogID int) error {
	sqlStatement := `
	DELETE
	FROM DogBreed 
	WHERE dog_id_fk = $1`

	rows, err := p.db.Query(sqlStatement, dogID)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

// GetBreedID ...
func (p *Postgres) GetBreedID(breed string) (int, error) {
	sqlStatement := `
	SELECT breed_id
	FROM Breed WHERE breed = $1`

	rows, err := p.db.Query(sqlStatement, breed)
	if err != nil {
		return 0, err
	}

	var breedID int

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&breedID); err != nil {
			return 0, err
		}
	}
	return breedID, nil
}

// GetExpenses
func (p *Postgres) GetExpenses() ([]models.Expense, error) {
	sqlStatement := `
	SELECT date, vendor, description, amount_in_cents, dog_id_fk
	FROM Expense`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Expenses := make([]models.Expense, 0)

	defer rows.Close()
	for rows.Next() {
		var Expense models.Expense
		if err := rows.Scan(
			&Expense.Date,
			&Expense.Vendor,
			&Expense.Description,
			&Expense.AmountInCents,
			&Expense.DogIdFk); err != nil {
			return nil, err
		}
		Expenses = append(Expenses, Expense)
	}
	return Expenses, nil
}

// GetExpensesByDogID
func (p *Postgres) GetExpensesByDogID(id int) ([]models.Expense, error) {
	sqlStatement := `
	SELECT date, vendor, description, amount_in_cents
	FROM Expense WHERE dog_id_fk = $1`

	rows, err := p.db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	Expenses := make([]models.Expense, 0)

	defer rows.Close()
	for rows.Next() {
		var Expense models.Expense
		if err := rows.Scan(
			&Expense.Date,
			&Expense.Vendor,
			&Expense.Description,
			&Expense.AmountInCents); err != nil {
			return nil, err
		}
		Expenses = append(Expenses, Expense)
	}
	return Expenses, nil
}

// GetAdoptions
func (p *Postgres) GetAdoptions() ([]models.Adoption, error) {
	sqlStatement := `
	SELECT date_adopted, application_number_fk, dog_id_fk
	FROM Adoption`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	Adoptions := make([]models.Adoption, 0)

	defer rows.Close()
	for rows.Next() {
		var Adoption models.Adoption
		if err := rows.Scan(
			&Adoption.DateAdopted,
			&Adoption.ApplicationNumberFk,
			&Adoption.DogIdFk); err != nil {
			return nil, err
		}
		Adoptions = append(Adoptions, Adoption)
	}
	return Adoptions, nil
}

// ChangeApplicationStatus
func (p *Postgres) ChangeApplicationStatus(id string, approve bool) error {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var state string
	if approve {
		state = "approved"
	} else {
		state = "rejected"
	}

	sqlStatement := `
		UPDATE Application 
		SET state = $2
		WHERE application_number = $1`

	if _, err := p.db.ExecContext(context, sqlStatement, id, state); err != nil {
		return err
	}

	return nil
}

// SaveApplicant
func (p *Postgres) SaveApplicant(applicant models.Applicant) (models.Applicant, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Applicant (zip, state, city, street, phone_number, email, last_name, first_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING applicant_id, zip, state, city, street, phone_number, email, last_name, first_name`
	row := p.db.QueryRowContext(context, sqlStatement, applicant.Zip, applicant.State, applicant.City, applicant.Street, applicant.PhoneNumber, applicant.Email, applicant.LastName, applicant.FirstName)

	var insertedApplicant models.Applicant
	err := row.Scan(&insertedApplicant.ID,
		&insertedApplicant.Zip, &insertedApplicant.State, &insertedApplicant.City, &insertedApplicant.Street, &insertedApplicant.PhoneNumber, &insertedApplicant.Email, &insertedApplicant.LastName, &insertedApplicant.FirstName)
	if err != nil {
		return models.Applicant{}, err
	}

	return insertedApplicant, nil
}

// SaveDog
func (p *Postgres) SaveDog(dog models.Dog) (models.Dog, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Dog (name, alteration_status, description, date_of_birth, sex, microchip_id, surrender_date, surrender_reason, surrender_was_by_animal_control, volunteer_id_fk)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING dog_id, name, alteration_status, description ,date_of_birth, sex, microchip_id, surrender_date, surrender_reason, surrender_was_by_animal_control, volunteer_id_fk`

	var row *sql.Row
	var insertedDog models.Dog

	if dog.MicrochipID.String == "" {
		row = p.db.QueryRowContext(context, sqlStatement,
			dog.Name,
			dog.AlterationStatus,
			dog.Description,
			dog.DateOfBirth,
			dog.Sex,
			nil,
			dog.SurrenderDate,
			dog.SurrenderReason,
			dog.SurrenderWasByAnimalControl,
			dog.VolunteerID)
	} else {
		row = p.db.QueryRowContext(context, sqlStatement,
			dog.Name,
			dog.AlterationStatus,
			dog.Description,
			dog.DateOfBirth,
			dog.Sex,
			dog.MicrochipID,
			dog.SurrenderDate,
			dog.SurrenderReason,
			dog.SurrenderWasByAnimalControl,
			dog.VolunteerID)
	}

	err := row.Scan(&insertedDog.ID,
		&insertedDog.Name,
		&insertedDog.AlterationStatus,
		&insertedDog.Description,
		&insertedDog.DateOfBirth,
		&insertedDog.Sex,
		&insertedDog.MicrochipID,
		&insertedDog.SurrenderDate,
		&insertedDog.SurrenderReason,
		&insertedDog.SurrenderWasByAnimalControl,
		&insertedDog.VolunteerID)
	if err != nil {
		return models.Dog{}, err
	}

	return insertedDog, nil
}

// UpdateDog
func (p *Postgres) UpdateDog(dog models.Dog) error {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		Update Dog 
		SET alteration_status = $1, sex = $2, microchip_id = $3
		WHERE dog_id = $4`

	if dog.MicrochipID.String == "" {
		if _, err := p.db.ExecContext(context,
			sqlStatement,
			dog.AlterationStatus,
			dog.Sex,
			nil,
			dog.ID); err != nil {
			return err
		}
	} else {
		if _, err := p.db.ExecContext(context,
			sqlStatement,
			dog.AlterationStatus,
			dog.Sex,
			dog.MicrochipID.String,
			dog.ID); err != nil {
			return err
		}
	}
	return nil
}

// SaveApplication
func (p *Postgres) SaveApplication(application models.Application) (models.Application, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Application (date, state, co_applicant_last_name, co_applicant_first_name, applicant_id_fk)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING application_number, date, state, co_applicant_last_name, co_applicant_first_name, applicant_id_fk`
	row := p.db.QueryRowContext(context, sqlStatement,
		application.Date,
		application.State,
		application.CoApplicantLastName,
		application.CoApplicantFirstName,
		application.ApplicantIdFk)

	var insertedApplication models.Application
	err := row.Scan(
		&insertedApplication.ID,
		&insertedApplication.Date,
		&insertedApplication.State,
		&insertedApplication.CoApplicantLastName,
		&insertedApplication.CoApplicantFirstName,
		&insertedApplication.ApplicantIdFk)
	if err != nil {
		return models.Application{}, err
	}

	return insertedApplication, nil
}

// SaveBreed
func (p *Postgres) SaveBreed(breed models.Breed) (models.Breed, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Breed (breed)
		VALUES ($1)
		RETURNING breed_id, breed`
	row := p.db.QueryRowContext(context, sqlStatement,
		breed.Breed)

	var insertedBreed models.Breed
	err := row.Scan(&insertedBreed.ID, &insertedBreed.Breed)
	if err != nil {
		return models.Breed{}, err
	}

	return insertedBreed, nil
}

// SaveDogBreed
func (p *Postgres) SaveDogBreed(dogBreed models.DogBreed) (models.DogBreed, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO DogBreed (dog_id_fk, breed_id_fk)
		VALUES ($1, $2)
		RETURNING dog_id_fk, breed_id_fk`
	row := p.db.QueryRowContext(context, sqlStatement, dogBreed.DogIdFk, dogBreed.BreedIdFk)

	var insertedDogBreed models.DogBreed
	err := row.Scan(&insertedDogBreed.DogIdFk, &insertedDogBreed.BreedIdFk)
	if err != nil {
		return models.DogBreed{}, err
	}

	return insertedDogBreed, nil
}

// SaveExpense
func (p *Postgres) SaveExpense(expense models.Expense) (models.Expense, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Expense (date, vendor, description, amount_in_cents, dog_id_fk)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING date, vendor, description, amount_in_cents, dog_id_fk`
	row := p.db.QueryRowContext(context, sqlStatement,
		expense.Date,
		expense.Vendor,
		expense.Description,
		expense.AmountInCents,
		expense.DogIdFk)

	var insertedExpense models.Expense
	err := row.Scan(&insertedExpense.Date,
		&insertedExpense.Vendor,
		&insertedExpense.Description,
		&insertedExpense.AmountInCents,
		&insertedExpense.DogIdFk)
	if err != nil {
		return models.Expense{}, err
	}

	return insertedExpense, nil
}

// SaveAdoption
func (p *Postgres) SaveAdoption(adoption models.Adoption) (models.Adoption, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqlStatement := `
		INSERT INTO Adoption (date_adopted, application_number_fk, dog_id_fk)
		VALUES ($1, $2, $3)
		RETURNING date_adopted, application_number_fk, dog_id_fk`
	row := p.db.QueryRowContext(context, sqlStatement,
		adoption.DateAdopted,
		adoption.ApplicationNumberFk,
		adoption.DogIdFk)

	var insertedAdoption models.Adoption
	err := row.Scan(&insertedAdoption.DateAdopted,
		&insertedAdoption.ApplicationNumberFk,
		&insertedAdoption.DogIdFk)
	if err != nil {
		return models.Adoption{}, err
	}

	return insertedAdoption, nil
}

// GetAnimalControlReport
func (p *Postgres) GetAnimalControlReport() ([]models.AnimalControlReport, error) {
	sqlStatement := `
	WITH RECURSIVE MyDates AS (
		SELECT 
		d
		FROM
		GENERATE_SERIES(
			now(),
			now() - interval '6 months',
			interval '-1 month'
		) AS d
		ORDER BY d desc)
	Select
	to_char(d, 'Month') AS month,
	EXTRACT(YEAR FROM d) AS year,
	(SELECT COUNT(dog_id)
		FROM dog
		WHERE surrender_date >= date_trunc('month', d)
		AND surrender_date <= date_trunc('month', d) + interval '1 month' - interval '1 day'
		AND surrender_was_by_animal_control = true) 
		AS TotalAnimalControlDogs,
	(SELECT COUNT(Dog.dog_id)
		FROM dog
		INNER JOIN Adoption 
		ON Adoption.dog_id_fk = Dog.dog_id
		INNER JOIN Expense 
		ON Expense.dog_id_fk = Dog.dog_id
		WHERE date_adopted <= date_trunc('month', d) + interval '1 month' - interval '1 day' 
		AND date_adopted >=date_trunc('month', d) 
		AND surrender_was_by_animal_control = true
		AND  (date_adopted - surrender_date) >= 60) 
		AS TotalAnimalControlDogsGreaterThan60Days,
	(SELECT CAST (COALESCE(SUM(amount_in_cents), 0) AS REAL)/100
		FROM Expense
		INNER JOIN Dog 
		ON Expense.dog_id_fk = Dog.dog_id
		WHERE surrender_was_by_animal_control = true 
		AND Expense.date >=date_trunc('month', d) 
		AND Expense.date <= date_trunc('month', d) + interval '1 month' - interval '1 day') 
		AS TotalExpenses
	FROM MyDates;`

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	AnimalControlReportRows := make([]models.AnimalControlReport, 0)

	defer rows.Close()
	for rows.Next() {
		var AnimalControlReport models.AnimalControlReport
		if err := rows.Scan(
			&AnimalControlReport.Month,
			&AnimalControlReport.Year,
			&AnimalControlReport.DogsTotalCount,
			&AnimalControlReport.DogsSixtyDaysCount,
			&AnimalControlReport.Expenses); err != nil {
			return nil, err
		}
		AnimalControlReportRows = append(AnimalControlReportRows, AnimalControlReport)
	}
	return AnimalControlReportRows, nil
}

// GetAnimalControlReportDrillDownOne
func (p *Postgres) GetAnimalControlReportDrillDownOne(startDate string, endDate string) ([]models.AnimalControlReportDrillDown, error) {

	sqlStatement := `
	SELECT dog_id, sex, alteration_status, COALESCE(microchip_id,''),TO_CHAR(surrender_date, 'MM/DD/YYYY') AS SurrenderDate, STRING_AGG(breed, '/' ORDER BY breed) AS breed
	FROM Dog
	INNER JOIN DogBreed ON dog_id = dog_id_fk
	INNER JOIN Breed ON breed_id_fk = breed_id
	WHERE surrender_date >= $1 
	AND surrender_date < $2 
	AND surrender_was_by_animal_control = true
	GROUP BY dog_id 
	ORDER BY dog_id;`

	rows, err := p.db.Query(sqlStatement, startDate, endDate)
	if err != nil {
		return nil, err
	}
	AnimalControlReportRows := make([]models.AnimalControlReportDrillDown, 0)

	defer rows.Close()
	for rows.Next() {
		var AnimalControlReport models.AnimalControlReportDrillDown
		if err := rows.Scan(
			&AnimalControlReport.DogID,
			&AnimalControlReport.Sex,
			&AnimalControlReport.AlterationStatus,
			&AnimalControlReport.MicrochipID,
			&AnimalControlReport.SurrenderDate,
			&AnimalControlReport.Breed); err != nil {
			return nil, err
		}
		AnimalControlReportRows = append(AnimalControlReportRows, AnimalControlReport)
	}
	return AnimalControlReportRows, nil
}

// GetAnimalControlReportDrillDownTwo
func (p *Postgres) GetAnimalControlReportDrillDownTwo(startDate string, endDate string) ([]models.AnimalControlReportDrillDown, error) {
	sqlStatement := `
	SELECT dog_id, sex, alteration_status, COALESCE(microchip_id,''),TO_CHAR(surrender_date, 'MM/DD/YYYY') AS SurrenderDate, STRING_AGG(breed, '/' ORDER BY breed) AS breed, TO_CHAR(date_adopted, 'MM/DD/YYYY') AS AdoptionDate
	FROM Dog
	INNER JOIN DogBreed ON dog_id = dog_id_fk
	INNER JOIN Breed ON breed_id_fk = breed_id
	INNER JOIN Adoption ON Adoption.dog_id_fk = dog_id
	WHERE date_adopted < $2
	AND date_adopted >= $1
	AND  (date_adopted - surrender_date) >= 60 
	AND surrender_was_by_animal_control = true
	GROUP BY dog_id, date_adopted
	ORDER BY (date_adopted-surrender_date) DESC, dog_id DESC;`

	rows, err := p.db.Query(sqlStatement, startDate, endDate)
	if err != nil {
		return nil, err
	}
	AnimalControlReportRows := make([]models.AnimalControlReportDrillDown, 0)

	defer rows.Close()
	for rows.Next() {
		var AnimalControlReport models.AnimalControlReportDrillDown
		if err := rows.Scan(
			&AnimalControlReport.DogID,
			&AnimalControlReport.Sex,
			&AnimalControlReport.AlterationStatus,
			&AnimalControlReport.MicrochipID,
			&AnimalControlReport.SurrenderDate,
			&AnimalControlReport.Breed,
			&AnimalControlReport.AdoptionDate); err != nil {
			return nil, err
		}
		AnimalControlReportRows = append(AnimalControlReportRows, AnimalControlReport)
	}
	return AnimalControlReportRows, nil
}
