# overview
This project implements a web server using **Reactjs** frontend and **Go** for the backend.

# Team
Ganesh, Dylan, Stevenson and Maichel

## Back-End
All back-end code exists in the `back-end/` directory. It is a prerequisite that the Postgres database be
running and loaded with the schema before the app can function properly. To run the Postgres database with
the schema, run `./StartDB.sh`.

The backend code is split into 4 distinct packages: *database*, *models*, *rest*, and *service*. The names of the packages indicate their respective charters: *database* handles the persistence layer, *models* defines the data structures, *rest* handles the transport layer, and *service* handles the business logic.

To run the app, simply cd into `back-end/` and
run `go run main.go`. It will spin up a server running on port 8080 and listen for HTTP requests. To 
check if it is working, you can make requests to create a volunteer(s) and then retrieve the volunteer(s) back.

### Create a Volunteer
```
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
```

### Get All Volunteers
```
curl --location --request GET 'http://localhost:8080/api/volunteer' \
--header 'Content-Type: application/json'
```
