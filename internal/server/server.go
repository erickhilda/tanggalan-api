package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	db "tanggalan-api/internal/database"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

type Server struct {
	port   int
	dbConn *sql.DB
	q      *db.Queries
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// 1. Setup SQLite connection
	conn, err := sql.Open("sqlite", "./data/tanggalan.db")
	if err != nil {
		panic(err)
	}

	// 2. Setup sqlc queries
	queries := db.New(conn)

	NewServer := &Server{
		port:   port,
		dbConn: conn,
		q:      queries,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
