# Education Management System app

## Description

This repository contains the server app of the Education School Management System. The work is still in progress.

Checkpoints:
* [x] Design the database schema, and apply the migrations
* [x] Implement SQL CRUD queries and use prepared statements
* [x] Implement all REST API endpoints.
* [x] Implement authentication middleware on the necessary endpoints
* [x] Implement Two Factor Authentication at Time-base OTP
* [x] Implement Role-Based-Access to REST endpoints
* [x] Dockerise the application
* [x] Deploy the Database to AWS RDS service
* [x] Deploy the dockerise application on AWS ECR on each push on main
* [x] Deploy the application service, ingress, and nginx lode balancer on AWS EKS
* [x] Configure automatic deploy to EKS


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
