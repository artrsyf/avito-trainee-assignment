package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	sessionUC usecase.SessionUsecaseI
	logger    *logrus.Logger
}

func NewSessionHandler(sessionUsecase usecase.SessionUsecaseI, logger *logrus.Logger) *SessionHandler {
	return &SessionHandler{
		sessionUC: sessionUsecase,
		logger:    logger,
	}
}

func (h *SessionHandler) Auth(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Incoming Auth request")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"errors": "bad request"})
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.WithError(err).Warn("Failed to close request body")
		}
	}()

	authRequest := &dto.AuthRequest{}
	if err = json.Unmarshal(body, authRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"errors": "bad request"})
		return
	}

	createdSessionEntity, err := h.sessionUC.LoginOrSignup(ctx, authRequest)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"stack": string(debug.Stack()),
		}).Debug("LoginOrSignup error handling")

		switch err {
		case sessionEntity.ErrWrongCredentials:
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "wrong credentials"})
		case userEntity.ErrAlreadyCreated:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"errors": "user conflict"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
		}
		return
	}

	setSessionCookies(w, createdSessionEntity, h.logger)

	response, err := json.Marshal(dto.SessionEntityToResponse(createdSessionEntity))
	if err != nil {
		h.logger.WithError(err).Error("Failed to marshal auth response")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		h.logger.WithError(err).Error("Failed to write auth response")
	}
}

func setSessionCookies(w http.ResponseWriter, session *sessionEntity.Session, logger *logrus.Logger) {
	cookies := []*http.Cookie{
		{
			Name:     "access_token",
			Value:    session.JWTAccess,
			Path:     "/",
			Expires:  session.AccessExpiresAt,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
		{
			Name:     "refresh_token",
			Value:    session.JWTRefresh,
			Path:     "/",
			Expires:  session.RefreshExpiresAt,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
		{
			Name:     "user_id",
			Value:    strconv.FormatUint(uint64(session.UserID), 10),
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	}

	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
		if cookie.Expires.Before(time.Now()) {
			logger.WithFields(logrus.Fields{
				"cookie": cookie.Name,
				"expiry": cookie.Expires,
			}).Warn("Setting expired cookie")
		}
	}
}
