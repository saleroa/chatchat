package middleware

import (
	"chatchat/app/global"
	"encoding/gob"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
)

var store *sessions.CookieStore

func SessionSetup() {
	gob.Register(url.Values{})
	session := global.Config.Session

	store = sessions.NewCookieStore([]byte(session.SecretKey))
	store.Options = &sessions.Options{
		Path: "/",
		// session 有效期
		// 单位秒
		MaxAge:   session.MaxAge,
		HttpOnly: true,
	}
}

func Get(r *http.Request, name string) (val interface{}, err error) {
	session1 := global.Config.Session
	// Get a session.
	session, err := store.Get(r, session1.Name)
	if err != nil {
		return
	}

	val = session.Values[name]

	return
}

func Set(w http.ResponseWriter, r *http.Request, name string, val interface{}) (err error) {
	// Get a session.
	session1 := global.Config.Session
	session, err := store.Get(r, session1.Name)
	if err != nil {
		return
	}

	session.Values[name] = val
	err = session.Save(r, w)

	return
}

func Delete(w http.ResponseWriter, r *http.Request, name string) (err error) {
	// Get a session.
	session1 := global.Config.Session
	session, err := store.Get(r, session1.Name)
	if err != nil {
		return
	}

	delete(session.Values, name)
	err = session.Save(r, w)

	return
}
