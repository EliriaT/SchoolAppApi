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
* [x] Hide errors messages that contain  implementation sensitive information
* [x] Check and implement input validation of the requests
* [x] Dockerise the application
* [x] Deploy the Database to AWS RDS service
* [x] Deploy the dockerized application on AWS ECR on each push on main
* [x] Deploy the application service, ingress, and nginx lode balancer on AWS EKS
* [x] Configure automatic deploy to EKS
* [x] Configure and issue the HTTPS certificate
* [x] Configure the application level rate limiter and the nginx rate limiter
* [x] Add some minor general test.
* [x] Add expiration token duration of 5-15  minutes for the access token, and implement refresh session tokens
* [x] Validate password based the high entropy value
- [ ] Implement grpc for create user, login user
- [ ] Make Swagger: API Documentation
- [ ] Analyze other possible security threats

## Run application

First run the PostgreSQL container:
 
 `make postgres`
 
To create the School database :
 
 `make createdb`
 
 To make the migrations :
 
 `make migrateup`
 
 To run the tests:
 
 `make test`

To run the server:
 
 `make server`

To run the dockerized application:
 
 `docker compose up`


