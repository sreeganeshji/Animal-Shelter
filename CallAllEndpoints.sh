# Call POST endpoints 1st,
# Then call GET endpoints

# Calls in this order:

# Volunteer
# Applicant
# Breed
# Dog
# Application
# DogBreed
# Expense
# Adoption

echo "Calling create volunteer"
curl --location --request POST 'http://localhost:8080/api/volunteer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "Dylan",
    "lastName": "Hantula",
    "email": "dhantula3@gatech.edu",
    "cell": "212-867-5309",
    "password": "test",
    "startDate": "2020-06-29",
    "isTrusted": true
}'


echo "Calling create Applicant"
curl --location --request POST 'http://localhost:8080/api/applicant' \
--header 'Content-Type: application/json' \
--data-raw '{
	"zip":"99999",
	"state":"AZ",
	"city":"Phoenix",
	"street":"123 Happy",
	"phoneNumber":"123 456 7890",
	"email":"Email@email.com",
	"lastName":"Nolan",
	"firstName":"Mike"
}'


echo "Calling create Breed"
curl --location --request POST 'http://localhost:8080/api/breed' \
--header 'Content-Type: application/json' \
--data-raw '{
    "breed":"fakebreed"
}'


echo "Calling create Dog"
curl --location --request POST 'http://localhost:8080/api/dog' \
--header 'Content-Type: application/json' \
--data-raw '{
	"name":"molly",
	"alterationStatus":true,
	"description":"Beautiful Golden Retriever",
	"dateOfBirth":"7-1-2020",
	"sex":"female",
	"microchipId":"123",
	"surrenderDate":"7-3-2020",
	"surrenderReason":"some reason",
	"surrenderWasByAnimalControl":false,
	"volunteerId":1
}'


echo "Calling create Application"
curl --location --request POST 'http://localhost:8080/api/application' \
--header 'Content-Type: application/json' \
--data-raw '{
    "date":"7-3-2020",
    "state":"pending approval",
    "coApplicantLastName":null,
    "coApplicantFirstName":null,
    "applicantIdFk":"1"
}'


echo "Calling create DogBreed"
curl --location --request POST 'http://localhost:8080/api/dogbreed' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dogIdFk":"1",
    "breedIdFk":"1"
}'


echo "Calling create Expense"
curl --location --request POST 'http://localhost:8080/api/expense' \
--header 'Content-Type: application/json' \
--data-raw '{
    "date":"7-3-2020",
	"vendor":"ABC Corp",
	"description":"Eye Surgry",
	"amountInCents":"10000",
	"dogIdFk":"1"
}'


echo "Calling create Adoption"
curl --location --request POST 'http://localhost:8080/api/adoption' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dateAdopted":"7-3-2020",
    "applicationNumberFk":"1",
    "dogIdFk":"1"
}'






echo "Calling get Adoption"
curl --location --request GET 'http://localhost:8080/api/volunteer' \
--header 'Content-Type: application/json'

echo "Calling get Applicant"
curl --location --request GET 'http://localhost:8080/api/applicant' \
--header 'Content-Type: application/json'

echo "Calling get Breed"
curl --location --request GET 'http://localhost:8080/api/breed' \
--header 'Content-Type: application/json'

echo "Calling get Dog"
curl --location --request GET 'http://localhost:8080/api/dog' \
--header 'Content-Type: application/json'

echo "Calling get Application"
curl --location --request GET 'http://localhost:8080/api/application' \
--header 'Content-Type: application/json'

echo "Calling get # DogBreed"
curl --location --request GET 'http://localhost:8080/api/dogbreed' \
--header 'Content-Type: application/json'

echo "Calling get Expense"
curl --location --request GET 'http://localhost:8080/api/expense' \
--header 'Content-Type: application/json'

echo "Calling get Adoption"
curl --location --request GET 'http://localhost:8080/api/adoption' \
--header 'Content-Type: application/json'