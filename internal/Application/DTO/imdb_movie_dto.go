package dto

import (
	"strconv"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Entity"
)

// IMDBMovieResponse is the DTO returned to the client via HTTP.
// All fields are strings to safely handle TEXT data from SQLite.
type IMDBMovieResponse struct {
	PosterLink   string   `json:"poster_link"`
	SeriesTitle  string   `json:"series_title"`
	ReleasedYear string   `json:"released_year"`
	Certificate  string   `json:"certificate"`
	Runtime      string   `json:"runtime"`
	Genre        string   `json:"genre"`
	IMDBRating   string   `json:"imdb_rating"`
	Overview     string   `json:"overview"`
	MetaScore    string   `json:"meta_score"`
	Director     string   `json:"director"`
	Stars        []string `json:"stars"`
	NoOfVotes    string   `json:"no_of_votes"`
	Gross        string   `json:"gross"`
}

// PaginationRequest represents pagination query parameters.
type PaginationRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// PaginatedResponse wraps a list of movies with pagination metadata.
type PaginatedResponse struct {
	Data       []IMDBMovieResponse `json:"data"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalCount int64               `json:"total_count"`
	TotalPages int                 `json:"total_pages"`
}

// NewIMDBMovieResponse converts a domain entity into a response DTO.
func NewIMDBMovieResponse(movie entity.IMDBMovie) IMDBMovieResponse {
	return IMDBMovieResponse{
		PosterLink:   movie.PosterLink,
		SeriesTitle:  movie.SeriesTitle,
		ReleasedYear: strconv.Itoa(movie.ReleasedYear),
		Certificate:  movie.Certificate,
		Runtime:      movie.Runtime,
		Genre:        movie.Genre,
		IMDBRating:   strconv.FormatFloat(movie.IMDBRating, 'f', 1, 64),
		Overview:     movie.Overview,
		MetaScore:    strconv.Itoa(movie.MetaScore),
		Director:     movie.Director,
		Stars:        movie.Stars,
		NoOfVotes:    strconv.FormatInt(movie.NoOfVotes, 10),
		Gross:        movie.Gross,
	}
}

// NewPaginatedResponse builds a paginated response from domain entities.
func NewPaginatedResponse(movies []entity.IMDBMovie, page, pageSize int, totalCount int64) PaginatedResponse {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	dtos := make([]IMDBMovieResponse, len(movies))
	for i, m := range movies {
		dtos[i] = NewIMDBMovieResponse(m)
	}

	return PaginatedResponse{
		Data:       dtos,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
