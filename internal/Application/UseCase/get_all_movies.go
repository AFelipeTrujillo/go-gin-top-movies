package usecase

import (
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Application/DTO"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Repository"
)

// GetAllMoviesUseCase orchestrates the paginated listing of all movies.
type GetAllMoviesUseCase struct {
	movieRepo repository.IMDBMovieRepository
}

// NewGetAllMoviesUseCase creates a new GetAllMoviesUseCase.
func NewGetAllMoviesUseCase(movieRepo repository.IMDBMovieRepository) *GetAllMoviesUseCase {
	return &GetAllMoviesUseCase{movieRepo: movieRepo}
}

// Execute returns a paginated list of movies.
// Defaults: page=1, pageSize=20 (max 100).
func (uc *GetAllMoviesUseCase) Execute(page, pageSize int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	movies, totalCount, err := uc.movieRepo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	response := dto.NewPaginatedResponse(movies, page, pageSize, totalCount)
	return &response, nil
}
