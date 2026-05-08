package usecase

import (
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Application/DTO"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Repository"
)

// SearchMoviesUseCase orchestrates keyword-based movie search.
type SearchMoviesUseCase struct {
	movieRepo repository.IMDBMovieRepository
}

// NewSearchMoviesUseCase creates a new SearchMoviesUseCase.
func NewSearchMoviesUseCase(movieRepo repository.IMDBMovieRepository) *SearchMoviesUseCase {
	return &SearchMoviesUseCase{movieRepo: movieRepo}
}

// Execute searches movies by keyword in title or overview.
// Returns a paginated response. Defaults: page=1, pageSize=20 (max 100).
func (uc *SearchMoviesUseCase) Execute(query string, page, pageSize int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	movies, totalCount, err := uc.movieRepo.Search(query, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := dto.NewPaginatedResponse(movies, page, pageSize, totalCount)
	return &response, nil
}
