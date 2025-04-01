package moviegrpc

import (
	"context"
	"errors"
	"log/slog"
	"movie-service/internal/model"
	repo "movie-service/internal/repository"
	"movie-service/pkg/pb"
	"movie-service/pkg/sl"

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
	l        *slog.Logger
	service  Service
	validate *validator.Validate
}

func Register(gRPCServer *grpc.Server, log *slog.Logger, service Service) {
	pb.RegisterMovieServiceServer(gRPCServer, &server{
		l:        log,
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	})
}

func (srv *server) CreateMovie(ctx context.Context, in *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	const op = "transport.grpc.CreateMovie"

	log := srv.l.With(
		slog.String("op", op),
		slog.Any("request_id", ctx.Value(reqIDKey)),
	)

	newMovie := pbToCreate(in)
	log.Debug("Converted CreateMovieRequest to dto", slog.Any("Request", newMovie))

	// Create request validation
	log.Debug("Validating CreateMovieRequest")
	if err := srv.validate.Struct(newMovie); err != nil {
		log.Error("Validation failed", sl.Err(err))

		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Add info about new movie to repository through the service layer
	log.Debug("Creating movie")
	newID, err := srv.service.CreateMovie(newMovie.ToModel())
	if err != nil {
		log.Error("Failed to create movie", sl.Err(err))

		return nil, status.Error(codes.Internal, "failed to add movie info")
	}

	log.Debug("Successfully created movie info", slog.String("movie_id", newID))

	return &pb.CreateMovieResponse{Id: newID}, nil
}

func (srv *server) GetMovie(ctx context.Context, in *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	const op = "transport.grpc.GetMovie"

	log := srv.l.With(
		slog.String("op", op),
		slog.Any("request_id", ctx.Value(reqIDKey)),
	)

	id := in.GetId()
	log.Debug("Got movie ID", slog.String("ID", id))

	// Check whether it's valid uuid
	log.Debug("Validating movie ID")
	if err := srv.validate.Var(id, "uuid"); err != nil {
		log.Error("Validation failed", sl.Err(err))

		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Get movie info from repository through the service layer
	log.Debug("Getting movie info by ID")
	movie, err := srv.service.GetMovie(id)
	if err != nil {
		log.Error("Failed to get movie info", sl.Err(err))

		if errors.Is(err, repo.ErrMovieNotExists) {
			return nil, status.Error(codes.NotFound, "movie not found")
		}

		return nil, status.Error(codes.Internal, "failed to get movie info")
	}

	log.Debug("Successfully found movie info", slog.Any("Movie", movie))

	return &pb.GetMovieResponse{Movie: toPb(movie)}, nil
}

func (srv *server) UpdateMovie(ctx context.Context, in *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	const op = "transport.grpc.UpdateMovie"

	log := srv.l.With(
		slog.String("op", op),
		slog.Any("request_id", ctx.Value(reqIDKey)),
	)

	movie := pbToUpdate(in)
	log.Debug("Converted UpdateMovieRequest to dto", slog.Any("request", movie))

	// Update request validation
	log.Debug("Validating UpdateMovieRequest")
	if err := srv.validate.Struct(movie); err != nil {
		log.Error("Validation failed", sl.Err(err))

		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Update movie info in repository through the service layer
	log.Debug("Updating movie info")
	newMovie, err := srv.service.UpdateMovie(movie.ID, movie.ToModel())
	if err != nil {
		log.Error("Failed to update movie info", sl.Err(err))

		if errors.Is(err, repo.ErrMovieNotExists) {
			return nil, status.Error(codes.NotFound, "movie not found")
		}

		return nil, status.Error(codes.Internal, "failed to update movie info")
	}

	log.Debug("Successfully updated movie info", slog.Any("New movie", newMovie))

	return &pb.UpdateMovieResponse{Movie: toPb(newMovie)}, nil
}

func (srv *server) DeleteMovie(ctx context.Context, in *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	const op = "transport.grpc.DeleteMovie"

	log := srv.l.With(
		slog.String("op", op),
		slog.Any("request_id", ctx.Value(reqIDKey)),
	)

	id := in.GetId()
	log.Debug("Got movie ID", slog.String("ID", id))

	// Check whether it's valid uuid
	log.Debug("Validating movie ID")
	if err := srv.validate.Var(id, "uuid"); err != nil {
		log.Error("validation failed", sl.Err(err))

		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Delete movie info from repository through the service layer
	log.Debug("Deleting movie info by ID")
	ok, err := srv.service.DeleteMovie(id)
	if err != nil && !errors.Is(err, repo.ErrMovieNotExists) {
		log.Error("Failed to delete movie info", sl.Err(err))

		return nil, status.Error(codes.Internal, "failed to delete movie info")
	}

	if ok {
		log.Debug("Succesfully deleted movie info")
	} else {
		log.Debug("No movie with this ID was found", slog.String("ID", id))
	}

	return &pb.DeleteMovieResponse{Success: ok}, nil
}

func (srv *server) GetMovies(in *pb.GetMoviesRequest, stream pb.MovieService_GetMoviesServer) error {
	const op = "transport.grpc.GetMovies"
	ctx := stream.Context()

	log := srv.l.With(
		slog.String("op", op),
		slog.Any("request_id", ctx.Value(reqIDKey)),
	)

	log.Debug("Getting all movies info from db")
	// TODO: change for dynamic reading (get rid of variant all-movies-in-slice)
	movies, err := srv.service.GetMovies()
	if err != nil {
		log.Error("Failed to get info about all movies", sl.Err(err))

		return nil
	}

	log.Debug("Starting stream...")
	for _, movie := range movies {
		if err := stream.Send(&pb.GetMovieResponse{Movie: toPb(&movie)}); err != nil {
			log.Error("Error during streaming movies")
			return err
		}
	}
	log.Debug("Finished stream")

	return nil
}
