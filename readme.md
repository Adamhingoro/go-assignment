## Golang Service Center Management Application

### Overview

This project is a Golang application designed to manage service centers for multiple companies. It provides an interface for an admin user to perform CRUD operations on companies and customer care representatives. The application also includes an authentication system and a background job to monitor user activity.

### Features

1. **Authentication System**: Secure authentication for the admin user and customer care representatives.
   
2. **CRUD Operations for Companies**: Admin can create, read, update, and delete company records.

3. **CRUD Operations for Customer Care Representatives**: Admin can manage customer care representatives associated with different companies.

4. **Background Job**: A scheduled job runs every 10 minutes to check the `last_seen_at` property of users. If a user has not been active for 5 hours, they are marked as unavailable.

5. **Unit Testing**: The application is equipped with unit tests, and coverage reports are provided to ensure code quality.

6. **Open-Ended Design**: Certain aspects of the application are intentionally left open for interpretation. Any ambiguities should be noted in this `README.md` file.

### Bonus Points

- **Docker Integration**: The application is set up using Docker for both the database and the application environment.
  
- **Dependency Injection**: Utilizes the Uber-Dig framework for dependency injection to enhance modularity and testability.

### Technical Details

- **Golang Version**: The project is built using Golang version 1.21.5.
  
- **Database Operations**: GORM is used for all database-related operations, providing a robust ORM solution.

- **Graceful Shutdowns**: The application handles graceful shutdowns using channels and signals to ensure a smooth termination process.

### Getting Started

To get started with the project, follow the setup instructions below.

#### Pre-requisits 
 - the application uses the Postgres database to store the data, make sure to spin-up the postgres container using the following command. 
 `docker-compose up -d postgres` 

#### Running the application locally 
 - To run the application locally use the command `go run main.go`, make sure `.env` file exists with correct values. (already supplied with initial values)

#### Running the application containerized 
 - To run the application using compose use the following command `docker-compose up -d my_service`

### Docker Compose Components

1. **Services**:
   - **postgres**: Runs the PostgreSQL database.
     - **image**: Uses PostgreSQL version 15.
     - **environment**: Sets up database credentials.
     - **ports**: Exposes port 5432 for database access.
   - **my_service**: Runs the Golang application.
     - **build**: Specifies the context and Dockerfile for building the service.
     - **ports**: Exposes port 8080 for the application.
     - **depends_on**: Ensures the application starts after the database.

2. **Volumes**: 
   - **postgres_data**: Persists PostgreSQL data.

3. **Networks**: 
   - **my_network**: Allows communication between services.



### Extras
 the file `testncover` is used to run all unit tests and build the coverage.html 


 Postman collection of API's can be found in `Golang Assignment.postman_collection.json`


