How to Run in Local:

    1. Clone the repository 
    2. Download your service account json key, rename to creds.json, place it in root directory
    3. Enable your service account to be able to access Cloud Run Api, Cloud Build Api, IAM secrets api, firestore etc
    4. If you want to use Secrets Manager, enable it for your Service account and place the following configs in it, otherwise place them in a .env file in root folder

        PORT=8080
        HTTP_TIMEOUT=12s
        API_URL=https://jsonplaceholder.typicode.com/posts
        SOURCE=placeholder_api

    -As Application:
        Run the following command by replacing the environment variables values for your account

        $env:GOOGLE_APPLICATION_CREDENTIALS = "C:/Users/HP/Cyderes/test-ingestion/creds.json"
        $env:GOOGLE_CLOUD_PROJECT = "cyderes-ingestion"
        go run ./cmd/main/main.go
