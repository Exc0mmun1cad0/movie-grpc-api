package dto

import "order-service/internal/model"

type CreateMovieRequest struct {
	Title    string
	Genre    string
	Director string
	Year     uint32
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
	ID       string
	Title    string
	Genre    string
	Director string
	Year     uint32
}

func (req *UpdateMovieRequest) ToModel() *model.Movie {
	return &model.Movie{
		Title:    req.Title,
		Genre:    req.Genre,
		Director: req.Director,
		Year:     req.Year,
	}
}
