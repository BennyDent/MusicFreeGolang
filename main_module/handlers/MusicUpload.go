package handlers

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"net/http"

	"mime"
	"mime/multipart"

	"github.com/beevik/guid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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
		var albumn_id string
		songs := make(map[int]string)
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
				var int, err6 = strconv.Atoi(p.FileName())
				if check_bool(err6, w) {
					return
				}
				songs[int] = object_ID.Hex()
			}
			p, err3 = mr.NextPart()

		}
		result := struct {
			albumn_image string
			songs        map[int]string
		}{
			albumn_image: albumn_id,
			songs:        songs}

		err6 := json.NewEncoder(w).Encode(result)
		if check_bool(err6, w) {
			return
		}
	}

	return http.HandlerFunc(fn)
}

type UploadedALbumnsInput struct {
	strings map[string]map[string]string
}

func UploadedIds(db *sql.DB) http.Handler {

	fn := func(w http.ResponseWriter, req *http.Request) {
		var input UploadedALbumnsInput
		err := json.NewDecoder(req.Body).Decode(input)
		if check_bool(err, w) {
			return
		}
		for key, value := range input.strings {

		}
	}

	return http.HandlerFunc(fn)
}
