package moviegrpc

import (
	"context"
	"errors"
	"movie-service/internal/model"
	repo "movie-service/internal/repository"
	"movie-service/pkg/pb"

	"github.com/go-playground/validator/v10"
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
	service  Service
	validate *validator.Validate
}

func Register(gRPCServer *grpc.Server, service Service) {
	pb.RegisterMovieServiceServer(gRPCServer, &server{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	})
}

func (srv *server) CreateMovie(ctx context.Context, in *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	newMovie := pbToCreate(in)

	// Create request validation
	if err := srv.validate.Struct(newMovie); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Add info about new movie to repository through the service layer
	newID, err := srv.service.CreateMovie(newMovie.ToModel())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add movie info")
	}

	return &pb.CreateMovieResponse{Id: newID}, nil
}

func (srv *server) GetMovie(ctx context.Context, in *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	id := in.GetId()

	// Check whether it's valid uuid
	if err := srv.validate.Var(id, "uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Get movie info from repository through the service layer
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
	movie := pbToUpdate(in)

	// Update request validation
	if err := srv.validate.Struct(movie); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Update movie info in repository through the service layer
	newMovie, err := srv.service.UpdateMovie(movie.ID, movie.ToModel())
	if err != nil {
		if errors.Is(err, repo.ErrMovieNotExists) {
			return nil, status.Error(codes.NotFound, "movie not found")
		}

		return nil, status.Error(codes.Internal, "failed to update movie info")
	}

	return &pb.UpdateMovieResponse{Movie: toPb(newMovie)}, nil
}

func (srv *server) DeleteMovie(ctx context.Context, in *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	id := in.GetId()

	// Check whether it's valid uuid
	if err := srv.validate.Var(id, "uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Delete movie info from repository through the service layer
	ok, err := srv.service.DeleteMovie(id)
	if err != nil && !errors.Is(err, repo.ErrMovieNotExists) {
		return nil, status.Error(codes.Internal, "failed to delete movie info")
	}

	return &pb.DeleteMovieResponse{Success: ok}, nil
}
