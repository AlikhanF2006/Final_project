package main

import (
	"log"

	ginframework "github.com/gin-gonic/gin"

	"github.com/AlikhanF2006/Final_project/configs"
	"github.com/AlikhanF2006/Final_project/pkg/db"

	"github.com/AlikhanF2006/Final_project/ginhandler"
	"github.com/AlikhanF2006/Final_project/postgres"
	"github.com/AlikhanF2006/Final_project/service"
)

func main() {
	configs.LoadConfig()
	db.Connect()
	defer db.Close()

	ginframework.SetMode(ginframework.ReleaseMode)

	movieRepo := postgres.NewMovieRepository()
	reviewRepo := postgres.NewReviewRepository()

	movieSvc := service.NewMovieService(movieRepo)
	reviewSvc := service.NewReviewService(reviewRepo, movieRepo)

	movieH := ginhandler.NewMovieHandler(movieSvc)
	reviewH := ginhandler.NewReviewHandler(reviewSvc)

	r := ginframework.Default()

	r.POST("/movies", movieH.CreateMovie)
	r.GET("/movies", movieH.GetMovies)
	r.GET("/movies/:id", movieH.GetMovieByID)

	r.POST("/movies/:id/reviews", reviewH.AddReview)
	r.GET("/movies/:id/reviews", reviewH.GetReviews)

	log.Println("server running on http://localhost:8080")
	r.Run(":8080")
}
