package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/beevik/guid"
)

var query_strings = struct {
	get_string_albumn    string
	post_string_albumn   string
	get_string_author    string
	post_string_author   string
	update_string_author string
}{
	get_string_albumn: `SELECT  COUNT(songs_views.id) AS songs_count FROM songs_views
INNER JOIN songs ON songs.id = songs_views.song_id
INNER JOIN albumns ON albumns.id = songs.albumn_id
 WHERE user_id=? AND albumns.id=?;`,
	post_string_albumn: `INSERT INTO albumns_views(id, albumn_id, user_id)(?,?,?)`,
	post_string_author: `INSERT INTO authors_views(id, author_id, user_id)(?,?,?)`,
	get_string_author: `SELECT  COUNT(songs_views.id) AS songs_count FROM songs_views
INNER JOIN songs ON songs.id = songs_views.song_id
INNER JOIN authors  ON authors.id = songs.author_id
 WHERE (user_id=? AND author.id=?) `,
	update_string_author: `UPDATE authors_views SET listened =listened+1, last_listened=NOW()
			WHERE author_id=? AND user_id=?;`}

type listenInput struct {
	song_id       string
	albumn_id     string
	main_author   string
	extra_authors []string
}

func (env *Env) AddListened(w http.ResponseWriter, req *http.Request) {
	user_id, is_guest_session := ReturnUserId(req.Context(), w)
	if user_id == "" {
		return
	}
	if is_guest_session {
		if env.GuestListened(w, req) {
			return
		}
	}
	ex_a := CreateExtraArgs(env.DB, req.Context(), w)
	var input listenInput
	json.NewDecoder(req.Body).Decode(input)
	query_rows, error := env.DB.QueryContext(req.Context(), `WITH albumns_cte (albumn_id) AS (
      SELECT albumn_id FROM songs WHERE id=? LIMIT 1
), authors_cte  AS (
      SELECT author_id FROM songs WHERE id=? LIMIT 1
)SELECT EXISTS(SELECT id FROM songs_views WHERE user_id=? AND song_id=?  LIMIT 1)
	 AS is_song_view,  EXISTS(SELECT id FROM albumns_views INNER JOIN albumns_cte ON albumns_views.albumn_id = albumns_cte.albumn_id
        WHERE  albumns_views.user_id=? AND albumns_views.albumn_id=albumns_cte.albumn_id  LIMIT 1 ) AS is_albumn_view,
 EXISTS(SELECT id FROM authors_views INNER JOIN authors_cte ON authors_views.author_id=authors_cte.author_id
  WHERE  authors_views.user_id=? AND authors_views.author_id=authors_cte.author_id LIMIT 1 ) AS is_author_view;`, input.song_id, input.song_id, user_id, input.song_id,
		user_id, user_id)
	if check_bool(error, w) {
		return
	}
	query_rows.Next()
	var is_song_view, is_albumn_view, is_author_view int
	err1 := query_rows.Scan(is_song_view, is_albumn_view, is_author_view)
	query_rows.Close()
	if check_bool(err1, w) {
		return
	}
	tx, err2 := env.DB.BeginTx(req.Context(), nil)
	ex_a.Tx = tx
	if check_bool(err2, w) {
		return
	}

	if is_song_view == 0 {
		if PostQueries(ex_a, `INSERT INTO songs_views (id, song_id, user_id)(?,?,?);`, guid.NewString(), user_id, input.song_id) {
			return
		} else {
			if PostQueries(ex_a, `UPDATE song_views 
SET listened = listened+1, last_listened=NOW()
WHERE song_id=? AND user_id=?;`, input.song_id, user_id) {
				return
			}
		}
		if is_albumn_view == 0 {
			if AlbumnMusicianAdd(query_strings.get_string_albumn, query_strings.post_string_albumn, user_id, input.albumn_id, ex_a) {
				return
			}
		} else {
			if PostQueries(ex_a, ` UPDATE albumns_views SET listened =listened+1, last_listened=NOW()
			WHERE albumn_id=? AND user_id=?;`, input.albumn_id, user_id) {
				return
			}
		}
		if is_author_view == 0 {
			if AlbumnMusicianAdd(query_strings.get_string_author, query_strings.post_string_author, user_id, input.main_author, ex_a) {
				return
			}
		} else {
			if PostQueries(ex_a, query_strings.update_string_author, input.main_author, user_id) {
				return
			}
		}

		for _, value := range input.extra_authors {
			if ExtraAuthors(value, user_id, ex_a) {
				return
			}
		}

	}

}

func AlbumnMusicianAdd(get_query_string string, post_query_string string, song_id string, user_id string, ex_a ExtraArgs) bool {

	albumn_rows, err5 := ex_a.Db.QueryContext(ex_a.Ctx, get_query_string, user_id, song_id)
	if check_bool(err5, ex_a.W) {
		return true
	}
	var albumn_count int
	albumn_rows.Next()
	if albumn_rows.Scan(albumn_count) == nil {
		if albumn_count > 2 {
			if PostQueries(ex_a, post_query_string, user_id, song_id) {
				return true
			}
		}
		return false
	}
	return true
}

func ReturnUserId(ctx context.Context, w http.ResponseWriter, is_session_accepted bool) (string, bool) {
	var id_field = ctx.Value("id_field")
	is_session := false
	if id_field == nil {
		id_field = ctx.Value("session_field")
		if !is_session_accepted {
			w.WriteHeader(403)
			return "", true
		}
		is_session = true
		if id_field == nil {
			w.WriteHeader(403)
			return "", false
		}
	}
	return id_field.(string), is_session

}

func PostQueries(ex_a ExtraArgs, query_string string, arg ...any) bool {
	_, error := ex_a.Tx.ExecContext(ex_a.Ctx, query_string, arg...)
	return check_bool(error, ex_a.W)
}

func ExtraAuthors(authors_id string, user_id string, ex_a ExtraArgs) bool {
	var is_exist, count int
	is_exist_row, error := ex_a.Db.QueryContext(ex_a.Ctx, ` SELECT EXISTS(SELECT id FROM authors_views WHERE author_id=? AND user_id=? LIMIT 1)`, authors_id, user_id)
	if check_bool(error, ex_a.W) {
		return true
	}
	is_exist_row.Next()
	if is_exist_row.Scan(is_exist) == nil {
		if is_exist == 0 {
			counter_row, err1 := ex_a.Db.QueryContext(ex_a.Ctx, `SELECT  COUNT(songs_views.id) AS songs_count FROM songs_views
INNER JOIN songs ON songs.id = songs_views.song_id
INNER JOIN authors_to_songs ON authors_to_songs.songs_id = songs.id
 WHERE songs_views.user_id=? AND authors_to_songs.author_id=?;`, user_id, authors_id)
			if check_bool(err1, ex_a.W) {
				return true
			}
			counter_row.Next()
			if counter_row.Scan(count) == nil {
				if count > 2 {
					_, err2 := ex_a.Tx.ExecContext(ex_a.Ctx, `INSERT INTO authors_views (id, author_id, user_id)(?,?,?)`, guid.NewString(), authors_id, user_id)
					if check_bool(err2, ex_a.W) {
						return true
					} else {

						return false
					}

				}
			}
		} else {
			if PostQueries(ex_a, query_strings.update_string_author, authors_id, user_id) {
				return true
			}
			return false
		}
	}
	return true
}
