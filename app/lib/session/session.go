package session

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

const CookieName = "SESSION_ID"

var store sessions.Store

func Init() {
	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
}

func Get(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, CookieName)
}

func IsAuth(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", nil
	}

	sess, err := store.Get(r, CookieName)
	if err != nil {
		return "", fmt.Errorf("session %s:%v not found: %v", CookieName, cookie.Value, err)
	}

	id, ok := sess.Values["id"]
	if !ok {
		return "", fmt.Errorf("wrong user ID in session: %#v", sess.Values["id"])
	}

	return id.(string), nil
}
