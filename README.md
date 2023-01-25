# Education Management System app

## Description

This repository contains the server app of the Education School Management System. The work is still in progress.

## Checkpoints:

* [x] Design the database schema, and apply the migrations
* [x] Implement SQL CRUD queries and use prepared statements
* [x] Implement all REST API endpoints.
* [x] Implement authentication middleware on the necessary endpoints
* [x] Implement Two Factor Authentication based on Time-based OTP
* [x] Implement Role-Based-Access to REST endpoints
* [x] Hide errors messages that contain  implementation sensitive information
* [x] Check and implement input validation of the requests
* [x] Dockerise the application
* [x] Deploy the Database to AWS RDS service
* [x] Deploy the dockerized application on AWS ECR on each push on main
* [x] Deploy the application service, ingress, and nginx load balancer on AWS EKS
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

## Endpoints
#### User routes:

`– POST /users` =>  Creates a users. Authentication middleware is applied to the endpoint. Only some roles can
create users. After registration, an email is send with the 2FA’s qrcode, and link for changing the
auto-generated password.

`– GET /users/:id` =>  Get a user information by id. Authentication middleware is applied to the endpoint. A user can
only get his user’s information.

`– POST /users/login` => Login to the service. This endpoint returns the Access token. This token is
further used to pass the 2FA endpoint.

`– POST /users/accountrecovery/:email/:token` =>  This endpoint is used to change a user’s
password. In the Body the new password is provided. The URI token parameter is uniquely
linked to the user’s email.

`– POST /users/accountrecovery` => Request a link to change the user’s password. In the body the
user’s email is provided. A link is sent to the user’s email, if such a user with such an email exists.

`– POST /users/twofactor` =>  Authentication middleware is applied to this endpoint. To access this endpoint, the user has to
provide the access token received from the login stage. After passing a valid Time-based OTP, the
access token will be authenticated, and returned as a response, along with a refresh token. A new
session with this refresh token is created in the database.

`– POST/token/renew access` => Request a new access token. In the body the refresh token is provided. Its expiration is checked,
as well as its validity.

#### School routes:
Authentication middleware is applied to the Schools routes.

`– POST /schools` => Creates a school. Admin only access

`– GET /schools/:id` => Get a school information by id. Admin or director’s school only access. 

`– GET /schools/?page id=1&page_size=10` => List all schools paginated. Admin only access.

#### Class routes:
Authentication middleware is applied to the Class routes.

`– POST /class ` => Creates a class. Director’s and manager’s school only access.

`– GET /class/:id` => Used to view information about a class with students list and contact information. Only a Head
Teacher, a director, or manager can view the response.

`– GET /class` => Response is based on user’s role. If it is a director or school manager, all classes are returned,
otherwise the user’s class (even in case of the head teacher).

`– PUT /class` =>  Used to change a class’ head teacher

#### Semester routes:
Authentication middleware is applied to the Semester routes.

`– POST /semester` => Creates a semester. Access limited to director and manager

`– GET /semester` => Get list of semesters

`– GET /semester/current` => Get current semester

#### Course routes:
Authentification middleware is applied to the Course routes.

`– POST /course` => Creates a course. Access limited to director and manager

`– GET /course/:id` => Get Student’s list together with their marks. A student will receive in the response only their
marks.

`– GET /course` => List courses of a teacher, of a student, of the class’s courses for headteacher, and
all courses for director and manager.

`– PUT /course` => Change a course teacher. Director and manager only access.

#### Lesson routes:
Authentification middleware is applied to the Lesson routes.

`– POST /lesson` => Creates a course. Access limited to director and manager

`– GET /lesson/course/:id` => Get all the lessons of a course. Director and manager only access.

`– GET /lesson` => Return a student and a teacher’s schedule. Head teachers can see class’s schedule.

`– PUT /lesson` => Change a lesson start, end hour , and day of week. Director and manager only access.

#### Marks routes:
Authentification middleware is applied to the Marks routes. Access limited to teacher who is teaching
the course.

`– POST /mark`  => Creates a mark or absence .

`– PUT /mark` => Change a mark or absence.

`– DELETE /mark` => Delete a mark or absence.

#### Other routes:
Authentification middleware is applied to the Marks routes.

`– GET /roles` => Get all the roles in the system.

`– GET /teachers` => Get all the teacher and stuff from a school
