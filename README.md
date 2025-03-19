# movie-grpc-api

GRPC API for creating/deleting/getting info about movies 

In database `Movie` is presented like this:
```go
type Movie struct {
    ID uuid
    Title string
    Genre string
    Director string
    Year int
}
```
Methods:
- `GET` Movie(s) (unary & streaming)
- `POST` Movie(s) (unary & streaming)
- `UPDATE` Movie
- `DELETE` Movie

## Configuration
Service reads config from `.env` file. Here is an example in `.env.example`

Vars:
- `GRPC_PORT` - grpc server port
- `HTTP_PORT` - http server (grpc-gateway) port
- `POSTGRES_HOST` - host of db *Postgres*
- `POSTGRES_PORT` - port of db *Postgres*
- `POSTGRES_USER` - username for *Postgres* connection
- `POSTGRES_PASSWORD` - password for *Postgres* connection
- `POSTGRES_DB` - name of *Postgres* database for connection
- `PGADMIN_PORT` - port for running *PGAdmin* to monitor *Postgres*
