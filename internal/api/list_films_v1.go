package api

import (
	"context"
	pb "github.com/Alex1472/ozon-film-service/pkg/film-service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (fa *filmAPI) ListFilmsV1(
	ctx context.Context,
	_ *pb.ListFilmsV1Request,
) (*pb.ListFilmsV1Response, error) {

	films, err := fa.s.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("ListFilmV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("ListFilmV1 - success")

	var result []*pb.Film
	for _, film := range films {
		result = append(result, filmToGrpc(film))
	}
	return &pb.ListFilmsV1Response{Items: result}, nil
}
