# go-api-practice
This project stasted as an exercise froked from [Kubucation]/(https://github.com/kubucation/go-rollercoaster-api)

## Requirements

To be able to show the desired features of curl this REST API must match a few
requirements:

* [x] `GET /coasters` returns list of coasters as JSON
* [x] `GET /coasters/{id}` returns details of specific coaster as JSON
* [x] `POST /coasters` accepts a new coaster to be added
* [x] `POST /coasters` returns status 415 if content is not `application/json`
* [x] `GET /admin` requires basic auth
* [x] `GET /coasters/random` redirects (Status 302) to a random coaster

### Data Types

A recipe object should look like this:
```json
{
  "id": "someid",
  "name": "name of the recipe",
  "ingredients": "a single string of all ingredients",
  "instructions": "a single string of all instructions",
}
```

### Persistence

There is no persistence, a temporary in-mem story is fine.
