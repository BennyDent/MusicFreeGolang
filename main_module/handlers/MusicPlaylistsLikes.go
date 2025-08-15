package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/beevik/guid"
)

type InputPlaylist struct {
	name        string
	description string
	songs       []string
}

func (env *Env) CreatePlaylist(w http.ResponseWriter, req *http.Request) {
	user_id, is_error := ReturnUserId(req.Context(), w, true)

	if is_error {
		return
	}
	var input_form InputPlaylist
	if check_bool(json.NewDecoder(req.Body).Decode(input_form), w) {
		return
	}
	tx, error := env.DB.BeginTx(req.Context(), nil)
	if check_bool(error, w) {
		return
	}
	playlist_id := guid.NewString()
	_, err := tx.ExecContext(req.Context(), `INSERT INTO playlist (id, name, description, user_id)(?,?,?)`, playlist_id, input_form.name, input_form.description, user_id)

	if check_bool(err, w) {
		return
	}
	for key, value := range input_form.songs {
		_, err := tx.ExecContext(req.Context(), `INSERT INTO playlist_songs (id,playlist_id, song_id,playlist_index)(?,?,?,?)`, guid.NewString(), playlist_id,
			value, key)
		if check_bool(err, w) {
			return
		}
	}
}

func (env *Env) AddToPlaylist(w http.ResponseWriter, req *http.Request) {
	user_id, _ := ReturnUserId(req.Context(), w, false)
	if user_id == "" {
		return
	}
	input_struct := struct {
		song_id     string
		playlist_id string
	}{}

	json.NewDecoder(req.Body).Decode(input_struct)
	result, error := env.DB.QueryContext(req.Context(), `SELECT FROM playlist_songs 
WHERE playlist_id=?
ORDER BY DESC playlist_index
LIMIT 1
`, input_struct.playlist_id)
	if check_bool(error, w) {
		return
	}
	result.Next()
	var playlist_index int
	result.Scan(playlist_index)
	_, err1 := env.DB.ExecContext(req.Context(), `INSERT INTO playlist_songs (id, plaliyst_id,song_id, user_id,
	playlist_index)`, guid.NewString(),
		input_struct.song_id, user_id, playlist_index+1)
	check_bool(err1, w)
}

func (env *Env) DeletePlaylist(w http.ResponseWriter, req *http.Request) {
	user_id, _ := ReturnUserId(req.Context(), w, false)

	var playlist_id string
	json.NewDecoder(req.Body).Decode(playlist_id)
	if user_id == "" {
		return
	}
	_, error := env.DB.ExecContext(req.Context(), `DELETE FROM playlist WHERE playlist.id = ? AND playlist.user_id = ?`,
		playlist_id, user_id)
	check_bool(error, w)
}
