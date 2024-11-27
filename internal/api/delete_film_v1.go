package api

import (
	"errors"
	"github.com/Alex1472/ozon-film-service/internal/service"
	pb "github.com/Alex1472/ozon-film-service/pkg/film-service"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (fa *filmAPI) RemoveFilmV1(ctx context.Context, req *pb.RemoveFilmV1Request) (*pb.RemoveFilmV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveFilmV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := fa.s.Remove(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			log.Debug().Uint64("filmId", req.GetId()).Msg("film not found")
			totalFilmNotFound.Inc()
			return &pb.RemoveFilmV1Response{Found: false}, nil
		}
		log.Error().Err(err).Msg("RemoveFilmV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("RemoveFilmV1 - success")

	return &pb.RemoveFilmV1Response{
		Found: true,
	}, nil
}
