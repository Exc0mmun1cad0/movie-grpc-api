package moviegrpc

import (
	"context"
	"errors"
	"movie-service/internal/model"
	repo "movie-service/internal/repository"
	"movie-service/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetMovie(id string) (*model.Movie, error)
	GetMovies() ([]model.Movie, error)
	CreateMovie(movie *model.Movie) (string, error)
	CreateMovies(movies []model.Movie) ([]string, error)
	UpdateMovie(id string, movie *model.Movie) (*model.Movie, error)
	DeleteMovie(id string) (bool, error)
}

type server struct {
	pb.UnimplementedMovieServiceServer
	service Service
}

func Register(gRPCServer *grpc.Server, service Service) {
	pb.RegisterMovieServiceServer(gRPCServer, &server{service: service})
}

func (srv *server) CreateMovie(ctx context.Context, in *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	// TODO: add validation
	newMovie := pbToCreate(in)

	newID, err := srv.service.CreateMovie(newMovie.ToModel())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add movie info")
	}

	return &pb.CreateMovieResponse{Id: newID}, nil
}

func (srv *server) GetMovie(ctx context.Context, in *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	// TODO: add uuid validation
	id := in.GetId()

	movie, err := srv.service.GetMovie(id)
	if err != nil {
		if errors.Is(err, repo.ErrMovieNotExists) {
			return nil, status.Error(codes.NotFound, "movie not found")
		}

		return nil, status.Error(codes.Internal, "failed to get movie info")
	}

	return &pb.GetMovieResponse{Movie: toPb(movie)}, nil
}

func (srv *server) UpdateMovie(ctx context.Context, in *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	// TODO: add model movie andd uuid validation
	id := in.GetId()
	movie := pbToUpdate(in)

	newMovie, err := srv.service.UpdateMovie(id, movie.ToModel())
	if err != nil {
		if errors.Is(err, repo.ErrMovieNotExists) {
			return nil, status.Error(codes.NotFound, "movie not found")
		}

		return nil, status.Error(codes.Internal, "failed to update movie info")
	}

	return &pb.UpdateMovieResponse{Movie: toPb(newMovie)}, nil
}

func (srv *server) DeleteMovie(ctx context.Context, in *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	// TODO: add uuid validation
	id := in.GetId()

	ok, err := srv.service.DeleteMovie(id)
	if err != nil && !errors.Is(err, repo.ErrMovieNotExists) {
		return nil, status.Error(codes.Internal, "failed to delete movie info")
	}

	return &pb.DeleteMovieResponse{Success: ok}, nil
}
