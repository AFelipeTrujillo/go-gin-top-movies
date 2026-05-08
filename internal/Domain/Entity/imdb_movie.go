package entity

// IMDBMovie represents a movie from the IMDB top list.
// Pure domain entity with no framework dependencies.
// All fields use proper Go types; parsing from DB text happens in the repository layer.
type IMDBMovie struct {
	PosterLink  string
	SeriesTitle string
	ReleasedYear int
	Certificate  string
	Runtime      string
	Genre        string
	IMDBRating   float64
	Overview     string
	MetaScore    int
	Director     string
	Stars        []string
	NoOfVotes    int64
	Gross        string
}
