This is the backEnd for the CRUD app [Movies](https://moviesplace.shop/).
## Getting Started

To load the database to postgreSQL(make sure that the movies db is created).
Choose one of the following:

```bash
psql -d movies -f 040\ go_movies.sql
psql -U postgres -d movies -f 040\ go_movies.sql
```

Change the config.db.dsn string in main.go to your postgreSQL database.

Config a unique secrete and run the development server:

```bash
export JWT_SECRET_KEY='my-secret-key'
go run ./cmd/api/.

```

To make the build and generate the binary file:

```bash
env GOOS=linux GOARCH=amd64 go build -o gomovies ./cmd/api
```


