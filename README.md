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
