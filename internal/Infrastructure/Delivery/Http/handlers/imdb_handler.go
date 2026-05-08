package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Application/DTO"
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
// @Summary      List all movies (paginated)
// @Description  Returns a paginated list of all movies from the IMDB database.
// @Description  Use the "page" and "page_size" query parameters to control pagination.
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number (default: 1)"  default(1)
// @Param        page_size query     int  false  "Items per page (default: 10, max: 100)"  default(10)
// @Success      200  {object}  dto.PaginatedResponse  "Paginated list of movies"
// @Failure      500  {object}  dto.ErrorResponse      "Internal server error"
// @Router       /imdb [get]
func (h *IMDBHandler) GetAll(c *gin.Context) {
	page, pageSize := parsePaginationParams(c)

	response, err := h.getAllMovies.Execute(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByTitle handles GET /imdb/:title
// @Summary      Get a movie by its exact title
// @Description  Returns a single movie matching the exact title specified in the URL path.
// @Description  The title is URL-decoded automatically by Gin.
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        title  path      string  true  "Movie title (URL encoded)"  example("The Shawshank Redemption")
// @Success      200  {object}  dto.IMDBMovieResponse  "Movie found"
// @Failure      400  {object}  dto.ErrorResponse       "Bad request (missing title)"
// @Failure      404  {object}  dto.ErrorResponse       "Movie not found"
// @Failure      500  {object}  dto.ErrorResponse       "Internal server error"
// @Router       /imdb/{title} [get]
func (h *IMDBHandler) GetByTitle(c *gin.Context) {
	title := c.Param("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "title is required"})
		return
	}

	response, err := h.getMovieByTitle.Execute(title)
	if err != nil {
		if errors.Is(err, exception.ErrMovieNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch movie"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Search handles GET /imdb/search?q=keyword&page=1&page_size=10
// @Summary      Search movies by keyword
// @Description  Searches movies where the title or overview contains the given keyword (case-insensitive).
// @Description  Supports pagination with "page" and "page_size" parameters.
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        q         query     string  true   "Search keyword"  example("Shawshank")
// @Param        page      query     int     false  "Page number (default: 1)"  default(1)
// @Param        page_size query     int     false  "Items per page (default: 10, max: 100)"  default(10)
// @Success      200  {object}  dto.PaginatedResponse  "Paginated search results"
// @Failure      400  {object}  dto.ErrorResponse       "Bad request (missing query)"
// @Failure      500  {object}  dto.ErrorResponse       "Internal server error"
// @Router       /imdb/search [get]
func (h *IMDBHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "query parameter 'q' is required"})
		return
	}

	page, pageSize := parsePaginationParams(c)

	response, err := h.searchMovies.Execute(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to search movies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTopRated handles GET /imdb/top?n=10
// @Summary      Get top rated movies
// @Description  Returns the top N rated movies from the IMDB database.
// @Description  Use the "n" query parameter to control how many movies to return (default: 10).
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        n    query     int  false  "Number of top movies to return (default: 10)"  default(10)
// @Success      200  {object}  dto.PaginatedResponse  "Top rated movies"
// @Failure      500  {object}  dto.ErrorResponse       "Internal server error"
// @Router       /imdb/top [get]
func (h *IMDBHandler) GetTopRated(c *gin.Context) {
	n := defaultTopN
	if nParam := c.Query("n"); nParam != "" {
		if parsed, err := strconv.Atoi(nParam); err == nil && parsed > 0 {
			n = parsed
		}
	}

	response, err := h.getAllMovies.Execute(1, n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch top rated movies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByGenre handles GET /imdb/genre?g=Action&page=1&page_size=10
// @Summary      Filter movies by genre
// @Description  Returns a paginated list of movies filtered by genre.
// @Description  Use the "g" query parameter to specify the genre (e.g., "Action", "Drama", "Comedy").
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        g         query     string  true   "Genre to filter by"  example("Action")
// @Param        page      query     int     false  "Page number (default: 1)"  default(1)
// @Param        page_size query     int     false  "Items per page (default: 10, max: 100)"  default(10)
// @Success      200  {object}  dto.PaginatedResponse  "Filtered movies by genre"
// @Failure      400  {object}  dto.ErrorResponse       "Bad request (missing genre)"
// @Failure      501  {object}  dto.ErrorResponse       "Not implemented"
// @Router       /imdb/genre [get]
func (h *IMDBHandler) GetByGenre(c *gin.Context) {
	genre := c.Query("g")
	if genre == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "query parameter 'g' is required"})
		return
	}

	page, pageSize := parsePaginationParams(c)
	_ = page
	_ = pageSize

	c.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "genre filtering not yet implemented"})
}

// GetByYear handles GET /imdb/year?y=2020&page=1&page_size=10
// @Summary      Filter movies by release year
// @Description  Returns a paginated list of movies filtered by release year.
// @Description  Use the "y" query parameter to specify the year (e.g., "1994", "2020").
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        y         query     int     true   "Release year to filter by"  example(1994)
// @Param        page      query     int     false  "Page number (default: 1)"  default(1)
// @Param        page_size query     int     false  "Items per page (default: 10, max: 100)"  default(10)
// @Success      200  {object}  dto.PaginatedResponse  "Filtered movies by year"
// @Failure      400  {object}  dto.ErrorResponse       "Bad request (missing or invalid year)"
// @Failure      501  {object}  dto.ErrorResponse       "Not implemented"
// @Router       /imdb/year [get]
func (h *IMDBHandler) GetByYear(c *gin.Context) {
	yearStr := c.Query("y")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "query parameter 'y' is required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid year format"})
		return
	}

	_ = year
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "year filtering not yet implemented"})
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
