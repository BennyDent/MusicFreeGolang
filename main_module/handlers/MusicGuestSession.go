package handlers

import (
	"net/http"
	"time"

	"github.com/beevik/guid"
)

func (env *Env) SessionAuthorise(w http.ResponseWriter, req *http.Request) {
	var string_id = guid.NewString()
	_, error := env.DB.ExecContext(req.Context(), `INSERT INTO guest_session (session_id)(?)`, string_id)
	if check_bool(error, w) {
		return
	}

	expires := time.Hour * 24 * 7
	cookie := &http.Cookie{
		Name:     "guest_session",
		Value:    string_id,
		HttpOnly: true,
		MaxAge:   int(expires.Seconds()),
	}
	http.SetCookie(w, cookie)
}

func (env *Env) GuestListened(w http.ResponseWriter, req *http.Request) bool {
	ex_a := CreateExtraArgs(env.DB, req.Context(), w)
	cookie, err1 := req.Cookie("guest_session")
	if check_bool(err1, w) {
		return true
	}
	cookie_string := cookie.Value
	result, err1 := env.DB.QueryContext(req.Context(), `
SELECT EXISTS(SELECT id FROM guest_session_view WHERE guest_session_view.session_id=? AND guest_session_view.user_id=? LIMIT 1 );`, cookie_string, req.PathValue("song"))
	if check_auth(err1, w) {
		return true
	}
	var is_listened int
	result.Next()
	result.Scan(is_listened)
	if is_listened == 0 {
		return PostQueries(ex_a, `INSERT INTO guest_session_view (session_id)(?)`, cookie_string)

	} else {
		return PostQueries(ex_a, `UPDATE guest_session_view SET listened = listened+1, last_listened=NOW() WHERE session_id = ? `, cookie_string)
	}
}

func check_auth(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(403)
		return true
	} else {
		return false
	}
}
