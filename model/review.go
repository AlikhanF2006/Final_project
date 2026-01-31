package model

import "time"

type Review struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movieId"`
	UserID    int       `json:"userId"`
	Score     int       `json:"score"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
