# Education Management System app

## Description

This repository contains the server app of the Education School Management System. The work is still in progress.

## Run application

First run the PostgreSQL container:
 
 `make postgres`
 
To create the School database :
 
 `make createdb`
 
 To make the migrations :
 
 `make migrateup`
 
 To run the tests:
 
 `make test`

`docker build -t school-api:latest .`

`docker run --name school-api -p 8080:8080 -e GIN_MODE=release  --network school-network -e DB_SOURCE="postgresql://root:secret@posgtgres:5432/school?sslmode=disable" school-api:latest`

`docker network create school-network`
