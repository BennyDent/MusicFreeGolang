package handlers

import (
	"context"
	"database/sql"
	"net/http"
)

type ExtraArgs struct {
	Db  *sql.DB
	Tx  *sql.Tx
	Ctx context.Context
	W   http.ResponseWriter
}

func CreateExtraArgs(db *sql.DB, ctx context.Context, w http.ResponseWriter) ExtraArgs {
	return ExtraArgs{
		Db: db, Ctx: ctx, W: w}
}

type AlbumnReturn struct {
	id             string
	name           string
	Main_author    AuthorReturn
	Extra_Authors  []AuthorReturn
	cover_filename string
	songs          []SongReturn
}

type AuthorReturn struct {
	id   string
	Name string
}

type GenreTagsReturn struct {
	Name string
}

type AuthorReturnWithFilename struct {
	id             string
	Name           string
	cover_filename string
}

type SongReturn struct {
	id             string
	Name           string
	Main_author    AuthorReturn
	Extra_Authors  []AuthorReturn
	song_filename  int64
	albumn_id      string
	file_name      string
	cover_filename string
}
