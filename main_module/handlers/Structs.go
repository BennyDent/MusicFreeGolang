package handlers

type AlbumnReturn struct {
	id             string
	Main_author    AuthorReturn
	Extra_Authors  []AuthorReturn
	cover_filename string
	songs          []SongReturn
	tags           []string
	genres         []string
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
	tags           []string
	genres         []string
}
