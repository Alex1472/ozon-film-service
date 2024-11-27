package api

import (
	"context"
	"errors"
	"github.com/Alex1472/ozon-film-service/internal/service"
	pb "github.com/Alex1472/ozon-film-service/pkg/film-service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (fa *filmAPI) DescribeFilmV1(
	ctx context.Context,
	req *pb.DescribeFilmV1Request,
) (*pb.DescribeFilmV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeFilmV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	film, err := fa.s.Describe(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			log.Debug().Uint64("filmId", req.GetId()).Msg("film not found")
			totalFilmNotFound.Inc()
			return nil, status.Error(codes.NotFound, "film not found")
		}
		log.Error().Err(err).Msg("DescribeFilmV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("DescribeFilmV1 - success")

	return &pb.DescribeFilmV1Response{
		Value: filmToGrpc(film),
	}, nil
}
