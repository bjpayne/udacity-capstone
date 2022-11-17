# UDACITY - Capstone

Capstone CRM API project. Performs CRUD operations on a single `customers` table in a sqlite database.

Dependencies:

`database/sql` `encoding/json` `fmt` `github.com/gorilla/mux` `github.com/mattn/go-sqlite3` `log` `net/http` `strconv` `time`

To build:

```bash
$ go build
```

To run tests:

```bash
$ go test
```

To run:

```bash
$ go run main.go
```

### Postman collection

Included is a postman collection to easily test the various endpoints:

[postman collection](postman_collection.json)

### SQLite database

A sqlite database is used to persist customer data

[app.db](app.db)

column | type
------ | ----
id | INT PK
first_name | TEXT
last_name | TEXT
email | TEXT
phone | TEXT
role | TEXT
street | TEXT
city | TEXT
state | TEXT
zip | TEXT
contacted | INT
created_at | TEXT CURRENT_TIMESTAMP
