package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/beevik/guid"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func check_bool(err error, w http.ResponseWriter) bool {
	if err != nil {

		w.WriteHeader(500)
		return true
	}
	return false
}

type AuthorInput struct {
	name string
}

func MusicianCreate(db *sql.DB) http.Handler {

	var cover_filename_return string = guid.NewString()
	fn := func(w http.ResponseWriter, req *http.Request) {

		var input AuthorInput

		var error = json.NewDecoder(req.Body).Decode(input)

		_, error = db.ExecContext(req.Context(), `INSERT INTO authors (id, name, cover_filename) VALUES ($1, $2, $3)`,
			guid.NewString(), input.name, cover_filename_return)
		if check_bool(error, w) {
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(cover_filename_return))
	}
	return http.HandlerFunc(fn)
}

type GenreTagInput struct {
	name   string
	is_tag bool
}

func TagorGenreCreate(db *sql.DB) http.Handler {

	fn := func(w http.ResponseWriter, req *http.Request) {
		var input GenreTagInput

		var table string
		if input.is_tag {
			table = `tags`
		} else {
			table = `genres`
		}
		var error = json.NewDecoder(req.Body).Decode(input)
		_, error = db.ExecContext(req.Context(), `INSERT INTO`+table+`( name) VALUES ($1, )`,
			input.name)
		if check_bool(error, w) {
			return
		}

	}
	return http.HandlerFunc(fn)
}

/*
interface AlbumnSendInterface{
name: string,

	        cover_filename: string
		    main_author?: string,
		    extra_authors: Array<string>,
		    songs: Array<SendSongInterface>,
		    tags: Array<String>,
		    genres: Array<String>,
		    date: string
		}
*/

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

func AlbumnCreate(db *sql.DB) http.Handler {

	fn := func(w http.ResponseWriter, req *http.Request) {

		var albumn_id = guid.NewString()
		var tx, error = db.BeginTx(req.Context(), nil)

		var input AlbumnInput
		error = json.NewDecoder(req.Body).Decode(input)
		if check_bool(error, w) {
			return
		}
		var cover_filename = guid.NewString()
		_, error = tx.ExecContext(req.Context(), `INSERT INTO albumns 
        (id, name, author_id, albumns_date, cover_filename, albumn_type) VALUES ($1, $2, $3, $4,$5, $6)`,
			albumn_id, input.name, input.main_author, input.albumn_date, cover_filename, input.albumn_type)
		if check_bool(error, w) {
			return
		}

		for _, value := range input.extra_authors {
			_, error = tx.ExecContext(req.Context(), `INSERT INTO authors_to_albumns (id, author_id, albumns_id) VALUES ($1, $2, $3)`,
				guid.NewString(),
				value, albumn_id)
			if check_bool(error, w) {
				return
			}

		}
		tags_bool := tags_genres_make(input.tags, w, req.Context(), tx, `INSERT INTO 
            tags_to_albumns(id, tag_name, albumn_id) VALUES ($1,$2,$3)`, albumn_id)

		if tags_bool {
			return
		}

		genres_bool := tags_genres_make(input.genres, w, req.Context(), tx, `INSERT INTO 
            genres_to_albumns(id, genre_name, albumn_id) VALUES ($1,$2,$3)`, albumn_id)

		if genres_bool {
			return
		}
		var is_error, songs = song_create(tx, w, error, cover_filename, albumn_id, req.Context(), input.songs)
		if is_error {
			return
		}
		error = tx.Commit()
		if check_bool(error, w) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			songs          []string
			cover_filename string
		}{songs: songs, cover_filename: cover_filename})
	}

	return http.HandlerFunc(fn)
}

func tags_create() {

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
		result_array = append(result_array, song_filename)
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

func MusicianCoverUpload(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		client, err1 := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
		if check_bool(err1, w) {
			return
		}
		defer func() {
			if err := client.Disconnect(req.Context()); err != nil {
				panic(err)
			}
		}()
		database := client.Database("")

		var bucket = database.GridFSBucket()

		err3 := req.ParseMultipartForm(10 << 20)
		if check_bool(err3, w) {
			return
		}
		var object_ID bson.ObjectID
		for new_name, files := range req.MultipartForm.File {

			var file, err5 = files[0].Open()
			if check_bool(err5, w) {
				return
			}
			object_ID, err5 = bucket.UploadFromStream(req.Context(), new_name, file)

			if check_bool(err5, w) {
				return
			}

		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(object_ID[:])
	}

	return http.HandlerFunc(fn)
}

func AlbumnUpload(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		client, err1 := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
		if check_bool(err1, w) {
			return
		}
		defer func() {
			if err := client.Disconnect(req.Context()); err != nil {
				panic(err)
			}
		}()
        var albumn_name: string;
		database := client.Database("")
		
		_, params, err1 := mime.ParseMediaType(req.Header.Get("Content-Type"))
		mr := multipart.NewReader(req.Body, params["boundary"])
		for {
			p, err3 := mr.NextPart()

			if check_bool(err3, w) {
				return
			}
           
            var bucket = database.GridFSBucket();

         object_ID, err5 =  bucket.UploadFromStream(req.Context(), guid.NewString(), p);

            if(p.FileName()=="albumn_cover"){

            }

		}
	}

	return http.HandlerFunc(fn)
}
