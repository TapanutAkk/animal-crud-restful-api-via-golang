# Getting Started

## Install go
```bash
https://go.dev/dl/
```

## Run Database Container (Local)
```bash
# Start PostgreSQL Container on port 5432
# Username: myuser
# Password: mypassword
# Database: go_db
docker run --name my-postgres -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=go_db -p 5432:5432 -d postgres:latest
```

## .ENV file
Copy .env.example to .env

## Run Dependencies
```bash
go mod tidy
```

## Run Server
```bash
air
```