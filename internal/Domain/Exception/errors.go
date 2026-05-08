package exception

import "errors"

var (
	// ErrInvalidYear is returned when a year is outside the valid range.
	ErrInvalidYear = errors.New("invalid year: must be between 1888 and 2026")
	// ErrInvalidRating is returned when an IMDB rating is outside the valid range.
	ErrInvalidRating = errors.New("invalid rating: must be between 0.0 and 10.0")
	// ErrMovieNotFound is returned when a movie is not found in the repository.
	ErrMovieNotFound = errors.New("movie not found")
)
