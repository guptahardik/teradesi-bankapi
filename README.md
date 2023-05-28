#Go Application with Docker Compose
This is a Go application that simulates a simple banking system with features like creating an account, depositing funds, withdrawing funds, and checking account details. It exposes a RESTful API for communication.

##Prerequisites
Docker: Ensure that Docker is installed on your system.

###1. Getting Started
Clone the repository:
git clone <repository_url>

###2. Navigate to the project directory:

Use code with caution. Learn more
content_copy
cd go-banking-app

###3. Build the Docker image and run the application using Docker Compose:

Use code with caution. Learn more
content_copy
docker-compose up

This command will build the Docker image for the application and start the containers.

Once the containers are up and running, you can access the application API at http://localhost:7000.

Running Tests
To run the tests for the Go application, you can use Docker Compose:

Make sure the application containers are running (you should have executed docker-compose up in the previous steps).
Open a new terminal window and navigate to the project directory.
Run the tests using Docker Compose:
docker-compose run app go test ./...

This command will execute all the tests in the Go application and display the test results in the terminal.

API Endpoints
The following API endpoints are available:

GET /accounts/{id}: Retrieve the account details for the specified account ID.
POST /accounts/{id}/deposit: Deposit funds into the specified account ID. The amount should be provided in the request body as a JSON payload.
POST /accounts/{id}/withdraw: Withdraw funds from the specified account ID. The amount should be provided in the request body as a JSON payload.
POST /accounts/{id}/create: Create a new account with the specified account ID.
POST /accounts/{id}/add-money: Add money to the specified account ID. The amount should be provided in the request body as a JSON payload.
Refer to the source code for more details on the request and response formats for each endpoint.
