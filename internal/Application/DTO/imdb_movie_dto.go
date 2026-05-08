package dto

import (
	"strconv"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Entity"
)

// IMDBMovieResponse is the DTO returned to the client via HTTP.
// All fields are strings to safely handle TEXT data from SQLite.
type IMDBMovieResponse struct {
	PosterLink   string   `json:"poster_link" example:"https://m.media-amazon.com/images/M/MV5BMDFkYTc0MGEtZmNhMC00ZDIzLWFmNTEtODM1ZmRlYWMwMWFmXkEyXkFqcGdeQXVyMTMxODk2OTU@._V1_UX67_CR0,0,67,98_AL_.jpg"`
	SeriesTitle  string   `json:"series_title" example:"The Shawshank Redemption"`
	ReleasedYear string   `json:"released_year" example:"1994"`
	Certificate  string   `json:"certificate" example:"R"`
	Runtime      string   `json:"runtime" example:"142 min"`
	Genre        string   `json:"genre" example:"Drama"`
	IMDBRating   string   `json:"imdb_rating" example:"9.3"`
	Overview     string   `json:"overview" example:"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency."`
	MetaScore    string   `json:"meta_score" example:"80"`
	Director     string   `json:"director" example:"Frank Darabont"`
	Stars        []string `json:"stars" example:"Tim Robbins,Morgan Freeman,Bob Gunton"`
	NoOfVotes    string   `json:"no_of_votes" example:"2343110"`
	Gross        string   `json:"gross" example:"$28,341,469"`
}

// PaginationRequest represents pagination query parameters.
type PaginationRequest struct {
	Page     int `form:"page" json:"page" example:"1"`
	PageSize int `form:"page_size" json:"page_size" example:"10"`
}

// PaginatedResponse wraps a list of movies with pagination metadata.
type PaginatedResponse struct {
	Data       []IMDBMovieResponse `json:"data"`
	Page       int                 `json:"page" example:"1"`
	PageSize   int                 `json:"page_size" example:"10"`
	TotalCount int64               `json:"total_count" example:"1000"`
	TotalPages int                 `json:"total_pages" example:"100"`
}

// ErrorResponse is the standard error response envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"movie not found"`
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
