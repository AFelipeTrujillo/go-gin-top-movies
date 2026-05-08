package persistence

import (
	"fmt"
	"strings"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Entity"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Exception"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Persistence/models"
	"gorm.io/gorm"
)

// GORMIMDBRepository implements repository.IMDBMovieRepository using GORM.
// It is the Infrastructure layer's concrete implementation of the Domain repository interface.
type GORMIMDBRepository struct {
	db *gorm.DB
}

// NewGORMIMDBRepository creates a new GORM-based IMDB repository.
func NewGORMIMDBRepository(db *gorm.DB) *GORMIMDBRepository {
	return &GORMIMDBRepository{db: db}
}

// GetAll returns a paginated list of movies and the total count.
func (r *GORMIMDBRepository) GetAll(page, pageSize int) ([]entity.IMDBMovie, int64, error) {
	var modelList []models.IMDBMovieModel
	var totalCount int64

	// Get total count
	if err := r.db.Model(&models.IMDBMovieModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movies: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.Order("Series_Title ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&modelList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list movies: %w", err)
	}

	return models.ToEntitySlice(modelList), totalCount, nil
}

// GetByTitle returns a movie by its exact title.
// Returns exception.ErrMovieNotFound if not found.
func (r *GORMIMDBRepository) GetByTitle(title string) (*entity.IMDBMovie, error) {
	var model models.IMDBMovieModel

	if err := r.db.Where("Series_Title = ?", title).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.ErrMovieNotFound
		}
		return nil, fmt.Errorf("failed to get movie by title: %w", err)
	}

	movie := model.ToEntity()
	return &movie, nil
}

// Search returns movies matching a keyword in title or overview.
func (r *GORMIMDBRepository) Search(query string, page, pageSize int) ([]entity.IMDBMovie, int64, error) {
	var modelList []models.IMDBMovieModel
	var totalCount int64

	likePattern := "%" + strings.TrimSpace(query) + "%"

	// Get total count of matching records
	baseQuery := r.db.Model(&models.IMDBMovieModel{}).
		Where("Series_Title LIKE ? OR Overview LIKE ?", likePattern, likePattern)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := baseQuery.
		Order("Series_Title ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&modelList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search movies: %w", err)
	}

	return models.ToEntitySlice(modelList), totalCount, nil
}

// GetTopRated returns the top n movies by IMDB rating.
func (r *GORMIMDBRepository) GetTopRated(n int) ([]entity.IMDBMovie, error) {
	var modelList []models.IMDBMovieModel

	// We need to cast the TEXT IMDB_Rating to REAL for proper numeric ordering.
	// SQLite supports CAST() for this purpose.
	if err := r.db.
		Order("CAST(IMDB_Rating AS REAL) DESC, No_of_Votes DESC").
		Limit(n).
		Find(&modelList).Error; err != nil {
		return nil, fmt.Errorf("failed to get top rated movies: %w", err)
	}

	return models.ToEntitySlice(modelList), nil
}

// GetByGenre returns movies filtered by genre.
func (r *GORMIMDBRepository) GetByGenre(genre string, page, pageSize int) ([]entity.IMDBMovie, int64, error) {
	var modelList []models.IMDBMovieModel
	var totalCount int64

	likePattern := "%" + strings.TrimSpace(genre) + "%"

	// Get total count of matching records
	baseQuery := r.db.Model(&models.IMDBMovieModel{}).
		Where("Genre LIKE ?", likePattern)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count genre results: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := baseQuery.
		Order("Series_Title ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&modelList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get movies by genre: %w", err)
	}

	return models.ToEntitySlice(modelList), totalCount, nil
}

// GetByYear returns movies filtered by release year.
func (r *GORMIMDBRepository) GetByYear(year int, page, pageSize int) ([]entity.IMDBMovie, int64, error) {
	var modelList []models.IMDBMovieModel
	var totalCount int64

	yearStr := fmt.Sprintf("%d", year)

	// Get total count of matching records
	baseQuery := r.db.Model(&models.IMDBMovieModel{}).
		Where("Released_Year = ?", yearStr)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count year results: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := baseQuery.
		Order("Series_Title ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&modelList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get movies by year: %w", err)
	}

	return models.ToEntitySlice(modelList), totalCount, nil
}
