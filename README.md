# CovidCases
This is an api server that stores and fetches the covid cases by state

# Steps to run
1. go get ..
2. setup mongo, update mongo uri in pkg/mongo/mongo.go
3. Set environment variables
    ```
    LOCATIONIQ_API_KEY = <your-token>
    ```
    Get your api key for geocode from https://locationiq.com/register

5. go run . 

# API documentation

1. Return the covid cases for user location

    GET /covid/user/{lat},{lon}

    Response (JSON):

    ```
    "state": {
        "code": "UP",
        "name": "Uttar Pradesh"
        "num_covid_cases": 26,
        "last_updated": "2021-01-01T15:30:00+05:30"
    },
    "country": {
        "code": "IN",
        "name": "India",
        "num_covid_cases": 1126,
        "last_updated": "2021-01-05T14:26:55+05:30"
    }
    ```
    
2. Update the covid data from source
    
    POST /covid/update

    Response:
    ```
    "Done"
    ```

# Improvements
 * Swagger api documentation
 * Add configuration parameters
 * Add api caching
 * Deploy to heroku
 * Setup environment in docker

