package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func PagesStringToNumbers(req *http.Request, w http.ResponseWriter) (int, int, bool) {

	var page_index_int, err = strconv.Atoi(req.PathValue("page_index"))
	if check_bool(err, w) {
		return 0, 0, false
	}
	var page_size_int, err1 = strconv.Atoi(req.PathValue("page_size"))
	if check_bool(err1, w) {
		return 0, 0, false
	}

	return page_index_int, page_size_int, true

}

func Musician_For_Search(db *sql.DB) http.Handler {

	fn := func(w http.ResponseWriter, req *http.Request) {

		var page_index, page_size, is_error = PagesStringToNumbers(req, w)
		if !is_error {
			return
		}
		var result, error = db.QueryContext(req.Context(), `SELECT authors.id, authors.name, authors.cover_filename 
FROM authors
WHERE name=$1
LIMIT $2 OFFSET $3;`, req.PathValue("name"), page_size, page_index*page_size)
		if check_bool(error, w) {
			return
		}

		var data = []AuthorReturnWithFilename{}
		var page_counter = 1
		for result.Next() {

			new := AuthorReturnWithFilename{}
			error = result.Scan(&new.id, &new.Name, &new.cover_filename)
			if check_bool(error, w) {
				return
			}
			data = append(data, new)
			page_counter++

		}

		if PageSearchReturn(data, page_size, page_counter, w) {
			return
		}
	}

	return http.HandlerFunc(fn)
}

func Tags_Genres__For_Search(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var genre_name = req.PathValue("type")

		var page_index, page_size, is_error = PagesStringToNumbers(req, w)
		if !is_error {
			return
		}

		var result, error = db.QueryContext(req.Context(), `SELECT  $1s.name
FROM $1s
INNER JOIN $1s_to_albumns ON $1s.name=  $1s_to_albumns.$1_name 
INNER JOIN $1s_to_songs ON $1s.name = $1s_to_songs.$1_name
INNER JOIN albumns ON $1s_to_albumns.albumn_id = albumns.id
INNER JOIN albumns_views ON albumns.id = albumns_views.albumn_id
WHERE $1s.name='$2'
GROUP BY $1s.name
ORDER BY COUNT($1s_to_albumns.id), COUNT($1s_to_songs.id), COUNT(albumns_views.id)
LIMIT $3 OFFSET $4;`, genre_name, req.PathValue("name"), page_size, page_index*page_size)
		if check_bool(error, w) {
			return
		}

		var data = []GenreTagsReturn{}
		var page_counter = 1
		for result.Next() {

			new := GenreTagsReturn{}
			error = result.Scan(&new.Name)
			if check_bool(error, w) {
				return
			}
			data = append(data, new)
			page_counter++

		}
		if PageSearchReturn(data, page_size, page_counter, w) {
			return
		}
	}

	return http.HandlerFunc(fn)
}

func PageSearchReturn[T any](data []T, page_size int, page_counter int, w http.ResponseWriter) bool {
	var isMore = true
	if page_size > page_counter {
		isMore = false
	}
	var for_return = struct {
		hasMore bool
		data    []T
	}{
		hasMore: isMore,
		data:    data,
	}

	err := json.NewEncoder(w).Encode(for_return)
	if check_bool(err, w) {
		return true
	} else {
		w.Header().Set("Content-Type", "application/json")

	}
	return false
}
