package dto

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Year        int    `json:"year" binding:"required"`
	Description string `json:"description"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

type MovieResponse struct {
	ID          int     `json:"id"`
	TMDBID      int     `json:"tmdb_id,omitempty"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
}
