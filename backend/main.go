package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Job struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

var db *sql.DB

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB Ping failed:", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	createTables()

	http.HandleFunc("/api/jobs", corsMiddleware(getJobsHandler))
	http.HandleFunc("/api/health", corsMiddleware(healthHandler))

	// Render yêu cầu health check ở root
	http.HandleFunc("/", corsMiddleware(healthHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		location VARCHAR(255),
		created_at TIMESTAMP DEFAULT NOW()
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Cannot create table:", err)
	}
	log.Println("✅ Bảng 'jobs' đã sẵn sàng")
}

func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, title, description, location FROM jobs ORDER BY id DESC")
	if err != nil {
		log.Println("Query error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	jobs := []Job{}
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Title, &j.Description, &j.Location); err != nil {
			continue
		}
		jobs = append(jobs, j)
	}

	if jobs == nil {
		jobs = []Job{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
