# Optii Project

This project is designed to create APIs using the Gin Web Framework and Air for live code reloading in Golang.

## IMPORTANT NOTE:
I tested the API with the hardcoded request body in the same format as provided by Swagger, but encountered an error. As a result, part of the code is left with TODO comments since it's not clear how the request body should be constructed.

## Prerequisites

Before running this project, you must have the following installed:

- Go programming language
- Git (for version control)

## Getting Started

To get the application up and running, follow these steps:

### 1. Clone the Repository

First, clone the repository to your local machine:

```sh
git clone git@github.com:shiguemori/optii.git
```

### 2. Run the Application

Navigate to the project directory and run the application with the following command:

```sh
go run main.go
```

### Accessing the API Documentation

Once the application is running, you can access the API documentation through Swagger at the following URL:

- http://localhost:8080/swagger/index.html

Here, you can view the list of available endpoints, their expected request formats, and try out API calls directly from your browser.

### Environment Variables

For security and configuration management, the application uses environment variables. Ensure that you have a .env file at the root of your project with the necessary variables set, such as `OPTII_URL`, `OPTII_CLIENT_ID`, `OPTII_CLIENT_SECRET`, and `OPTII_AUTHENTICATION_URL`.

### Testing

Our project comes with a comprehensive test suite designed to ensure the highest standards of quality. To execute the tests and verify that all components behave as expected, follow the steps below:

Run All Tests: To execute all tests across the project, use the following command at the root of the project:

- ```go test ./...```
    - This will recursively run all tests in all subdirectories.

Test Coverage: To assess test coverage across the project, you can generate a coverage report by running:

- ```go test ./... -cover```
    - This command provides a summary of the coverage percentage per file.

Coverage Profile: If you need a detailed coverage report, you can generate a coverage profile with:

- ```go test ./... -coverprofile=coverage.out```
    - This command creates a file named coverage.out containing the coverage data.

View Coverage in Browser: For a visual representation of coverage, you can render the coverage profile in HTML format and view it in your web browser:
- ```go tool cover -html=coverage.out```
    - This will open a webpage displaying code coverage, with untested lines highlighted.

By regularly running these tests and analyzing coverage, we can maintain and improve the project's quality over time. Make sure to run the test suite before pushing code changes to ensure that all functionalities are working as intended.

#### Contact Information
- Email: shiguemori@hotmail.com
- Phone: +55 11 975676977
- LinkedIn: https://www.linkedin.com/in/shiguemori-fullstack-dev/
- GitHub: https://github.com/shiguemori

Your README is now clear and free of significant grammar errors. Be sure to adjust project and environment details to fit your specific needs.