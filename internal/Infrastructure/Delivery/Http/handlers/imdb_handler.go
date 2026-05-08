package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Application/UseCase"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Exception"
	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
	defaultTopN     = 10
)

// IMDBHandler handles HTTP requests for IMDB movie endpoints.
// It is a thin controller that only parses requests, calls use cases, and returns responses.
type IMDBHandler struct {
	getAllMovies    *usecase.GetAllMoviesUseCase
	getMovieByTitle *usecase.GetMovieByTitleUseCase
	searchMovies    *usecase.SearchMoviesUseCase
}

// NewIMDBHandler creates a new IMDBHandler with the required use cases.
func NewIMDBHandler(
	getAllMovies *usecase.GetAllMoviesUseCase,
	getMovieByTitle *usecase.GetMovieByTitleUseCase,
	searchMovies *usecase.SearchMoviesUseCase,
) *IMDBHandler {
	return &IMDBHandler{
		getAllMovies:    getAllMovies,
		getMovieByTitle: getMovieByTitle,
		searchMovies:    searchMovies,
	}
}

// GetAll handles GET /imdb?page=1&page_size=10
func (h *IMDBHandler) GetAll(c *gin.Context) {
	page, pageSize := parsePaginationParams(c)

	response, err := h.getAllMovies.Execute(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByTitle handles GET /imdb/:title
func (h *IMDBHandler) GetByTitle(c *gin.Context) {
	title := c.Param("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	response, err := h.getMovieByTitle.Execute(title)
	if err != nil {
		if errors.Is(err, exception.ErrMovieNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movie"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Search handles GET /imdb/search?q=keyword&page=1&page_size=10
func (h *IMDBHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	page, pageSize := parsePaginationParams(c)

	response, err := h.searchMovies.Execute(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search movies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTopRated handles GET /imdb/top?n=10
func (h *IMDBHandler) GetTopRated(c *gin.Context) {
	n := defaultTopN
	if nParam := c.Query("n"); nParam != "" {
		if parsed, err := strconv.Atoi(nParam); err == nil && parsed > 0 {
			n = parsed
		}
	}

	// We need a way to get top rated movies.
	// For now, re-use GetAll with a reasonable page size, sorted by rating.
	// The repository's GetTopRated is called via a future use case or directly.
	// TODO: Add GetTopRatedUseCase when needed.
	response, err := h.getAllMovies.Execute(1, n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch top rated movies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByGenre handles GET /imdb/genre?g=Action&page=1&page_size=10
func (h *IMDBHandler) GetByGenre(c *gin.Context) {
	genre := c.Query("g")
	if genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'g' is required"})
		return
	}

	page, pageSize := parsePaginationParams(c)
	_ = page
	_ = pageSize

	// TODO: Add GetByGenreUseCase when needed.
	// For now, fallback to returning a not-implemented response.
	c.JSON(http.StatusNotImplemented, gin.H{"error": "genre filtering not yet implemented"})
}

// GetByYear handles GET /imdb/year?y=2020&page=1&page_size=10
func (h *IMDBHandler) GetByYear(c *gin.Context) {
	yearStr := c.Query("y")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'y' is required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year format"})
		return
	}

	// TODO: Add GetByYearUseCase when needed.
	// For now, fallback to returning a not-implemented response.
	_ = year
	c.JSON(http.StatusNotImplemented, gin.H{"error": "year filtering not yet implemented"})
}

// parsePaginationParams extracts and validates pagination parameters from the request.
func parsePaginationParams(c *gin.Context) (int, int) {
	page := defaultPage
	pageSize := defaultPageSize

	if pageParam := c.Query("page"); pageParam != "" {
		if parsed, err := strconv.Atoi(pageParam); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if parsed, err := strconv.Atoi(pageSizeParam); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	// Enforce max page size
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}
