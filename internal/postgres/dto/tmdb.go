package dto

type TMDBMovieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
}

type TMDBVideosResponse struct {
	Results []TMDBVideo `json:"results"`
}

type TMDBVideo struct {
	Key  string `json:"key"`
	Site string `json:"site"`
	Type string `json:"type"`
}
