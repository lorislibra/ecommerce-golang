package services

import (
	"context"

	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct {
	userRepo domains.UserRepository
}

func NewUserService(userRepo domains.UserRepository) domains.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) All(ctx context.Context) ([]*dtos.UserMeResp, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, web.ErrInternal
	}

	userDtos := make([]*dtos.UserMeResp, 0, len(users))

	for _, user := range users {
		userDtos = append(userDtos, dtos.UserMeRespFromEntity(user))
	}

	return userDtos, nil
}

func (s *userService) Me(ctx context.Context, oid primitive.ObjectID) (*dtos.UserMeResp, error) {
	user, err := s.userRepo.Find(ctx, oid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, web.ErrNotFound
		}
		return nil, web.ErrInternal
	}

	return dtos.UserMeRespFromEntity(user), nil
}
