package valueobject

import (
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Exception"
)

// Year represents a validated release year for a movie.
type Year int

const (
	// MinYear is the year the first motion picture was created.
	MinYear = 1888
	// MaxYear is the current year plus a small buffer.
	MaxYear = 2026
)

// NewYear creates a new Year after validating the input.
// Returns ErrInvalidYear if the year is outside the valid range.
func NewYear(y int) (Year, error) {
	if y < MinYear || y > MaxYear {
		return 0, exception.ErrInvalidYear
	}
	return Year(y), nil
}

// Value returns the year as an int.
func (y Year) Value() int {
	return int(y)
}

// IsZero returns true if the year is the zero value (unset/invalid).
func (y Year) IsZero() bool {
	return y == 0
}
