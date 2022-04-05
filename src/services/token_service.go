package services

import (
	"context"
	"time"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/internal/utils"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tokenService struct {
	refreshTokenRepo domains.RefreshTokenRepository
	userRepo         domains.UserRepository
}

func NewAuthService(refreshTokenRepo domains.RefreshTokenRepository, userRepo domains.UserRepository) domains.AuthService {
	return &tokenService{
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
	}
}

func (s *tokenService) createTokenPair(ctx context.Context, user *entities.User, oldTokenId ...string) (*models.JwtPair, error) {
	refreshTokenId := uuid.New().String()

	refreshToken, err := s.createRefreshToken(ctx, user, refreshTokenId)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.createAccessToken(ctx, user, refreshTokenId)
	if err != nil {
		return nil, err
	}

	if len(oldTokenId) > 0 {
		// save refresh token in redis with TTL and delete the old refresh token
		err = s.refreshTokenRepo.SaveAndDelete(ctx, user.Oid.Hex(), refreshTokenId, refreshToken.Value, time.Until(refreshToken.Expires), oldTokenId[0])
	} else {
		// save refresh token in redis with TTL
		err = s.refreshTokenRepo.Save(ctx, user.Oid.Hex(), refreshTokenId, refreshToken.Value, time.Until(refreshToken.Expires))
	}

	if err != nil {
		return nil, err
	}

	pair := models.JwtPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &pair, nil
}

func (s *tokenService) createRefreshToken(ctx context.Context, user *entities.User, tokenId string) (*fiber.Cookie, error) {
	cfg := config.Get()

	refreshToken, expireTime, err := utils.CreateJwtRefreshToken(user, tokenId)
	if err != nil {
		return nil, err
	}

	maxAge := time.Until(expireTime).Seconds()
	return &fiber.Cookie{
		HTTPOnly: true, // not accesible from javascript
		// Secure:   true, // only http(s) sites
		Name:    cfg.JwtRefreshTokenCookieName,
		Value:   refreshToken,
		Expires: expireTime,
		MaxAge:  int(maxAge),
	}, nil
}

func (s *tokenService) createAccessToken(ctx context.Context, user *entities.User, tokenId string) (*dtos.JwtResponse, error) {
	accessToken, accessTokenExpire, err := utils.CreateJwtAccessToken(user, tokenId)
	if err != nil {
		return nil, err
	}

	expireIn := time.Until(accessTokenExpire).Seconds()
	return &dtos.JwtResponse{
		AccessToken: accessToken,
		ExpireIn:    int64(expireIn),
	}, nil
}

func (s *tokenService) Signin(ctx context.Context, body *dtos.UserSigninBody) (*models.JwtPair, error) {
	user, err := s.userRepo.FindByUserOrEmail(ctx, body.Username, body.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, web.NewWebErrMessage(fiber.StatusUnauthorized, "user not exist")
		}
		return nil, web.ErrInternal
	}

	if !utils.TestPasswordHash(body.Password, user.Password) {
		return nil, web.NewWebErrMessage(fiber.StatusUnauthorized, "bad password")
	}

	pair, err := s.createTokenPair(ctx, user)
	if err != nil {
		return nil, web.ErrInternal
	}

	return pair, nil
}

func (s *tokenService) Signup(ctx context.Context, body *dtos.UserSignupBody) (*models.JwtPair, error) {
	hashPassword, err := utils.GenPasswordHash(body.Password)
	if err != nil {
		return nil, web.ErrInternal
	}

	user := body.ToEntity()
	user.Password = hashPassword
	user.CreatedAt = time.Now()

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, web.NewWebErrMessage(fiber.StatusBadRequest, "user already exist")
		}
		return nil, web.ErrInternal
	}

	pair, err := s.createTokenPair(ctx, user)
	if err != nil {
		return nil, web.ErrInternal
	}

	return pair, nil
}

func (s *tokenService) Signout(ctx context.Context, oid primitive.ObjectID, refreshTokenId string) error {
	var ok bool
	var err error

	if refreshTokenId == "" {
		ok, err = s.refreshTokenRepo.DeleteAll(ctx, oid.Hex())
	} else {
		ok, err = s.refreshTokenRepo.Delete(ctx, oid.Hex(), refreshTokenId)
	}

	if err != nil {
		return web.ErrInternal
	}

	if !ok {
		return web.NewWebErrMessage(fiber.StatusNotFound, "already logged out")
	}

	return nil
}

func (s *tokenService) Refresh(ctx context.Context, refreshTokenRaw string) (*models.JwtPair, error) {
	oldRefreshTokenClaim, err := utils.ParseRefreshToken(refreshTokenRaw)
	if err != nil {
		return nil, web.ErrInternal
	}

	exist, err := s.refreshTokenRepo.Exist(ctx, oldRefreshTokenClaim.Oid().Hex(), oldRefreshTokenClaim.ID)
	if err != nil {
		return nil, web.ErrInternal
	}

	if !exist {
		return nil, web.NewWebErrMessage(fiber.StatusUnauthorized, "refresh token doesn't exists")
	}

	// find user by refresh token oid
	user, err := s.userRepo.Find(ctx, oldRefreshTokenClaim.Oid())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, web.ErrNotFound
		}
		return nil, web.ErrInternal
	}

	// create new token pair and delete the old refresh token
	pair, err := s.createTokenPair(ctx, user, oldRefreshTokenClaim.ID)
	if err != nil {
		return nil, web.ErrInternal
	}

	return pair, nil
}
