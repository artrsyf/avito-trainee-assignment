package dto

import (
	"os"
	"strconv"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
	"github.com/dgrijalva/jwt-go"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func createJWT(authRequest *AuthRequest, ttl time.Time, userID uint) (string, error) {
	jwtTokenKey := []byte(os.Getenv("TOKEN_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{
			"username": authRequest.Username,
			"id":       strconv.FormatUint(uint64(userID), 10),
		},
		"iat": time.Now().Unix(),
		"exp": ttl,
	})
	tokenString, err := token.SignedString(jwtTokenKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createSignedSession(authRequest *AuthRequest, userID uint, accessTokenTTL time.Time, refreshTokenTTL time.Time) (*entity.Session, error) {
	accessToken, err := createJWT(authRequest, accessTokenTTL, userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := createJWT(authRequest, refreshTokenTTL, userID)
	if err != nil {
		return nil, err
	}

	return &entity.Session{
		JWTAccess:        accessToken,
		JWTRefresh:       refreshToken,
		UserID:           userID,
		Username:         authRequest.Username,
		AccessExpiresAt:  accessTokenTTL,
		RefreshExpiresAt: refreshTokenTTL,
	}, nil
}

func SessionEntityToModel(sessionEntity *entity.Session) *model.Session {
	return &model.Session{
		JWTAccess:        sessionEntity.JWTAccess,
		JWTRefresh:       sessionEntity.JWTRefresh,
		UserID:           sessionEntity.UserID,
		Username:         sessionEntity.Username,
		AccessExpiresAt:  sessionEntity.AccessExpiresAt,
		RefreshExpiresAt: sessionEntity.RefreshExpiresAt,
	}
}

func SessionModelToEntity(sessionModel *model.Session) *entity.Session {
	return &entity.Session{
		JWTAccess:        sessionModel.JWTAccess,
		JWTRefresh:       sessionModel.JWTRefresh,
		UserID:           sessionModel.UserID,
		Username:         sessionModel.Username,
		AccessExpiresAt:  sessionModel.AccessExpiresAt,
		RefreshExpiresAt: sessionModel.RefreshExpiresAt,
	}
}

func AuthRequestToEntity(authRequest *AuthRequest, userID uint, accessTokenTTL time.Time, refreshTokenTTL time.Time) (*entity.Session, error) {
	signedSessionEntity, err := createSignedSession(authRequest, userID, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return signedSessionEntity, nil
}

func SessionEntityToResponse(sessionEntity *entity.Session) *AuthResponse {
	return &AuthResponse{
		Token: sessionEntity.JWTAccess,
	}
}
