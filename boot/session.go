package boot

import (
	"chatchat/app/global"
	"encoding/gob"
	"github.com/gorilla/sessions"
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
