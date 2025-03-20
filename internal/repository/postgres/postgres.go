package postgresrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"movie-service/internal/model"
	repo "movie-service/internal/repository"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) GetMovie(id string) (*model.Movie, error) {
	const op = "repository.postgres.GetMovie"

	query, args, err := r.builder.Select("*").
		From("movies").
		Where(sq.Eq{"movie_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to form sql query: %w", op, err)
	}

	var movie model.Movie
	err = r.db.Get(&movie, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: failed to get movie info by id: %w", op, repo.ErrMovieNotExists)
		}

		return nil, fmt.Errorf("%s: failed to get movie info by id: %w", op, err)
	}

	return &movie, nil
}

func (r *Repository) GetMovies() ([]model.Movie, error) {
	const op = "repository.postgres.GetMovies"

	query, _, err := r.builder.Select("*").
		From("movies").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to form sql query: %w", op, err)
	}

	var movies []model.Movie
	err = r.db.Select(&movies, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get info about movies: %w", op, err)
	}

	return movies, nil
}

func (r *Repository) CreateMovie(movie *model.Movie) (string, error) {
	const op = "repository.postgres.CreateMovie"

	query, args, err := r.builder.Insert("movies").
		Columns("title", "genre", "director", "year").
		Values(movie.Title, movie.Genre, movie.Director, movie.Year).
		Suffix("RETURNING movie_id").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: failed to form sql query: %w", op, err)
	}

	var movieID string
	err = r.db.Get(&movieID, query, args...)
	if err != nil {
		return "", fmt.Errorf("%s: failed to add movie info: %w", op, err)
	}

	return movieID, nil
}

func (r *Repository) CreateMovies(movies []model.Movie) ([]string, error) {
	const op = "repository.postgres.CreateMovies"

	builder := r.builder.Insert("movies").Columns("title", "genre", "director", "year")
	for _, movie := range movies {
		builder = builder.Values(movie.Title, movie.Genre, movie.Director, movie.Year)
	}
	builder = builder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to form query: %w", op, err)
	}

	var movieIDs []string
	err = r.db.Select(&movieIDs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to add info about movies: %w", op, err)
	}

	return movieIDs, nil
}

func (r *Repository) UpdateMovie(id string, movie *model.Movie) (*model.Movie, error) {
	const op = "repository.postgres.UpdateMovide"

	builder := r.builder.Update("movies")
	if movie.Title != "" {
		builder = builder.Set("title", movie.Title)
	}
	if movie.Genre != "" {
		builder = builder.Set("genre", movie.Genre)
	}
	if movie.Director != "" {
		builder = builder.Set("director", movie.Director)
	}
	if movie.Year != 0 {
		builder = builder.Set("year", movie.Year)
	}

	query, args, err := builder.
		Where(sq.Eq{"movie_id": id}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to form query: %w", op, err)
	}

	var newMovie model.Movie
	err = r.db.Get(&newMovie, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: failed to update movie info: %w", op, repo.ErrMovieNotExists)
		}

		return nil, fmt.Errorf("%s: failed to update movie info: %w", op, err)
	}

	return &newMovie, nil
}

// Returning bool val indicates whether movie info was deleted or not
func (r *Repository) DeleteMovie(id string) (bool, error) {
	const op = "repository.postgres.DeleteMovie"

	query, args, err := r.builder.Delete("movies").
		Where(sq.Eq{"movie_id": id}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: failed to form query: %w", op, err)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("%s: failed to delete movie info: %w", op, err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: failed to get number of updated rows: %w", op, err)
	}

	if num == 0 {
		return false, nil
	}

	return true, nil
}
