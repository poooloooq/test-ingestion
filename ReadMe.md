How to Run in Local:

    1. Clone the repository 
    2. Download your service account json key, rename to creds.json, place it in root directory
    3. Enable your service account to be able to access Cloud Run Api, Cloud Build Api, IAM secrets api, firestore etc
    4. If you want to use Secrets Manager, enable it for your Service account and place the following configs in it, otherwise place them in a .env file in root folder

        PORT=8080
        HTTP_TIMEOUT=12s
        API_URL=https://jsonplaceholder.typicode.com/posts
        SOURCE=placeholder_api
    
    5. Run the application using one of the below methods (Powershell commands)  
    6. Make a http get request to /posts to start the ingestion
    7. After ingestion is complete and we get a response make a http get request to /posts/get

    -As Application:
        Run the following command by replacing the environment variables values for your account

        go test ./...
        $env:GOOGLE_APPLICATION_CREDENTIALS = "C:/Users/HP/Cyderes/test-ingestion/creds.json"
        $env:GOOGLE_CLOUD_PROJECT = "cyderes-ingestion"
        go run ./cmd/main/main.go

    -As Docker Container
        Start Docker Engine and Run the following command. Env variables are set in docker-compose

        docker compose up --build


How to Run in Cloud

    1. This repo has a github actions yaml file which has steps to deploy the application to cloud run.
    2. All secrets mentioned in the deployment yaml file is stored as repository secrets in github.
    3. Start the workflow to deploy the application to cloud run and perform integration tests.
    4. The cloud env doesnt have a .env file so it loads all config from Secrets Manager.
    5. The application is already deployed in Pratik's GCP account. Test the endpoints using the following url :

        ingestion: https://test-ingestion-service-507986862354.us-central1.run.app/posts
        fetch:     https://test-ingestion-service-507986862354.us-central1.run.app/posts/get


Documentation-

    Endpoints:

        1. /posts - Starts the ingestion process, calls api, transforms data and stores in firestore.

        2. /posts/get - Gets all data in Firestore database.

    Transformation Logic:

        1. Adds source field whose value is a config
        2. Adds ingested_at field and populates current time

    Database Schema: 
        
        userId          int    
	    id              int    
	    title           string 
	    body            string 
        ingested_at     Time 
	    source          string 

    Hardest Part to Implement:

        1. The hardest part was coding in Golang. Coming from Java background, the Golang naming standards, dependency management and method/exception handling were conceptually different from a Java Springboot environment where spring handles beans and configuration.

        2. Slightly less difficult was to deploy application in cloud run and integrating the secrets manager to fetch configs. Lot of credential management and IAM role management was required for the cloud setup.

    Trade Offs:

        1. Responsiveness vs Memory Efficiency

        A multi threaded approach can be considered since the ingestion takes a lot of time due to high volume of data which can increase in a live environment. Using multiple thread will decrease response time but also increase the cpu and memory utilization

        2. Simplicity vs Scalability

        The logic which inserts and retrieves from firestore can be moved to separate service for modularity. Since most memory utilisation is during database calls, to scale it up we have to scale entire service. Having database logic in separate service helps to scale only that particular service for high availability.

        3. Maintainability vs Optimization

        Currently the deployment steps are all in one file. While easier to maintain, it is not optimised aince traditionally we segregate the CI and CD process.

    Improvement Points:

        1. A scheduler could be added along with the /posts endpoint to trigger the ingestion process at regular intervals. A cron can be introduced to set intervals.

        2. A multithreaded approach can be introduced to insert data to firestore to decrease the response times.

        3. Currently the cloud endpoint is not authenticated. It is better to have some authentication like MTls  to ensure only permitted users are able to access the endpoints.

        4. Automate the deployment process, currently we have to manually start pipeline to deploy in cloud. CI and CD can be segregated and CI can be triggered for branch pushes and merge requests.

        5. Expose the endpoint via a gateway and use load balancer infront of it instead of directly accessing it unauthenticated.

        6. Use kubernetes manifests preferably helm to deploy apps in a clustered environment. That will help with autoscalling and autodeploying

        7. Have delete functionality

Bonus-

    1. Implemented CI/CD using github actions
    2. Added rest endpoint (/posts/get) to fetch all data from firestore
    3. How to track ingestion?

            1. If you need data of individual record last ingested then query using the ingested_at for latest record.

            2. If you need data about the ingestion itself then create a separate metadata collection and store the details of ingestion like count and time for each ingestion

            3. If you need data for all records that are part of ingestion, there are two ways

                - Add a trace_id field to the existing collection and populate with unique value, then quesry using this fiend and ingested_at to get latest ingestion records.
                - The ingestion method returns all the ingested records, that can be stored in a json file and updated in some s3 bucket which gets replaced with each new ingestion.

    4. Added integration tests to the deployment yaml file. It calls the /posts to ingest data and then calls /posts/get to see if data is ingested to firestore.