// internal/postgres/dto/review_dto.go
package dto

type AddReviewRequest struct {
	Score int    `json:"score" binding:"required,min=1,max=5"`
	Text  string `json:"text"`
}

type UpdateReviewRequest struct {
	Score int    `json:"score" binding:"required,min=1,max=5"`
	Text  string `json:"text"`
}

type ReviewResponse struct {
	ID        int    `json:"id"`
	MovieID   int    `json:"movie_id"`
	UserID    int    `json:"user_id"`
	Score     int    `json:"score"`
	Text      string `json:"text,omitempty"`
	CreatedAt string `json:"created_at"`
}
