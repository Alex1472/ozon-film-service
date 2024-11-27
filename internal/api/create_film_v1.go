package api

import (
	pb "github.com/Alex1472/ozon-film-service/pkg/film-service"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (fa *filmAPI) CreateFilmV1(ctx context.Context, req *pb.CreateFilmV1Request) (*pb.CreateFilmV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateFilmV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := fa.s.Create(ctx, req.GetName(), req.GetRating(), req.GetShortDescription())
	if err != nil {
		log.Error().Err(err).Msg("CreateFilmV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("CreateFilmV1 - success")

	return &pb.CreateFilmV1Response{
		FilmId: id,
	}, nil
}
