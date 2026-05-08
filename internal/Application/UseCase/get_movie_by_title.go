package usecase

import (
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Application/DTO"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Exception"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Repository"
)

// GetMovieByTitleUseCase orchestrates fetching a single movie by its exact title.
type GetMovieByTitleUseCase struct {
	movieRepo repository.IMDBMovieRepository
}

// NewGetMovieByTitleUseCase creates a new GetMovieByTitleUseCase.
func NewGetMovieByTitleUseCase(movieRepo repository.IMDBMovieRepository) *GetMovieByTitleUseCase {
	return &GetMovieByTitleUseCase{movieRepo: movieRepo}
}

// Execute returns a movie by its exact title.
// Returns exception.ErrMovieNotFound if no movie matches.
func (uc *GetMovieByTitleUseCase) Execute(title string) (*dto.IMDBMovieResponse, error) {
	if title == "" {
		return nil, exception.ErrMovieNotFound
	}

	movie, err := uc.movieRepo.GetByTitle(title)
	if err != nil {
		return nil, err
	}

	response := dto.NewIMDBMovieResponse(*movie)
	return &response, nil
}
