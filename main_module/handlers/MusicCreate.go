package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beevik/guid"
)

type Env struct {
	DB *sql.DB
}

func check_bool(err error, w http.ResponseWriter) bool {
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(500)
		return true
	}
	return false
}

type AuthorInput struct {
	name string
}

type GenreTagSimilarInput struct {
	first_tag  string
	second_tag string
	is_tag     bool
}

func (env *Env) SimilarTagsGenres(w http.ResponseWriter, req *http.Request) {
	var input GenreTagSimilarInput
	input_error := json.NewDecoder(req.Body).Decode(input)
	if check_bool(input_error, w) {
		return
	}

	var query_string string

	if input.is_tag {
		query_string = `INSERT INTO similar_tags (id, first_tag, second_tag)(?,?,?);
`
	} else {
		query_string = `INSERT INTO similar_genres (id, first_genre, second_genre)(?,?,?);`
	}

	_, error := env.DB.ExecContext(req.Context(), query_string, guid.NewString(), input.first_tag, input.second_tag)
	check_bool(error, w)
}

func (env *Env) MusicianCreate(w http.ResponseWriter, req *http.Request) {
	var id = guid.NewString()
	var input AuthorInput

	var error = json.NewDecoder(req.Body).Decode(input)

	_, error = env.DB.ExecContext(req.Context(), `INSERT INTO authors (id, name) VALUES (?, ?);`,
		id, input.name)
	if check_bool(error, w) {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(id))

}

func SetCorsHeader(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
}

type AlbumnInput struct {
	name          string
	main_author   string
	extra_authors []string
	albumn_date   string
	albumn_type   string
	tags          []string
	genres        []string
	songs         []SongInput
}

type SongInput struct {
	name          string
	author_id     string
	song_index    int
	tags          []string
	genres        []string
	extra_authors []string
}

func (env *Env) AlbumnCreate(w http.ResponseWriter, req *http.Request) {

	var albumn_id = guid.NewString()
	var tx, error = env.DB.BeginTx(req.Context(), nil)

	var input AlbumnInput
	error = json.NewDecoder(req.Body).Decode(input)
	if check_bool(error, w) {
		return
	}
	var cover_filename = guid.NewString()
	_, error = tx.ExecContext(req.Context(), `INSERT INTO albumns 
        (id, name, author_id, albumns_date, cover_filename, albumn_type) VALUES (?, ?, ?, ?,?, ?);`,
		albumn_id, input.name, input.main_author, input.albumn_date, cover_filename, input.albumn_type)
	if check_bool(error, w) {
		return
	}
	type Result struct {
		albumn_id string
		songs     []string
	}
	var result_struct Result
	for _, value := range input.extra_authors {
		_, error = tx.ExecContext(req.Context(), `INSERT INTO authors_to_albumns (id, author_id, albumns_id) VALUES (?, ?, ?)`,
			guid.NewString(),
			value, albumn_id)
		if check_bool(error, w) {

		}

	}
	tags_bool := tags_genres_make(input.tags, w, req.Context(), tx, `INSERT INTO 
            tags_to_albumns(id, tag_name, albumn_id) VALUES ($1,$2,$3)`, albumn_id)

	if tags_bool {

	}

	genres_bool := tags_genres_make(input.genres, w, req.Context(), tx, `INSERT INTO 
            genres_to_albumns(id, genre_name, albumn_id) VALUES ($1,$2,$3)`, albumn_id)

	if genres_bool {

	}
	var is_error bool
	is_error, result_struct.songs = song_create(tx, w, error, cover_filename, albumn_id, req.Context(), input.songs)
	if is_error {

	}
	error = tx.Commit()
	if check_bool(error, w) {

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result_struct)

}

type GenreTagInput struct {
	genre_string string
	name         string
}

func (env *Env) GenreTagCreation(w http.ResponseWriter, req *http.Request) {
	var input GenreTagInput
	err := json.NewDecoder(req.Body).Decode(input)
	if check_bool(err, w) {
		return
	}
	_, err1 := env.DB.ExecContext(req.Context(), `INSERT INTO $1(name)VALUES($2)`, input.genre_string, input.name)
	check_bool(err1, w)

}

func song_create(tx *sql.Tx, w http.ResponseWriter, err error, cover_filename string, albumn_id string, ctx context.Context, songs []SongInput) (bool, []string) {

	var result_array []string
	for _, value := range songs {
		song_filename := guid.NewString()
		song_id := guid.NewString()
		_, err = tx.ExecContext(ctx, `INSERT INTO songs (id, name, author_id, albumn_id, song_index, cover_filename, song_filename)
         VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			song_id, value.name, value.author_id, albumn_id, value.song_index, cover_filename, song_filename)
		if check_bool(err, w) {
			return true, result_array
		}
		result_array = append(result_array, song_id)
		for _, author := range value.extra_authors {
			_, err = tx.ExecContext(ctx, `INSERT INTO authors_to_songs (id, author_id, songs_id)
             VALUES($1,$2,$3)`, guid.NewString(), author, song_id)
			if check_bool(err, w) {
				return true, result_array
			}
		}
		tags_result := tags_genres_make(value.tags, w, ctx, tx, `INSERT INTO 
            tags_to_songs(id, tag_name, song_id) VALUES ($1,$2,$3)`, song_id)
		if tags_result {
			return true, result_array
		}
		genres_result := tags_genres_make(value.genres, w, ctx, tx, `INSERT INTO 
            genres_to_songs(id, genre_name, song_id) VALUES ($1,$2,$3)`, song_id)
		if genres_result {
			return true, result_array
		}
	}
	return false, result_array
}

func tags_genres_make(tags []string, w http.ResponseWriter, ctx context.Context, tx *sql.Tx, input string,
	first_string string) bool {
	for _, tag := range tags {

		_, error := tx.ExecContext(ctx, input, guid.NewString(), tag, first_string)
		if check_bool(error, w) {
			return true
		}
	}
	return false
}
