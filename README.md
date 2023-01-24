# Docs API (Golang)
A REST API reated with Golang to demonstrate CRUD functionality with MongoDB.

## Setup
To start of, one needs to [install Go](https://go.dev/doc/install)
Once installed and setup, `git clone git@github.com:n8e/docs_api_golang.git`
`cd docs_api_golang` and install dependencies:
```
go install
```
Create a `.env` file such as below:
```
JWT_SECRET_KEY="XXXXXXXX"
MONGO_URL="XXXXXXX" # This could be a Mongo Cluster created on https://cloud.mongodb.com, or use local Mongo DB
# MONGO_URL="mongodb://localhost:27017"
```
With the env file containing correct values we are ready to run the API
```
go run main.go
```