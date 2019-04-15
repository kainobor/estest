package user

import (
	"github.com/kainobor/estest/app/lib/logger"
	"github.com/kainobor/estest/app/lib/session"
	"github.com/kainobor/estest/app/pkg/user"
	"github.com/nu7hatch/gouuid"
	"net/http"
)

type User struct {
	userSrv user.Service
}

func New(userSrv user.Service) *User {
	return &User{userSrv: userSrv}
}

func (ctrl *User) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.New(r.Context())

	login := r.URL.Query().Get("login")
	if login == "" {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("Login can't be empty")); err != nil {
			log.Errorw("Failed to write login response", "err", err)
		}
		return
	}
	pass := r.URL.Query().Get("pass")
	if pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("Login can't be empty")); err != nil {
			log.Errorw("Failed to write login response", "err", err)
		}
		return
	}

	id, err := ctrl.userSrv.Auth(r.Context(), login, pass)
	if err != nil {
		log.Errorw("failed to auth user", "login", login, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Failed to load user")); err != nil {
			log.Errorw("Failed to auth by login", "err", err, "login", login)
		}
		return
	}

	sessUUID, _ := uuid.NewV4()
	sessID := sessUUID.String()

	sess, _ := session.Get(r)
	sess.Values["session_id"] = sessID
	sess.Values["id"] = id
	if err := sess.Save(r, w); err != nil {
		log.Errorw("failed to save user session", "login", login, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
