# CovidCases
This is an api server that stores and fetches the covid cases store by state

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
TODO



# Improvements
 * Swagger api documentation
 * Add configuration parameters
 * Add api caching
 * Deploy to heroku
 * Setup environment in docker

