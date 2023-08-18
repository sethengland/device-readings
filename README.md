# Thank you for checking out my project.

# Summary
This project implements a local API using Go and the Gin web framework. The API provides endpoints for posting and retrieving readings data from different devices.

# Instructions
- make sure you have Go installed
- run `go build` and then `./brightwheel-device-readings` or run it manually with `go run main.go`
- run my test cases with `go test`
- the api will be running on `localhost:8080`

### Here is an example of some curl commands I used to manually test my API
Formatting `curl` commands in Windows is difficult so you may have to change yours slightly.

Getting readings:

    C:\Users\sethe>curl -v localhost:8080/get-reading/adieu

Posting readings:

    C:\Users\sethe>curl -v -X POST localhost:8080/post-readings -d "{\"id\":\"adieu\", \"readings\":[{\"timestamp\":\"2021-09-29T16:08:15+01:00\", \"count\":5}]}"

# API Documentation
## POST /post-readings

This endpoint allows you to post readings data for a specific device.

Request:

    Method: POST
    Path: /post-readings
    Headers:
        Content-Type: application/json

Request Body Example:

```json
{
  "id": "adieu",
  "readings": [
    {
      "timestamp": "2023-08-16T12:00:00Z",
      "count": 42
    }
  ]
}
```

Response:

    Status Code: 200 (OK)
    Body: Returns the list of readings for the posted device.

GET /get-reading/:id

This endpoint allows you to retrieve readings data for a specific device.

Request:

    Method: GET
    Path: /get-reading/:id

URL Parameters:

    id (string): The device ID for which to retrieve readings.

Response:

    Status Code: 200 (OK)
    Body: Returns the list of readings for the specified device.

In case of a non-existing device ID, a 404 status with the following message will be returned:


```json
{
  "message": "no device with that id"
}
```

# Project Reflection
### What were the major roadblocks/where did I spend the bulk of the time
- Setting up the router and assigning routes was made very easy by using Gin.
- I spent the most amount of time designing the data structure for my data map. I wanted to map each device to a `Set` of readings, but sets do not exist on their own in Go. After some research, I found that using a nested map with empty structs as the map values can give me identical functionality to a true `Set`.
- I wanted to map each device to a `Set` of readings because of the considerations suggested by the project description. A `Set` can easily ignore duplicate readings.
- I chose to use the unique combination of `timestamp` and `count` as the key to my `Set` implementation. It was a concious choice to allow `Readings` with identical `timestamps` but different `counts` to exist independently.
- The other thing I spent the most time on was creating a robust enough test suite (I also found a small bug along the way which made it very worth it!)
- This documentation is putting me slightly over the time cap but it's worth it because it looks nice :)

### If you had more time, what part of your project would you refactor? What other tradeoffs did you make?
- I regretfully ignored returning the `Readings` sorted by timestamp, that would probably be the first improvement I would make even though it wasn't a requirement.
- I would also add an integrated test suite that alternatively calls each endpoint to verify that the data map state is correct after every `POST`.