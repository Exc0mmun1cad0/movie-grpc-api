package dto

import "movie-service/internal/model"

type CreateMovieRequest struct {
	Title    string `validate:"required"`
	Genre    string `validate:"required"`
	Director string `validate:"required"`
	Year     uint32 `validate:"required,gte=1911"`
}

func (req *CreateMovieRequest) ToModel() *model.Movie {
	return &model.Movie{
		Title:    req.Title,
		Genre:    req.Genre,
		Director: req.Director,
		Year:     req.Year,
	}
}

type UpdateMovieRequest struct {
	ID       string `validate:"required,uuid"`
	Title    string `validate:"omitempty"`
	Genre    string `validate:"omitempty"`
	Director string `validate:"omitempty"`
	Year     uint32 `validate:"omitempty,gte=1911"`
}

func (req *UpdateMovieRequest) ToModel() *model.Movie {
	return &model.Movie{
		Title:    req.Title,
		Genre:    req.Genre,
		Director: req.Director,
		Year:     req.Year,
	}
}
