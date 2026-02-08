package tmdb

type MovieDTO struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	Rating      float64 `json:"vote_average"`
}
