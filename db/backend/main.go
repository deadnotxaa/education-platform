package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	queryDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "query_duration_seconds",
			Help: "Time taken to execute database queries",
		},
		[]string{"query"},
	)
	queryErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_errors_total",
			Help: "Total number of query errors",
		},
		[]string{"query"},
	)
)

var (
	query1 = `
      SELECT c.name, u.name || ' ' || u.surname as username, e.id as employee_id
      FROM course c
      JOIN purchase p ON c.course_id = p.course_id
      JOIN "user" u ON p.user_id = u.account_id
      JOIN employee e ON u.account_id = e.user_id
      WHERE p.purchase_date > NOW() - INTERVAL '30 days'
      LIMIT 10;`

	query2 = `
      SELECT c.name, u.name || ' ' || u.surname as username, e.id as employee_id, r.name as role_name
      FROM course c
      JOIN purchase p ON c.course_id = p.course_id
      JOIN "user" u ON p.user_id = u.account_id
      JOIN employee e ON u.account_id = e.user_id
      JOIN role r ON e.role_id = r.id
      WHERE r.name LIKE 'HR%'
      ORDER BY p.purchase_date DESC
      LIMIT 5;`

	query3 = `
      SELECT c.name, u.name || ' ' || u.surname as username, e.id as employee_id, r.name as role_name, cs.name as specialization_name
      FROM course c
      JOIN purchase p ON c.course_id = p.course_id
      JOIN "user" u ON p.user_id = u.account_id
      JOIN employee e ON u.account_id = e.user_id
      JOIN role r ON e.role_id = r.id
      JOIN course_specialization cs ON c.specialization_id = cs.id
      WHERE cs.name = 'Programming'
      AND p.purchase_date BETWEEN NOW() - INTERVAL '90 days' AND NOW()
      LIMIT 100;`
)

func init() {
	prometheus.MustRegister(queryDuration)
	prometheus.MustRegister(queryErrors)
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(db)

	http.HandleFunc("/query1", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rows, err := db.Query(query1)
		duration := time.Since(start).Seconds()
		queryDuration.WithLabelValues("query1").Observe(duration)

		if err != nil {
			queryErrors.WithLabelValues("query1").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(rows *sql.Rows) {
			err = rows.Close()
			if err != nil {
				log.Printf("Error closing rows: %v", err)
			}
		}(rows)

		_, err = fmt.Fprintf(w, "Query1 executed in %.2f seconds", duration)
		if err != nil {
			log.Printf("Error executing query1: %v", err)
		}
	})

	http.HandleFunc("/query2", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rows, err := db.Query(query2)
		duration := time.Since(start).Seconds()
		queryDuration.WithLabelValues("query2").Observe(duration)

		if err != nil {
			queryErrors.WithLabelValues("query2").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(rows *sql.Rows) {
			err = rows.Close()
			if err != nil {
				log.Printf("Error closing rows: %v", err)
			}
		}(rows)

		_, err = fmt.Fprintf(w, "Query2 executed in %.2f seconds", duration)
		if err != nil {
			log.Printf("Error executing query2: %v", err)
		}
	})

	http.HandleFunc("/query3", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rows, err := db.Query(query3)
		duration := time.Since(start).Seconds()
		queryDuration.WithLabelValues("query3").Observe(duration)

		if err != nil {
			queryErrors.WithLabelValues("query3").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(rows *sql.Rows) {
			err = rows.Close()
			if err != nil {
				log.Printf("Error closing rows: %v", err)
			}
		}(rows)

		_, err = fmt.Fprintf(w, "Query3 executed in %.2f seconds", duration)
		if err != nil {
			return
		}
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8081", nil))
}
