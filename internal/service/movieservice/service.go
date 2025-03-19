package movieservice

import "order-service/internal/model"

type movieRepo interface {
	GetMovie(id string) (*model.Movie, error)
	GetMovies() ([]model.Movie, error)
	CreateMovie(movie *model.Movie) (string, error)
	CreateMovies(movies []model.Movie) ([]string, error)
	UpdateMovie(id string, movie *model.Movie) (*model.Movie, error)
	DeleteMovie(id string) (bool, error)
}

type Service struct {
	movieRepo movieRepo
}

func New(movieRepo movieRepo) *Service {
	return &Service{
		movieRepo: movieRepo,
	}
}

func (s *Service) GetMovie(id string) (*model.Movie, error) {
	return s.movieRepo.GetMovie(id)
}

func (s *Service) GetMovies() ([]model.Movie, error) {
	return s.movieRepo.GetMovies()
}

func (s *Service) CreateMovie(movie *model.Movie) (string, error) {
	return s.movieRepo.CreateMovie(movie)
}

func (s *Service) CreateMovies(movies []model.Movie) ([]string, error) {
	return s.movieRepo.CreateMovies(movies)
}

func (s *Service) UpdateMovie(id string, movie *model.Movie) (*model.Movie, error) {
	return s.movieRepo.UpdateMovie(id, movie)
}
func (s *Service) DeleteMovie(id string) (bool, error) {
	return s.movieRepo.DeleteMovie(id)
}
