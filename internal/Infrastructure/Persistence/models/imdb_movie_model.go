package models

import (
	"strconv"
	"strings"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Domain/Entity"
)

// IMDBMovieModel is the GORM model mapping to the existing "imdb" table.
// This struct lives in Infrastructure layer and is NOT used in Domain.
// It handles all the TEXT-to-Go-type parsing from SQLite.
type IMDBMovieModel struct {
	PosterLink   string `gorm:"column:Poster_Link;primaryKey"`
	SeriesTitle  string `gorm:"column:Series_Title"`
	ReleasedYear string `gorm:"column:Released_Year"`
	Certificate  string `gorm:"column:Certificate"`
	Runtime      string `gorm:"column:Runtime"`
	Genre        string `gorm:"column:Genre"`
	IMDBRating   string `gorm:"column:IMDB_Rating"`
	Overview     string `gorm:"column:Overview"`
	MetaScore    string `gorm:"column:Meta_score"`
	Director     string `gorm:"column:Director"`
	Star1        string `gorm:"column:Star1"`
	Star2        string `gorm:"column:Star2"`
	Star3        string `gorm:"column:Star3"`
	Star4        string `gorm:"column:Star4"`
	NoOfVotes    string `gorm:"column:No_of_Votes"`
	Gross        string `gorm:"column:Gross"`
}

// TableName overrides the default table name for GORM.
// The existing table is named "imdb".
func (IMDBMovieModel) TableName() string {
	return "imdb"
}

// ToEntity converts the GORM model into a pure Domain entity.
// Parsing errors for numeric fields are silently ignored, falling back to zero values.
func (m *IMDBMovieModel) ToEntity() entity.IMDBMovie {
	// Parse Released_Year (TEXT) -> int
	releasedYear := 0
	if m.ReleasedYear != "" {
		if year, err := strconv.Atoi(strings.TrimSpace(m.ReleasedYear)); err == nil {
			releasedYear = year
		}
	}

	// Parse IMDB_Rating (TEXT) -> float64
	imdbRating := 0.0
	if m.IMDBRating != "" {
		if rating, err := strconv.ParseFloat(strings.TrimSpace(m.IMDBRating), 64); err == nil {
			imdbRating = rating
		}
	}

	// Parse Meta_score (TEXT) -> int
	metaScore := 0
	if m.MetaScore != "" {
		if score, err := strconv.Atoi(strings.TrimSpace(m.MetaScore)); err == nil {
			metaScore = score
		}
	}

	// Parse No_of_Votes (TEXT) -> int64
	noOfVotes := int64(0)
	if m.NoOfVotes != "" {
		if votes, err := strconv.ParseInt(strings.TrimSpace(m.NoOfVotes), 10, 64); err == nil {
			noOfVotes = votes
		}
	}

	// Build stars slice from the 4 star columns, filtering out empty strings
	stars := make([]string, 0, 4)
	if m.Star1 != "" {
		stars = append(stars, m.Star1)
	}
	if m.Star2 != "" {
		stars = append(stars, m.Star2)
	}
	if m.Star3 != "" {
		stars = append(stars, m.Star3)
	}
	if m.Star4 != "" {
		stars = append(stars, m.Star4)
	}

	return entity.IMDBMovie{
		PosterLink:   m.PosterLink,
		SeriesTitle:  m.SeriesTitle,
		ReleasedYear: releasedYear,
		Certificate:  m.Certificate,
		Runtime:      m.Runtime,
		Genre:        m.Genre,
		IMDBRating:   imdbRating,
		Overview:     m.Overview,
		MetaScore:    metaScore,
		Director:     m.Director,
		Stars:        stars,
		NoOfVotes:    noOfVotes,
		Gross:        m.Gross,
	}
}

// ToEntitySlice converts a slice of GORM models into a slice of Domain entities.
func ToEntitySlice(models []IMDBMovieModel) []entity.IMDBMovie {
	entities := make([]entity.IMDBMovie, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities
}
