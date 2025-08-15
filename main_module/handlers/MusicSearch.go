package handlers

import (
	"database/sql"
	"net/http"
	"strings"
)

func (env *Env) SearchSong(w http.ResponseWriter, req *http.Request) {
	page_index, page_size, is_err := PagesStringToNumbers(req, w)
	if is_err {
		return
	}
	songs_result, error := env.DB.QueryContext(req.Context(), `SELECT songs.id, songs.name ,songs.albumn_id, songs.author_id,authors.name GROUP_CONCAT(authors_to_songs.author_id SEPARATOR ',') AS extra_authors_id,
GROUP_CONCAT(extra_authors.name SEPARATOR ',') AS extra_authors_name  FROM songs
INNER JOIN authors_to_songs ON songs.id = authors_to_songs.songs_id
INNER JOIN songs_views ON songs.id = songs_views.song_id
INNER JOIN authors  ON songs.author_id=authors.id 
INNER JOIN authors AS extra_authors ON authors_to_songs.author_id=extra_authors.id
WHERE songs.name LIKE '%?%'
GROUP BY songs.id
ORDER BY COUNT(songs_views.song_id)
LIMIT ? OFFSET ?`, req.PathValue("name"), page_size, page_index*page_size)
	if check_bool(error, w) {
		return
	}

	var counter = 0
	var result_slice []SongReturn

	for songs_result.Next() {
		var song_struct SongReturn
		var id_strings string
		var name_string string
		songs_result.Scan(song_struct.id, song_struct.Name, song_struct.albumn_id, song_struct.Main_author.id, song_struct.Main_author.Name, id_strings, name_string)
		song_struct.Extra_Authors = extra_authors_return(id_strings, name_string)
		result_slice = append(result_slice, song_struct)
		counter++
	}

	PageSearchReturn(result_slice, page_size == counter, w)
}

func extra_authors_return(id_string string, name_string string) []AuthorReturn {
	ids := strings.Split(id_string, ",")
	names := strings.Split(name_string, ",")
	var for_return []AuthorReturn
	for key, value := range ids {
		var author_return = AuthorReturn{id: value, Name: names[key]}
		for_return = append(for_return, author_return)
	}
	return for_return
}

func (env *Env) AuthorsSearch(w http.ResponseWriter, req *http.Request) {

	page_index, page_size, is_err := PagesStringToNumbers(req, w)
	if is_err {
		return
	}
	result_rows, error := env.DB.QueryContext(req.Context(), `SELECT authors.id, authors.name, authors.cover_filename FROM authors
INNER JOIN authors_views ON authors.id = authors_views.author_id
WHERE authors.name Like '%?%'
GROUP BY authors.id
ORDER  BY COUNT(authors_views.id)
LIMIT ? OFFSET ?`, req.PathValue("name"), page_size, page_index*page_size)
	if check_bool(error, w) {
		return
	}

	result_array, is_err := RowsToSong(result_rows, w)
	if is_err {
		return
	}
	PageSearchReturn(result_array, page_size == len(result_array), w)

}

func (env *Env) AlbumnSearch(w http.ResponseWriter, req *http.Request) {
	ex_a := CreateExtraArgs(env.DB, req.Context(), w)
	page_index, page_size, is_err := PagesStringToNumbers(req, w)
	if is_err {
		return
	}
	results, error := env.DB.QueryContext(req.Context(), `SELECT albumns.id, albumns.name, authors.id, authors.name, 
GROUP_CONCAT(extra_authors.id SEPARATOR ','), GROUP_CONCAT(extra_authors.name SEPARATOR ','), albumns.cover_filename  FROM albumns
INNER JOIN albumns_views ON albumns.id = albumns_views.albumn_id
INNER JOIN songs ON albumns.id=songs.albumn_id 
INNER JOIN authors ON albumns.author_id = authors.id
INNER JOIN authors_to_albumns ON albumns.id = authors_to_albumns.albumn_id
INNER JOIN authors AS extra_authors ON authors.id = extra_authors.id
WHERE albumns.name Like '%?%'
GROUP BY albumns.id
ORDER  BY COUNT(albumns_views.id) 
LIMIT ? OFFSET ?`, req.PathValue("name"), page_size, page_index*page_size)
	if check_bool(error, w) {
		return
	}
	result_return, is_err1 := ex_a.AlbumnReturn(results)
	if is_err1 {
		return
	}

	PageSearchReturn(result_return, page_size, len(result_return), w)
}

func SongAlbumnAdd(ex_a ExtraArgs, albumn_id string) ([]SongReturn, bool) {
	result, error := ex_a.Db.QueryContext(ex_a.Ctx, `SELECT songs.id, songs.name, songs.cover_filename, songs.song_filename,  GROUP_CONCAT(extra_authors.id SEPARATOR ','), GROUP_CONCAT(extra_authors.name SEPARATOR ',') FROM songs
INNER JOIN authors_to_songs ON songs.id = authors_to_songs.songs_id
INNER JOIN authors ON songs.author_id = authors.id
INNER JOIN authors AS extra_authors ON authors_to_songs.author_id = extra_authors.id
WHERE songs.albumn_id = ''
GROUP BY songs.id
ORDER BY songs.albumn_index; `, albumn_id)
	if check_bool(error, ex_a.W) {
		return nil, false

	}

	result_array, err1 := RowsToSong(result, ex_a.W)
	if err1 {
		return nil, true
	}
	return result_array, false
}

func RowsToSong(rows *sql.Rows, w http.ResponseWriter) ([]SongReturn, bool) {
	var for_return []SongReturn
	for rows.Next() {
		var part SongReturn
		var ids string
		var names string
		error := rows.Scan(part.id, part.Name, part.cover_filename, part.song_filename, ids, names)
		if check_bool(error, w) {

		}
		part.Extra_Authors = extra_authors_return(ids, names)
	}
	return for_return, false
}

func (ex_a ExtraArgs) AlbumnReturn(results *sql.Rows) ([]AlbumnReturn, bool) {
	var result_return []AlbumnReturn
	for results.Next() {
		var albumn_result AlbumnReturn
		var ids string
		var names string
		err := results.Scan(albumn_result.id, albumn_result.name, albumn_result.Main_author.id, albumn_result.Main_author.Name, ids, names, albumn_result.cover_filename)
		if check_bool(err, ex_a.W) {
			return nil, true
		}
		albumn_result.Extra_Authors = extra_authors_return(ids, names)
		albumn_songs, is_er := SongAlbumnAdd(ex_a, albumn_result.id)
		if is_er {
			return nil, true
		}
		albumn_result.songs = albumn_songs
		result_return = append(result_return, albumn_result)

	}
	return result_return, false
}
