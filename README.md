# Log Regex Parse + Upload + Top Errors

## Set up and Run API locally

1. Open your command line terminal.
2. Navigate to the directory where you want to clone the project.
3. Clone the GitHub repository using its URL
4. Navigate to the project directory that you just cloned.
5. Install project dependencies using the below command:

    ```
    go mod tidy
    ```
6. Inside the project folder, create a new file named `.env`. This file should be located at the same level as the `go.mod` file. In the `.env` file, make sure to set all the required environment variables according to the project's configuration. You can refer to the `sample.env` file as a reference for the environment variables that need to be defined in your `.env` file.
7. This application uses AWS SDK to communicate with AWS S3 service. This SDK requires you to have AWS credentials stored in the `~/.aws/credentials` file. Please create this file with the following content:
    ```
    [default]
    aws_access_key_id = <YOUR_ACCESS_KEY_ID>
    aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
    ```

    Replace <YOUR_ACCESS_KEY_ID> and <YOUR_SECRET_ACCESS_KEY> with your credentials.
7. Start the API server:

    ```
    go run .
    ```
    The terminal should display the message `Server is running on <SERVER_PORT>`

You can now access the API endpoints using a tool like Postman or via the `curl` command.

## API Documentation 

An overview of the endpoints, required parameters, response formats, and example requests/responses of the API.

### Base URL

```
http://<HOSTNAME>:<SERVER_PORT>
```

### Endpoints

#### Health Check
- **URL**: `/health`
- **Method**: `GET`
- **Description**: This API endpoint allows users to check the health of the server to ensure it is running properly.
- **Example Request**:
    ```
    GET /health
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "message": "Server is up and running"
    }
    ```

#### Upload Logs
- **URL**: `/upload`
- **Method**: `POST`
- **Description**: This API endpoint allows you to upload logs to S3.
- **Request Body**: The request must be a `multipart/form-data` type, with a key named `sample-file` containing the log file to be uploaded.
- **Example Request**:
    ```
    POST /upload
    Content-Type: multipart/form-data

    Form Data:
    Key: sample-file, Value: [log file content]
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "message": "logs uploaded to S3 successfully"
    }
    ```

#### Get Top Error
- **URL**: `/top-error`
- **Method**: `GET`
- **Description**: This API endpoint allows you to retrieve the most frequently occurring error from the S3.
- **Example Request**:
    ```
    GET /top-error
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "top-error": "Request from 10.0.0.2 failed with status code 404 for /page-not-found.html",
        "count": 3,
        "service": "apache",
        "severity": "INFO"
    }
    ```
