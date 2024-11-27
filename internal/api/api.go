package api

import (
	"context"
	"github.com/Alex1472/ozon-film-service/internal/model"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	pb "github.com/Alex1472/ozon-film-service/pkg/film-service"
)

var (
	totalFilmNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "film_service_api_template_not_found_total",
		Help: "Total number of templates that were not found",
	})
)

type Service interface {
	List(ctx context.Context) ([]*model.Film, error)
	Describe(ctx context.Context, filmID uint64) (*model.Film, error)
	Create(ctx context.Context, name string, rating float64, shortDescription string) (uint64, error)
	Remove(ctx context.Context, filmID uint64) error
}

type filmAPI struct {
	pb.UnimplementedFilmServiceServer
	s Service
}

func NewFilmAPI(service Service) pb.FilmServiceServer {
	return &filmAPI{s: service}
}
