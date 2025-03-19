package dto

import "order-service/internal/model"

type CreateMovieRequest struct {
	Title    string
	Genre    string
	Director string
	Year     uint
}

func (req *CreateMovieRequest) ToModel() model.Movie {
	return model.Movie{
		Title:    req.Title,
		Genre:    req.Genre,
		Director: req.Director,
		Year:     req.Year,
	}
}

type UpdateMovieRequest struct {
	Title    string
	Genre    string
	Director string
	Year     uint
}

func (req *UpdateMovieRequest) ToModel() model.Movie {
	panic("implement me")
	return model.Movie{}
}
