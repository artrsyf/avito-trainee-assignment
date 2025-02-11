package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"
)

type SessionHandler struct {
	sessionUC usecase.SessionUsecaseI
}

func (h *SessionHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		/*Handle*/
	}

	authRequest := &dto.AuthRequest{}
	err = json.Unmarshal(body, authRequest)
	if err != nil {
		/*Handle*/
	}

	createdSessionEntity, err := h.sessionUC.Signup(authRequest)
	if err != nil {
		/*Handle*/
	}

	authResponse := dto.SessionEntityToResponse(createdSessionEntity)

	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   createdSessionEntity.JWTAccess,
		Path:    "/",
		Expires: time.Now().Add(15 * time.Minute), /*TODO*/
		Secure:  false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   createdSessionEntity.JWTRefresh,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), /*TODO*/
		Secure:  false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  strconv.FormatUint(uint64(createdSessionEntity.UserID), 10),
		Path:   "/",
		Secure: false,
	})

	response, err := json.Marshal(authResponse)
	if err != nil {
		/*Handle*/
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
