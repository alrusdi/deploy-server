# Go Web Server for Automated Application Deployment Using GitHub Actions

Develop a Go-based web server for automating application deployment via GitHub Actions. The server should meet the following requirements:

## API with Basic Auth
- Protect the API using Basic Auth.
- The username and password for authentication must be read from a `.htpasswd` file located in the same directory as the executable.

## Handling GitHub Requests
- The server should accept POST requests with a parameter `id` that specifies the application ID.
- Using the provided ID, the server should retrieve deployment commands from a local SQLite database and create a record in the `deployment_session` table to log the start of the deployment. The API should return the ID of the created record.

## Command Execution
- Deployment commands must be executed sequentially in the background.
- The result of each command execution should be recorded in a separate table `deployment_steps` and reference the created `deployment_session` record.
- The shell name and its arguments for command execution should be specified in a `config.yaml` configuration file.
- If a command execution fails, log the error.

## Logging
- All logs must be printed to `stdout` and simultaneously written to a text file named `error.log`.
- Logs should include information about command execution, errors, and deployment status.

## Database
- Use SQLite as the database for data storage.
- Create three tables:
  - `applications`: stores the application ID, name, and deployment commands.
  - `deployment_session`: stores information about the current deployment (unique record ID, application ID, start time, end time, overall status).
  - `deployment_steps`: stores the results of each command execution, including the unique record ID, `deployment_session` ID, command, output, status, and execution time.

## Configuration
- Shell configuration (name and arguments) should be stored in a `config.yaml` file.
- The configuration file should allow specifying multiple arguments for command execution.

## Additional Requirements
- Use Go Modules for dependency management.
- The program must be easy to build and run, with clear instructions provided.
- Write instructions for building and running the application.

## Example Usage
After starting the server, a POST request to `/deploy?id=1` should trigger the server to:
- Find commands for the application with ID=1 in the database.
- Execute the commands sequentially.
- Record the results in the database and logs.

## Expected Result
A Go web server that meets all the requirements. The code should be well-structured, with comments and examples. Include instructions for building and running the application.

## Automated Tests
- Include automated tests to verify the application's functionality.
- The tests should use an in-memory database, and each test should assume an empty database, requiring all necessary data to be created from scratch for each test.

## GitHub Action
- Include a GitHub Action workflow that automatically builds and tests the application.
- The workflow must install SQLite for testing and include a step to check test coverage.
- The GitHub Action should run on pushes to the `main` branch and pull requests targeting the `main` branch.
