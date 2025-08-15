package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/beevik/guid"
)

func (env *Env) SimilarAuthors(w http.ResponseWriter, req *http.Request) {
	var first_query = `SELECT DISTINCT authors.id, authors.name, authors.cover_filename FROM authors 
WHERE are_author_similar(?,authors.id)>? AND are__authors_similar_genre(?,authors.id)> ? LIMIT 15`
	name := req.PathValue("name")

	for i := 0; i < 3; i++ {
		var index int
		switch i {
		case 0:
			index = 5
		case 1:
			index = 3
		case 2:
			index = 3
			first_query = `SELECT DISTINCT authors.id, authors.name, authors.cover_filename FROM authors 
WHERE are__authors_similar_ws_genres(?,authors.id)>? AND are__authors_similar_ws(?,authors.id)> ? LIMIT 15`
		}
		query_rows, error := env.DB.QueryContext(req.Context(), first_query, name, index, name, index)
		if check_bool(error, w) {
			return
		}
		var author_return []AuthorReturnWithFilename
		var is_error bool
		author_return, is_error = convert_to_authors(author_return, query_rows, w)
		if is_error {
			return
		}
		if len(author_return) == 15 {
			break
		}
		err5 := json.NewEncoder(w).Encode(author_return)
		check_bool(err5, w)

	}

}

func ReturnMusicRecomendation(query_string_without_ws string, query_string_ws string) (string, string) {

	select_string := "SELECT songs.id, songs.albumn_id, main_author.id, main_author.name, " +
		"GROUP_CONCAT(extra_authors.id SEPARATOR ',') GROUP_CONCAT(extra_authors.name SEPARATOR ','), songs.song_filename, songs.cover_filename FROM  songs "
	var join_string string
	var second_query_string string = "OR (RAND()*(10-5)+5 AND for_songs_views.user_id=?)"
	var joined_string = "INNER JOIN  authors_to_songs ON songs.id = authors_to_songs.song_id " +
		"INNER JOIN authors AS main_author ON songs.author_id = main_author.id" +
		"INNER JOIN authors AS extra_authors ON authors_to_songs.author_id = extra_authors.id"
	join_string = "INNER JOIN songs_views  AS for_songs_views ON songs.id = songs_views.song_id  INNER JOIN radio_listened AS radio ON songs.id = radio_listened.song_id" + joined_string
	second_query_string += "AND (NOT radio.user_id = ? AND readio.expires < NOW())"

	return select_string + join_string + query_string_ws + second_query_string + "GROUP BY songs.id  ORDER BY COUNT(for_songs_views.id), RAND()*(100)" + "LIMIT ? OFFSET ?;",
		select_string + join_string + query_string_without_ws + second_query_string + "GROUP BY songs.id  ORDER BY COUNT(for_songs_views.id), RAND()*(100)" + "LIMIT ? OFFSET ?;"
}

func (ex_a *ExtraArgs) RecommendationReturn(req *http.Request, search_string string, search_string_ws string, search_args ...any) {
	user_id, _ := ReturnUserId(ex_a.Ctx, ex_a.W, false)
	if user_id == "" {

		return
	}

	var page_size, page_index, is_error = PagesStringToNumbers(req, ex_a.W)
	if is_error {
		return
	}

	search_args = append(search_args, user_id, page_size, page_index*page_size)

	query_string, query_string_ws := ReturnMusicRecomendation(search_string, search_string_ws)

	row_result, error := ex_a.Db.QueryContext(ex_a.Ctx, query_string, search_args...)
	var result_array []SongReturn
	result_array, is_error = RowsToSong(row_result, ex_a.W)
	if is_error {
		return
	}

	var is_similar = false
	if len(result_array) < page_size {
		is_similar = true
		row_result, error = ex_a.Db.QueryContext(ex_a.Ctx, query_string_ws, search_args...)

		if check_bool(error, ex_a.W) {
			return
		}
		new_array, is_err := RowsToSong(row_result, ex_a.W)
		if is_err {
			return
		}
		result_array = append(result_array, new_array...)
		if is_error {
			return
		}
	}
	ex_a.Tx, error = ex_a.Db.BeginTx(req.Context(), nil)
	is_error = ex_a.AddToRadioListened(result_array, user_id)
	PageSearchReturn(result_array, len(result_array) == page_size, ex_a.W, is_similar)
}
func (env *Env) RecomendationAuthor(w http.ResponseWriter, req *http.Request) {

	ex_a := CreateExtraArgs(env.DB, req.Context(), w)
	ex_a.RecommendationReturn(req, "WHERE are_author_similar(?, songs.author_id)>5 AND are__authors_similar_genre(?, songs.author_id)>3 OR are_author_similar(?, authors_to_songs.author_id)>3 AND are__authors_similar_genre(?, authors_to_songs.author_id)>3",
		"WHERE are_author_similar(?, songs.author_id)>5 AND are__authors_similar_genre(?, songs.author_id)>3 OR are_author_similar(?, authors_to_songs.author_id)>3 AND are__authors_similar_genre(?, authors_to_songs.author_id)>3",
		req.PathValue("author_id"), req.PathValue("author_id")) //доделать

}

func (env *Env) RecommendationAlbumn(w http.ResponseWriter, req *http.Request) {
	ex_a := CreateExtraArgs(env.DB, req.Context(), w)
	ex_a.RecommendationReturn(req, "WHERE (are_songs_albumns_tags_similar(?, songs.id) >5 AND are__songs_albumns_genres_similar(?, songs.id)>5)",
		"WHERE (are_songs_albumns_tags_similar_ws(?, songs.id) >5 AND are__songs_albumns_genres_similar_ws(?, songs.id)>5)", req.PathValue("song_id"), req.PathValue("song_id"))

}

func convert_to_authors(author_return []AuthorReturnWithFilename, query_rows *sql.Rows, w http.ResponseWriter) ([]AuthorReturnWithFilename, bool) {
	for query_rows.Next() {
		var author AuthorReturnWithFilename
		err := query_rows.Scan(author.id, author.Name, author.cover_filename)
		if check_bool(err, w) {
			return author_return, true
		}
		author_return = append(author_return, author)
	}
	return author_return, false
}

func (ex_a *ExtraArgs) AddToRadioListened(data []SongReturn, user_id string) bool {

	query_string := "INSERT INTO radio_listened(id, user_id, song_id, expires)(?,?,?);"
	expires := time.Now().Add(time.Hour * 24 * 2)
	for _, value := range data {
		_, error := ex_a.Tx.ExecContext(ex_a.Ctx, query_string, value, guid.NewString(), user_id, value, expires)
		if check_bool(error, ex_a.W) {
			return true
		}
	}

	error := ex_a.Tx.Commit()
	if check_bool(error, ex_a.W) {
		return true
	}

	return false
}
