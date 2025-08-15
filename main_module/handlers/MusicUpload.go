package handlers

import (
	"encoding/json"

	"mime"
	"mime/multipart"
	"net/http"

	"github.com/beevik/guid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func MusicianCoverUpload(w http.ResponseWriter, req *http.Request) {

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

func AlbumnUpload(w http.ResponseWriter, req *http.Request) {

	client, err1 := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if check_bool(err1, w) {
		return
	}
	defer func() {
		if err := client.Disconnect(req.Context()); err != nil {
			panic(err)
		}
	}()
	var albumn_id string
	var songs_ids = []string{}
	database := client.Database("")
	_, params, err1 := mime.ParseMediaType(req.Header.Get("Content-Type"))
	mr := multipart.NewReader(req.Body, params["boundary"])
	p, err3 := mr.NextPart()
	for err3 != nil && err3.Error() != "io.EOF" {

		if check_bool(err3, w) {
			return
		}
		var bucket = database.GridFSBucket()

		object_ID, err5 := bucket.UploadFromStream(req.Context(), guid.NewString(), p)
		if err5 != nil {
			break
		}
		if p.FileName() == "albumn_cover" {

			albumn_id = object_ID.Hex()
		} else {

			songs_ids = append(songs_ids, object_ID.Hex())
		}
		p, err3 = mr.NextPart()

	}
	result := struct {
		albumn_image string
		songs        []string
	}{
		albumn_image: albumn_id,
		songs:        songs_ids}

	err6 := json.NewEncoder(w).Encode(result)
	if check_bool(err6, w) {
		return
	}
}

type UploadedAuthorInput struct {
	strings map[string]string
}

func (env *Env) UploadMusicianIds(w http.ResponseWriter, req *http.Request) {

	var input UploadedAuthorInput
	err := json.NewDecoder(req.Body).Decode(input)
	if check_bool(err, w) {
		return
	}
	for key, value := range input.strings {

		_, err1 := env.DB.ExecContext(req.Context(), `UPDATE authors
					SET cover_filename=$1
					WHERE authors.id = $2
					`, value, key)

		if check_bool(err1, w) {
			return
		}
	}

}

type UploadedAlbumnInput struct {
	albumn_cover map[string]string
	songs        map[string]string
}

func (env *Env) UploadedAlbumnIds(w http.ResponseWriter, req *http.Request) {

	var input UploadedAlbumnInput
	err := json.NewDecoder(req.Body).Decode(input)
	if check_bool(err, w) {
		return
	}
	var tx, err1 = env.DB.BeginTx(req.Context(), nil)
	if check_bool(err1, w) {
		return
	}
	var cover_filename string
	for key, value := range input.albumn_cover {
		var _, err2 = tx.ExecContext(req.Context(), `UPDATE albumns
					SET cover_filename=$1
					WHERE id = $2
					`, value, key)
		if check_bool(err2, w) {
			return
		}
		cover_filename = value
	}
	for key, value := range input.songs {

		_, err2 := tx.ExecContext(req.Context(), `UPDATE songs
					SET  cover_filename=$1, song_filename=$2
					WHERE id = $3
					`, cover_filename, value, key)
		if check_bool(err2, w) {
			return
		}
	}
	err3 := tx.Commit()
	check_bool(err3, w)
}
