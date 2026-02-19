# Go server

This is a basic example of how to implement a simple server in Golang with the standard library of `"net/http"`. For this example I chose to create a simple CRUD todo list you can run on the terminal. For an easier aproach on the database side I used sqlite so that you can have a local SQL db on your pc instead of a `CSV` or `JSON` file.

## How to run the server

First you need to download the repo:

```bash
git clone https://github.com/DJ2513/go_server.git
```

Open the project on you IDE, I recommend `VSC` or `Cursor` and inside the terminal run this command:

```bash
# This complies and runs all go files needed
go run cmd/*.go 

# You can also run the make file if you make any changes
# that way it formats, verifies, builds and runs the code for you
make
```

On the terminal, it will apear a link you can click by combining `Control + Click`. On Mac use `command + click`.

## How to consume the API

Once you run the server, now you can visit the following paths:

```text
- GET /health
- POST /todos
- GET /todos
- GET /todos/{id}
- DELETE /todos/{id}
- POST /todos/{id}/tasks
- DELETE /todos/{id}/tasks/{task_id}
```

Now you can use `Curl` or `Postman`, whatever you like, to consume the API.

```bash
# Curl Example
curl -X POST localhost:8080/todos -d '{"title":"Groceries"}'

# Get all lists
curl localhost:8080/todos

# Delete a list
curl -X DELETE localhost:8080/todos/<id>
```
