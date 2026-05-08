package valueobject

import (
	"math"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Exception"
)

// IMDBRating represents a validated IMDB rating (0.0 - 10.0).
type IMDBRating float64

const (
	// MinRating is the minimum possible IMDB rating.
	MinRating = 0.0
	// MaxRating is the maximum possible IMDB rating.
	MaxRating = 10.0
)

// NewIMDBRating creates a new IMDBRating after validating the input.
// The rating is rounded to 1 decimal place.
// Returns ErrInvalidRating if the rating is outside the valid range.
func NewIMDBRating(r float64) (IMDBRating, error) {
	if r < MinRating || r > MaxRating {
		return 0, exception.ErrInvalidRating
	}
	// Round to 1 decimal place
	r = math.Round(r*10) / 10
	return IMDBRating(r), nil
}

// Value returns the rating as a float64.
func (r IMDBRating) Value() float64 {
	return float64(r)
}
