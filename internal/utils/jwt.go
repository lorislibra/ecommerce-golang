package utils

import (
	"time"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJwtAccessToken(user *entities.User, refreshTokenId string) (string, time.Time, error) {
	cfg := config.Get()
	currentTime := time.Now()
	expireTime := currentTime.Add(cfg.JwtAccessTokenExpire)

	accessClaim := &models.JwtAccessClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Subject:   user.Oid.Hex(),
			ID:        refreshTokenId,
		},
		Role: user.Role,
	}

	jwtToken := jwt.NewWithClaims(middleware.JwtSignMethod, accessClaim)
	accessToken, err := jwtToken.SignedString(cfg.PrivateKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, expireTime, nil
}

func CreateJwtRefreshToken(user *entities.User, tokenId string) (string, time.Time, error) {
	cfg := config.Get()
	currentTime := time.Now()
	expireTime := currentTime.Add(cfg.JwtRefreshTokenExpire)

	refreshClaim := &models.JwtRefreshClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId,
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Subject:   user.Oid.Hex(),
		},
	}

	jwtToken := jwt.NewWithClaims(middleware.JwtSignMethod, refreshClaim)
	accessToken, err := jwtToken.SignedString(cfg.PrivateKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, expireTime, nil
}

func ParseRefreshToken(refreshToken string) (*models.JwtRefreshClaim, error) {
	claim := new(models.JwtRefreshClaim)

	_, err := jwt.ParseWithClaims(refreshToken, claim, middleware.GetJwtKey)
	if err != nil {
		return nil, err
	}
	return claim, nil
}

func ParseAccessToken(accessToken string) (*models.JwtAccessClaim, error) {
	claim := new(models.JwtAccessClaim)

	_, err := jwt.ParseWithClaims(accessToken, claim, middleware.GetJwtKey)
	if err != nil {
		return nil, err
	}
	return claim, nil
}
