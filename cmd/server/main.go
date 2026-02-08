// main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/AlikhanF2006/Final_project/configs"
	"github.com/AlikhanF2006/Final_project/pkg/db"

	"github.com/AlikhanF2006/Final_project/internal/ginhandler"
	"github.com/AlikhanF2006/Final_project/internal/middleware"
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
	userRepo := postgres.NewUserRepository()

	tmdbClient := tmdb.NewClient(configs.AppConfig.TMDB.ApiKey)

	movieSvc := service.NewMovieService(movieRepo, tmdbClient)
	reviewSvc := service.NewReviewService(reviewRepo, movieRepo)
	userSvc := service.NewUserService(userRepo)

	reviewSvc.StartRatingWorker()

	movieH := ginhandler.NewMovieHandler(movieSvc)
	reviewH := ginhandler.NewReviewHandler(reviewSvc)
	userH := ginhandler.NewUserHandler(userSvc)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./web/static")

	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", userH.Register)
			authGroup.POST("/login", userH.Login)
		}

		public := api.Group("")
		{
			public.GET("/movies", movieH.GetMovies)
			public.GET("/movies/:id", movieH.GetMovieByID)
			public.GET("/movies/tmdb/popular", movieH.GetPopularFromTMDB)
			public.GET("/movies/:id/reviews", reviewH.GetReviews)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(configs.AppConfig.Auth.JWTSecret))
		{
			protected.POST("/movies", movieH.CreateMovie)
			protected.PUT("/movies/:id", movieH.UpdateMovie)
			protected.DELETE("/movies/:id", movieH.DeleteMovie)

			protected.POST("/movies/:id/reviews", reviewH.AddReview)
			protected.PUT("/movies/:id/reviews", reviewH.UpdateReview)
			protected.DELETE("/movies/:id/reviews", reviewH.DeleteReview)
			protected.DELETE("/reviews/:review_id", reviewH.AdminDeleteReview)

			protected.GET("/me", userH.Me)
			protected.PUT("/me", userH.UpdateMe)
			protected.PUT("/me/password", userH.ChangePassword)
			protected.DELETE("/me", userH.DeleteMe)

			protected.GET("/users/:id", userH.GetUserByID)
			protected.DELETE("/users/:id", userH.AdminDeleteUser)
		}
	}

	r.GET("/", func(c *gin.Context) {
		movies := movieSvc.ListMovies()
		c.HTML(200, "index.tmpl", gin.H{"Movies": movies})
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	log.Println("server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
