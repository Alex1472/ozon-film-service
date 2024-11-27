package api

import (
	"github.com/Alex1472/ozon-film-service/internal/model"
	pb "github.com/Alex1472/ozon-film-service/pkg/film-service"
)

func filmToGrpc(film *model.Film) *pb.Film {
	return &pb.Film{
		Id:               film.ID,
		Name:             film.Name,
		Rating:           film.Rating,
		ShortDescription: film.ShortDescription,
	}
}
