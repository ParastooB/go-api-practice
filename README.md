# go-api-practice
This project stasted as an exercise froked from [Kubucation](https://github.com/kubucation/go-rollercoaster-api)

## Requirements

To be able to show the desired features of curl this REST API must match a few
requirements:

* [x] `GET /recipes` returns list of recipe as JSON
* [x] `GET /recipes/{id}` returns details of specific recipe as JSON
* [x] `POST /recipes` accepts a new recipe to be added
* [x] `DELETE /recipes/{id}` deletes the recipe with the given id 
* [x] `UPDATE /recipes/{id}` updates all the fields of the recipe with the given id 

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
### How to Run
Clone the repository and in the folder run the following command.
```console
go run server.go
```
You will need [Go](https://golang.org/doc/install) installed before you can run this code. If you desire, you could install [cURL](https://curl.se/) and use it for testing instead of a web browser.

### Useful Curl Commands
```console
curl localhost:8080/recipes
curl localhost:8080/recipes/id1 -X DELETE
curl localhost:8080/recipes -X GET | jq
curl localhost:8080/recipes/id1 -X UPDATE -d '{"Name": "Chicken Stew","Instructions":"X","Ingredients":"Y"}' -H "Content-Type: application/json"
curl localhost:8080/recipes -X POST -d '{"name": " Chicken Stew ","instructions":"mix, cook, stir, shred, garnish, serve","Ingredients":"butter,carrot, celery, salt, black pepper, garlic, flour, chicken, thyme, bay leaf, potato, chicken broth, parsley"}' -H "Content-Type: application/json"
```
### Persistence

There is no persistence, a temporary in-mem story is fine.
