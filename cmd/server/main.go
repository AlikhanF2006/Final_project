package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/AlikhanF2006/Final_project/configs"
	"github.com/AlikhanF2006/Final_project/pkg/db"

	"github.com/AlikhanF2006/Final_project/internal/ginhandler"
	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/internal/service"
	"github.com/AlikhanF2006/Final_project/internal/tmdb"
)

func main() {
	configs.LoadConfig()

	db.Connect()
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)

	movieRepo := postgres.NewMovieRepository()
	reviewRepo := postgres.NewReviewRepository()

	tmdbClient := tmdb.NewClient(configs.AppConfig.TMDB.ApiKey)

	movieSvc := service.NewMovieService(movieRepo, tmdbClient)
	reviewSvc := service.NewReviewService(reviewRepo, movieRepo)

	movieH := ginhandler.NewMovieHandler(movieSvc)
	reviewH := ginhandler.NewReviewHandler(reviewSvc)

	r := gin.Default()

	r.POST("/movies", movieH.CreateMovie)
	r.GET("/movies", movieH.GetMovies)
	r.GET("/movies/:id", movieH.GetMovieByID)
	r.PUT("/movies/:id", movieH.UpdateMovie)
	r.DELETE("/movies/:id", movieH.DeleteMovie)

	r.POST("/movies/:id/reviews", reviewH.AddReview)
	r.GET("/movies/:id/reviews", reviewH.GetReviews)

	r.GET("/movies/tmdb/popular", movieH.GetPopularFromTMDB)

	log.Println("server running on http://localhost:8080")
	_ = r.Run(":8080")
}
