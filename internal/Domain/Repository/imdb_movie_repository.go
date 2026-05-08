package repository

import (
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Entity"
)

// IMDBMovieRepository defines the contract for IMDB movie data access.
// Implementations belong to the Infrastructure layer.
type IMDBMovieRepository interface {
	// GetAll returns a paginated list of movies and the total count.
	GetAll(page, pageSize int) ([]entity.IMDBMovie, int64, error)

	// GetByTitle returns a movie by its exact title.
	// Returns exception.ErrMovieNotFound if not found.
	GetByTitle(title string) (*entity.IMDBMovie, error)

	// Search returns movies matching a keyword in title or overview.
	Search(query string, page, pageSize int) ([]entity.IMDBMovie, int64, error)

	// GetTopRated returns the top n movies by IMDB rating.
	GetTopRated(n int) ([]entity.IMDBMovie, error)

	// GetByGenre returns movies filtered by genre.
	GetByGenre(genre string, page, pageSize int) ([]entity.IMDBMovie, int64, error)

	// GetByYear returns movies filtered by release year.
	GetByYear(year int, page, pageSize int) ([]entity.IMDBMovie, int64, error)
}
