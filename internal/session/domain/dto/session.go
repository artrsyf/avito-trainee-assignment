package dto

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (req *AuthRequest) ValidateAuthRequest(validate *validator.Validate) error {
	err := validate.Struct(req)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				field := err.Field()

				switch err.Tag() {
				case "required":
					return errors.New(field + " is required")
				case "min":
					return errors.New(field + " is too short")
				case "max":
					return errors.New(field + " is too long")
				default:
					return errors.New(field + " is invalid")
				}
			}
		}

		return err
	}

	return nil
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
