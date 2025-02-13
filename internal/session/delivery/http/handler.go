package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"
)

type SessionHandler struct {
	sessionUC usecase.SessionUsecaseI
}

func NewSessionHandler(sessionUsecase usecase.SessionUsecaseI) *SessionHandler {
	return &SessionHandler{
		sessionUC: sessionUsecase,
	}
}

func (h *SessionHandler) Auth(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	authRequest := &dto.AuthRequest{}
	err = json.Unmarshal(body, authRequest)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	createdSessionEntity, err := h.sessionUC.LoginOrSignup(authRequest)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	authResponse := dto.SessionEntityToResponse(createdSessionEntity)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    createdSessionEntity.JWTAccess,
		Path:     "/",
		Expires:  createdSessionEntity.AccessExpiresAt, /*TODO*/
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    createdSessionEntity.JWTRefresh,
		Path:     "/",
		Expires:  createdSessionEntity.RefreshExpiresAt, /*TODO*/
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.FormatUint(uint64(createdSessionEntity.UserID), 10),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	response, err := json.Marshal(authResponse)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
