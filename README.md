# httpmiddleware-app-example

Simple microservices that insert and find record on MongoDB to ilustrade how is simple to use [httpmiddleware lib](https://github.com/LeoCBS/httpmiddleware/)

## Endpoints avaliable:

    GET "/storage/mongodb/database/:database/collection/:collection/id/:id"
    POST "/storage/mongodb/database/:database/collection/:collection"


## How to run the server

To start server you must export variable `MONGO_PORT` and just run go command:

    > go run main.go
    INFO[0000] starting server on port 8080

### Start server using docker-compose

To start this server using docker compose on port 8080 just run command below

    make run-compose

And run a wget to check if server was started

```
wget "http://localhost:8080/storage/mongodb/database/test/collection/person"    
--2023-04-11 16:42:07--  http://localhost:8080/storage/mongodb/database/test/collection/person
Resolving localhost (localhost)... 127.0.0.1
Connecting to localhost (localhost)|127.0.0.1|:8080... connected.
HTTP request sent, awaiting response... 404 Not Found
2023-04-11 16:42:07 ERROR 404: Not Found.
```

## Integration tests

To run integration tests:

    make check-integration

Integration tests will run a docker compose with this application and MongoDB
