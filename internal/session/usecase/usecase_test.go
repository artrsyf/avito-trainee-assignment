package usecase

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/artrsyf/avito-trainee-assignment/config"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionModel "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
	mockSession "github.com/artrsyf/avito-trainee-assignment/internal/session/repository/mock_repository"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	userModel "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	mockUser "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/mock_repository"
)

func TestSessionUsecase_LoginOrSignup(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mockSession.NewMockSessionRepositoryI(ctrl)
	mockUserRepo := mockUser.NewMockUserRepositoryI(ctrl)

	cfg := config.UserConfig{
		InitCoinsBalance: 100,
		Auth: config.AuthConfig{
			AccessTokenExpiration:  "1h",
			RefreshTokenExpiration: "24h",
		},
	}

	uc := NewSessionUsecase(mockSessionRepo, mockUserRepo, cfg)

	ctx := context.Background()
	testAuthRequest := &dto.AuthRequest{
		Username: "testuser",
		Password: "testpass",
	}

	// Set JWT secret for testing
	os.Setenv("TOKEN_KEY", "test-secret-key")

	t.Run("successful login with existing session", func(t *testing.T) {
		user := &userModel.User{ID: 1, Username: "testuser", PasswordHash: hashPassword("testpass")}
		session := &sessionModel.Session{UserID: 1}

		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(user, nil)
		mockSessionRepo.EXPECT().Check(ctx, uint(1)).Return(session, nil)

		result, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.UserID != 1 {
			t.Errorf("expected userID 1, got %d", result.UserID)
		}
	})

	t.Run("successful login without existing session", func(t *testing.T) {
		user := &userModel.User{ID: 1, Username: "testuser", PasswordHash: hashPassword("testpass")}

		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(user, nil)
		mockSessionRepo.EXPECT().Check(ctx, uint(1)).Return(nil, sessionEntity.ErrNoSession)
		mockSessionRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, s *sessionEntity.Session) (*sessionModel.Session, error) {
			return dto.SessionEntityToModel(s), nil
		})

		result, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.UserID != 1 {
			t.Errorf("expected userID 1, got %d", result.UserID)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		user := &userModel.User{ID: 1, Username: "testuser", PasswordHash: hashPassword("wrongpass")}

		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(user, nil)

		_, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if !errors.Is(err, sessionEntity.ErrWrongCredentials) {
			t.Errorf("expected ErrWrongCredentials, got %v", err)
		}
	})

	t.Run("successful signup new user", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(nil, userEntity.ErrIsNotExist)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(&userModel.User{ID: 2}, nil)
		mockSessionRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, s *sessionEntity.Session) (*sessionModel.Session, error) {
			return dto.SessionEntityToModel(s), nil
		})

		result, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.UserID != 2 {
			t.Errorf("expected new userID 2, got %d", result.UserID)
		}
	})

	t.Run("user repo error", func(t *testing.T) {
		testErr := errors.New("database error")
		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(nil, testErr)

		_, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if !errors.Is(err, testErr) {
			t.Errorf("expected error %v, got %v", testErr, err)
		}
	})

	t.Run("session create error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(nil, userEntity.ErrIsNotExist)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(&userModel.User{ID: 3}, nil)
		mockSessionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("session error"))

		_, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("token generation error", func(t *testing.T) {
		os.Unsetenv("TOKEN_KEY") // Cause JWT generation error
		defer os.Setenv("TOKEN_KEY", "test-secret-key")

		mockUserRepo.EXPECT().GetByUsername(ctx, "testuser").Return(nil, userEntity.ErrIsNotExist)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(&userModel.User{ID: 4}, nil)
		mockSessionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("session error"))

		_, err := uc.LoginOrSignup(ctx, testAuthRequest)
		if err == nil {
			t.Error("expected JWT error but got nil")
		}
	})
}

func hashPassword(pass string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash)
}
