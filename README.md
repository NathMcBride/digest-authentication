# Go server side implementation of RFC7616 HTTP Digest Access Authentication

## Features
* Supports SHA-256
* Supports User name hashing
* Provides Authentication middleware, for use with standard library or gorilla mux.
* Integration and Unit tests written with Ginkgo

## Running the project

This project uses Go modules running ```go build``` or ```go run``` will download project dependencies and build/run the project.

## Running tests

To run integration tests ```ginkgo -r ./integration```

To run unit tests ```ginkgo -r ./src ```

To run all tests ```ginkgo -r```


