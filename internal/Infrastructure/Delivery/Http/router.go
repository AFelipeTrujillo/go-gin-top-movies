package http

import (
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Delivery/Http/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all movie-related endpoints.
// It receives the handler struct (which holds the use cases) and registers routes.
func SetupRouter(movieHandler *handlers.IMDBHandler) *gin.Engine {
	router := gin.Default()

	// IMDB movie routes
	imdbGroup := router.Group("/imdb")
	{
		// GET /imdb — List movies (paginated)
		imdbGroup.GET("", movieHandler.GetAll)

		// GET /imdb/top — Top rated movies
		imdbGroup.GET("/top", movieHandler.GetTopRated)

		// GET /imdb/search — Search by keyword (title/overview)
		imdbGroup.GET("/search", movieHandler.Search)

		// GET /imdb/genre — Filter by genre
		imdbGroup.GET("/genre", movieHandler.GetByGenre)

		// GET /imdb/year — Filter by release year
		imdbGroup.GET("/year", movieHandler.GetByYear)

		// GET /imdb/:title — Get movie by exact title
		// Must be last to avoid route conflicts with the above specific routes
		imdbGroup.GET("/:title", movieHandler.GetByTitle)
	}

	return router
}
