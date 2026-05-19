package grpc

import (
	"context"
	"time"

	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	"github.com/cinema-booking-system/movie-service/internal/domain"
	apperrors "github.com/cinema-booking-system/movie-service/internal/errors"
	"github.com/cinema-booking-system/movie-service/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	moviepb.UnimplementedMovieServiceServer
	movies  *usecase.MovieUseCase
	sessions *usecase.SessionUseCase
}

func NewHandler(movies *usecase.MovieUseCase, sessions *usecase.SessionUseCase) *Handler {
	return &Handler{movies: movies, sessions: sessions}
}

func (h *Handler) CreateMovie(ctx context.Context, req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {
	isAdmin := req.GetRequesterRole() == moviepb.Role_ROLE_ADMIN
	age := int(req.GetAgeRating())
	if age == 0 {
		age = 12
	}
	m := &domain.Movie{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Genre:       req.GetGenre(),
		Duration:    int(req.GetDuration()),
		Rating:      req.GetRating(),
		AgeRating:   age,
	}
	out, err := h.movies.CreateMovie(ctx, isAdmin, m)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.CreateMovieResponse{Movie: toProtoMovie(out)}, nil
}

func (h *Handler) GetMovie(ctx context.Context, req *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {
	out, err := h.movies.GetMovie(ctx, req.GetId())
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.GetMovieResponse{Movie: toProtoMovie(out)}, nil
}

func (h *Handler) UpdateMovie(ctx context.Context, req *moviepb.UpdateMovieRequest) (*moviepb.UpdateMovieResponse, error) {
	isAdmin := req.GetRequesterRole() == moviepb.Role_ROLE_ADMIN
	age := int(req.GetAgeRating())
	if age == 0 {
		age = 12
	}
	m := &domain.Movie{
		ID:          req.GetId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Genre:       req.GetGenre(),
		Duration:    int(req.GetDuration()),
		Rating:      req.GetRating(),
		AgeRating:   age,
	}
	out, err := h.movies.UpdateMovie(ctx, isAdmin, m)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.UpdateMovieResponse{Movie: toProtoMovie(out)}, nil
}

func (h *Handler) DeleteMovie(ctx context.Context, req *moviepb.DeleteMovieRequest) (*moviepb.DeleteMovieResponse, error) {
	isAdmin := req.GetRequesterRole() == moviepb.Role_ROLE_ADMIN
	if err := h.movies.DeleteMovie(ctx, isAdmin, req.GetId()); err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.DeleteMovieResponse{Success: true}, nil
}

func (h *Handler) ListMovies(ctx context.Context, req *moviepb.ListMoviesRequest) (*moviepb.ListMoviesResponse, error) {
	limit := int(req.GetLimit())
	if limit == 0 {
		limit = 20
	}
	var after, before *time.Time
	if ts := req.GetCreatedAfter(); ts != nil {
		t := ts.AsTime().UTC()
		after = &t
	}
	if ts := req.GetCreatedBefore(); ts != nil {
		t := ts.AsTime().UTC()
		before = &t
	}
	list, total, err := h.movies.ListMovies(ctx, limit, int(req.GetOffset()), req.GetGenre(), after, before)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	out := make([]*moviepb.Movie, 0, len(list))
	for _, m := range list {
		out = append(out, toProtoMovie(m))
	}
	return &moviepb.ListMoviesResponse{Movies: out, Total: total}, nil
}

func (h *Handler) SearchMovies(ctx context.Context, req *moviepb.SearchMoviesRequest) (*moviepb.SearchMoviesResponse, error) {
	limit := int(req.GetLimit())
	if limit == 0 {
		limit = 20
	}
	list, total, err := h.movies.SearchMovies(ctx, req.GetTitle(), req.GetGenre(), limit, int(req.GetOffset()))
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	out := make([]*moviepb.Movie, 0, len(list))
	for _, m := range list {
		out = append(out, toProtoMovie(m))
	}
	return &moviepb.SearchMoviesResponse{Movies: out, Total: total}, nil
}

func (h *Handler) CreateHall(ctx context.Context, req *moviepb.CreateHallRequest) (*moviepb.CreateHallResponse, error) {
	isAdmin := req.GetRequesterRole() == moviepb.Role_ROLE_ADMIN
	hall, err := h.sessions.CreateHall(ctx, isAdmin, req.GetName(), int(req.GetCapacity()), req.GetCity(), req.GetCinemaName())
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.CreateHallResponse{Hall: toProtoHall(hall)}, nil
}

func (h *Handler) GetHall(ctx context.Context, req *moviepb.GetHallRequest) (*moviepb.GetHallResponse, error) {
	hall, seats, err := h.sessions.GetHall(ctx, req.GetHallId())
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	ss := make([]*moviepb.Seat, 0, len(seats))
	for _, s := range seats {
		ss = append(ss, toProtoSeat(s))
	}
	return &moviepb.GetHallResponse{Hall: toProtoHall(hall), Seats: ss}, nil
}

func (h *Handler) CreateSession(ctx context.Context, req *moviepb.CreateSessionRequest) (*moviepb.CreateSessionResponse, error) {
	if req.GetStartTime() == nil {
		return nil, apperrors.ToGRPC(domain.ErrInvalidArgument)
	}
	s := &domain.Session{
		MovieID:   req.GetMovieId(),
		HallID:    req.GetHallId(),
		StartTime: req.GetStartTime().AsTime().UTC(),
		Price:     req.GetPrice(),
	}
	out, err := h.sessions.CreateSession(ctx, s)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.CreateSessionResponse{Session: toProtoSession(out)}, nil
}

func (h *Handler) GetSession(ctx context.Context, req *moviepb.GetSessionRequest) (*moviepb.GetSessionResponse, error) {
	out, err := h.sessions.GetSession(ctx, req.GetId())
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &moviepb.GetSessionResponse{Session: toProtoSession(out)}, nil
}

func (h *Handler) ListSessions(ctx context.Context, req *moviepb.ListSessionsRequest) (*moviepb.ListSessionsResponse, error) {
	limit := int(req.GetLimit())
	if limit == 0 {
		limit = 20
	}
	var day *time.Time
	if d := req.GetDate(); d != nil {
		t := d.AsTime().UTC()
		day = &t
	}
	list, total, err := h.sessions.ListSessions(ctx, req.GetMovieId(), req.GetCity(), day, limit, int(req.GetOffset()))
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	out := make([]*moviepb.Session, 0, len(list))
	for _, s := range list {
		out = append(out, toProtoSession(s))
	}
	return &moviepb.ListSessionsResponse{Sessions: out, Total: total}, nil
}

func (h *Handler) GetAvailableSeats(ctx context.Context, req *moviepb.GetAvailableSeatsRequest) (*moviepb.GetAvailableSeatsResponse, error) {
	seats, err := h.sessions.GetAvailableSeats(ctx, req.GetSessionId())
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	out := make([]*moviepb.Seat, 0, len(seats))
	for _, s := range seats {
		out = append(out, toProtoSeat(s))
	}
	return &moviepb.GetAvailableSeatsResponse{Seats: out}, nil
}

func toProtoMovie(m *domain.Movie) *moviepb.Movie {
	return &moviepb.Movie{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Genre:       m.Genre,
		Duration:    int32(m.Duration),
		Rating:      m.Rating,
		AgeRating:   int32(m.AgeRating),
		CreatedAt:   timestamppb.New(m.CreatedAt.UTC()),
	}
}

func toProtoHall(h *domain.Hall) *moviepb.Hall {
	return &moviepb.Hall{
		Id:          h.ID,
		Name:        h.Name,
		Capacity:    int32(h.Capacity),
		City:        h.City,
		CinemaName:  h.CinemaName,
	}
}

func toProtoSeat(s *domain.Seat) *moviepb.Seat {
	return &moviepb.Seat{
		Id:           s.ID,
		HallId:       s.HallID,
		Row:          s.Row,
		Number:       int32(s.Number),
		IsAvailable:  s.IsAvailable,
	}
}

func toProtoSession(s *domain.Session) *moviepb.Session {
	return &moviepb.Session{
		Id:        s.ID,
		MovieId:   s.MovieID,
		HallId:    s.HallID,
		StartTime: timestamppb.New(s.StartTime.UTC()),
		Price:     s.Price,
	}
}
