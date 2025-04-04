package moviegrpc

import (
	"movie-service/internal/model"
	"movie-service/internal/transport/dto"
	"movie-service/pkg/pb"
)

func toPb(movie *model.Movie) *pb.Movie {
	return &pb.Movie{
		Id:       movie.ID,
		Title:    movie.Title,
		Genre:    movie.Genre,
		Director: movie.Director,
		Year:     movie.Year,
	}
}

func pbToCreate(in *pb.CreateMovieRequest) *dto.CreateMovieRequest {
	return &dto.CreateMovieRequest{
		Title:    in.GetTitle(),
		Genre:    in.GetGenre(),
		Director: in.GetDirector(),
		Year:     in.GetYear(),
	}
}

func pbToUpdate(in *pb.UpdateMovieRequest) *dto.UpdateMovieRequest {
	return &dto.UpdateMovieRequest{
		ID:       in.GetId(),
		Title:    in.GetTitle(),
		Genre:    in.GetGenre(),
		Director: in.GetDirector(),
		Year:     in.GetYear(),
	}
}
