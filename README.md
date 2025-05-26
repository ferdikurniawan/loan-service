# Loan Service

This is a simple implementation of a Loan Service which capables to provide the following use cases:

1. Borrower submits a new Loan `POST v1/loans`
2. Internal Staff to Approve a Loan `PATCH v1/loans/:loan_id/status`
3. Investor(s) to pledge fund to a loan based on the principal amount `POST v1/loans/:loan_id/investments`
4. Staff to disburse Loan to borrower `POST v1/loans/:loan_id/disburse`
5. Get Loan Detail `GET v1/loans/:loan_id`

## Project Structure

This app is structured by the way of Clean Architecture that is the controller / request handler, service layer and repository layer are separated. 

## How to Start the App

1. Rename `env.example` to `.env` file
2. Here, replace the value of `POSTGRES_URL` into the PostgreSQL DSN of your own (you need to set up an empty PostgreSQL DB for this one)
3. Use Golang [Migrate](https://github.com/golang-migrate/migrate) to migrate DB on your local like this `migrate -path migrations -database "your local DB DSN" -verbose up`
4. Build & run the app by run this command from your terminal `make all`. The app will be accessible via localhost:8080. Ensure that your Go version is at least 1.23.3
5. Your app is running and you can import Postman collection on this repo to look around the API specs of loan-service

## Unit Test

To verify the accuracy of the feature, the unit test is present inside `internal/services/loan_service_test.go` file. It covers all use cases of loan_service.go functionality