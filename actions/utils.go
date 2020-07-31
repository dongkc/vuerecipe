package actions

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// 设置cookie默认到关闭浏览器
func CustomCookieStore(keyPairs ...[]byte) *sessions.CookieStore {
	cs := &sessions.CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 0,
		},
	}

	cs.MaxAge(cs.Options.MaxAge)
	return cs
}

