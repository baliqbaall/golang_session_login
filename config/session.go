package config

import "github.com/gorilla/sessions"

const SESSION_ID = "session_login"

var Store = sessions.NewCookieStore([]byte("asd123fghjkl"))
