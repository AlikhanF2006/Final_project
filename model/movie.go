package model

type Movie struct {
	ID          int     `json:"id"`
	TMDBID      int     `json:"tmdb_id"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
}
