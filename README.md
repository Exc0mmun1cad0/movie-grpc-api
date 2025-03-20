### 🚀 **movie-grpc-api**
![Go Version](https://img.shields.io/badge/Go-1.23.4-blue)  
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-✔%ef%b8%8f-blue)  
![License](https://img.shields.io/badge/License-MIT-green)  

**GRPC API for creating, deleting, and retrieving movie data.**  
Supports both **unary** and **streaming** requests.

---

## 📌 **Movie Model**
In the **PostgreSQL** database, a `Movie` is stored as:

```go
type Movie struct {
    ID       uuid.UUID `db:"movie_id"`
    Title    string    `db:"title"`
    Genre    string    `db:"genre"`
    Director string    `db:"director"`
    Year     int       `db:"year"`
}
```

## 🔥 **API Methods**
| Method  | Type         | Description |
|---------|-------------|-------------|
| `GET`   | Unary & Streaming | Retrieve movie(s) |
| `POST`  | Unary & Streaming | Create new movie(s) |
| `UPDATE` | Unary | Update existing movie |
| `DELETE` | Unary | Remove movie from database |

---

## ⚙ **Configuration**
Service reads configuration from a `.env` file. Example from `.env.example`:  
```ini
ENV=prod                            # app environment
GRPC_PORT=50051                     # grpc server port
HTTP_PORT=8080                      # http server (grpc-gateway) port
POSTGRES_HOST=localhost             # host of db Postgres
POSTGRES_PORT=5432                  # port of db Postgres
POSTGRES_USER=your_user             # username for Postgres connection
POSTGRES_PASSWORD=your_password     # password for Postgres connection
POSTGRES_DB=movies_db               # name of Postgres database for connection
PGADMIN_PORT=5050                   # port for running PGAdmin to monitor Postgres
```

---

## 🛠 **Database**
- **PostgreSQL**  
- **Migrations** via [`golang-migrate`](https://github.com/golang-migrate/migrate)  
- **SQL Builder**: [`Squirrel`](https://github.com/Masterminds/squirrel) *(Used for learning purposes)*  

Migrations are done by the application itself at startup


---

## 🏰 **Project Structure**
Following **Clean Architecture**:
```
📺 movie-grpc-api
├── 📂 api/             # Protobuf contracts
├── 📂 cmd/             # Application entry point
├── 📂 internal/
│   ├── 📂 config/      # Configuration parsing
│   ├── 📂 model/       # Database models
│   ├── 📂 repository/  # Data access layer
│   ├── 📂 service/     # Business logic layer
│   ├── 📂 transport/   # GRPC & HTTP handlers
│   │   ├── 📂 dto/     # Data Transfer Objects
│   │   ├── 📂 grpc/
├── 📂 migrations/      # DB migration files
├── 📂 tests/           # Integration & unit tests
├── 📂 pkg/             # Reusable parts of code

├── Dockerfile
├── Makefile
└── README.md
```

---

## 🔍 **Health Check**
Health check endpoint for GRPC Gateway:
```
GET http://localhost:8080/health
```

---
