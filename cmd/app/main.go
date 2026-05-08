package main

import (
	"fmt"
	"log"

	"github.com/afelipetrujillo/go-gin-top-movies/internal/Application/UseCase"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Config"
	deliveryHttp "github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Delivery/Http"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Delivery/Http/handlers"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Persistence"
	"github.com/afelipetrujillo/go-gin-top-movies/internal/Infrastructure/Persistence/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// --------------------------------------------------
	// 1. Load configuration
	// --------------------------------------------------
	cfg := config.Load()

	// --------------------------------------------------
	// 2. Database connection (GORM)
	// --------------------------------------------------
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Verify the connection and that the "imdb" table exists
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Check that the imdb table exists (we don't auto-migrate, it's pre-existing)
	if !db.Migrator().HasTable(&models.IMDBMovieModel{}) {
		log.Fatal("the 'imdb' table does not exist in the database")
	}

	fmt.Printf("Connected to database: %s\n", cfg.DBPath)

	// --------------------------------------------------
	// 3. Dependency Injection (wiring)
	// --------------------------------------------------
	// Infrastructure
	movieRepo := persistence.NewGORMIMDBRepository(db)

	// Application (UseCases)
	getAllMoviesUC := usecase.NewGetAllMoviesUseCase(movieRepo)
	getMovieByTitleUC := usecase.NewGetMovieByTitleUseCase(movieRepo)
	searchMoviesUC := usecase.NewSearchMoviesUseCase(movieRepo)

	// Delivery (Handlers)
	movieHandler := handlers.NewIMDBHandler(
		getAllMoviesUC,
		getMovieByTitleUC,
		searchMoviesUC,
	)

	// --------------------------------------------------
	// 4. Router setup
	// --------------------------------------------------
	router := deliveryHttp.SetupRouter(movieHandler)

	// --------------------------------------------------
	// 5. Start server
	// --------------------------------------------------
	addr := cfg.Address()
	fmt.Printf("Starting server on %s\n", addr)
	fmt.Printf("API endpoints:\n")
	fmt.Printf("  GET %s/imdb?page=1&page_size=10\n", addr)
	fmt.Printf("  GET %s/imdb/top?n=10\n", addr)
	fmt.Printf("  GET %s/imdb/search?q=keyword\n", addr)
	fmt.Printf("  GET %s/imdb/genre?g=Action\n", addr)
	fmt.Printf("  GET %s/imdb/year?y=2020\n", addr)
	fmt.Printf("  GET %s/imdb/:title\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
