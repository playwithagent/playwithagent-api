package httpserver

import (
	"context"
	"fmt"
	"os"
	"playwithagent-xo/config"
	"playwithagent-xo/internal/api"
	"playwithagent-xo/internal/repository"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config config.Config
	
}

func NewServer(config config.Config) *Server {	
	return &Server{
		config: config,
	}
}

func (s Server) Serve() {
	ctx := context.Background()

	dbPool := connectDB(ctx)
	defer dbPool.Close()

	rdb := connectRedis(ctx)
	defer rdb.Close()

	gameRepo := repository.NewRedisGameRepository(rdb)
	gameHandler := api.NewGameHandler(gameRepo)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	gameGroup := e.Group("/api/v1")
	gameGroup.POST("/games", gameHandler.CreateGame)
	gameGroup.GET("/games/:game_id", gameHandler.GetGame)

	e.GET("/health", s.healthCheck)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
	
}


func connectDB(ctx context.Context) *pgxpool.Pool {
	dbConnString := os.Getenv("DATABASE_URL")
	if dbConnString == ""{
		dbConnString = "postgres://myuser:mypassword@localhost:5432/tictactoe_db?sslmode=disable"
		fmt.Println("WARNING: DATABASE_URL env var not set. Using default.")
	}

	var err error
	dbPool, err := pgxpool.Connect(ctx, dbConnString)
	if err != nil {
	
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := dbPool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to PostgreSQL pool.")
	return dbPool
}

func connectRedis(ctx context.Context) *redis.Client {

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == ""{
		redisURL = "redis://localhost:6379/0"
		fmt.Println("WARNING: REDIS_URL env var not set. Using default.")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse Redis URL: %v\n", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(opt)
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping Redis: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to Redis.")
	return rdb
} 