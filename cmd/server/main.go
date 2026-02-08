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

	api := r.Group("/api")
	{
		api.POST("/movies", movieH.CreateMovie)
		api.GET("/movies", movieH.GetMovies)
		api.GET("/movies/:id", movieH.GetMovieByID)
		api.PUT("/movies/:id", movieH.UpdateMovie)
		api.DELETE("/movies/:id", movieH.DeleteMovie)

		api.POST("/movies/:id/reviews", reviewH.AddReview)
		api.GET("/movies/:id/reviews", reviewH.GetReviews)

		api.GET("/movies/tmdb/popular", movieH.GetPopularFromTMDB)
	}

	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	log.Println("server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
