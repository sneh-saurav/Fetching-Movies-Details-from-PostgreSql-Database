package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Saurav@123"
	dbname   = "postgres"
)

type Movies struct {
	Title         string  `json:"title"`
	Released_year int     `json:"released_year"`
	Rating        float64 `json :"rating"`
	Movie_id      string  `json :"movie_id"`
	Genres        string  `json :"genres"`
}

var db *sql.DB

func main() {

	router := mux.NewRouter()
	// API for getting all movies details
	router.HandleFunc("/AllMovies", getMovies).Methods("GET")

	// API for searching movies
	router.HandleFunc("/MovieByTitle/{title}", getMovieByTitle).Methods("GET")
	router.HandleFunc("/MovieByReleasedYear/{year}", getMovieByReleasedYear).Methods("GET")
	router.HandleFunc("/MovieByRating/{rating}", getMovieByRating).Methods("GET")
	router.HandleFunc("/MovieById/{movie_id}", getMovieById).Methods("GET")

	// API for updating movies
	router.HandleFunc("/UpdateRating/{movie_id}", updateMovieRating).Methods("PUT")
	router.HandleFunc("/UpdateGenres/{movie_id}", updateMovieGenre).Methods("PUT")

	http.ListenAndServe(":8000", router)
}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println(psqlInfo)
	fmt.Println("Successfully connected!")
	return db
}

// Fetching all movies details

func getMovies(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")

	var movies []Movies
	result, err := db.Query("SELECT* from movies")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var movie Movies
		err := result.Scan(&movie.Title, &movie.Released_year, &movie.Rating, &movie.Movie_id, &movie.Genres)
		if err != nil {
			panic(err.Error())
		}
		movies = append(movies, movie)
	}
	json.NewEncoder(w).Encode(movies)
}

// Fetching movie details using title

func getMovieByTitle(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT* FROM movies WHERE title = $1;", params["title"])
	if err != nil {
		panic(err.Error())
	}
	var movie Movies
	for result.Next() {
		err := result.Scan(&movie.Title, &movie.Released_year, &movie.Rating, &movie.Movie_id, &movie.Genres)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(movie)
}

// Fetching movie details using released year

func getMovieByReleasedYear(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(mux.Vars(r)["year"])
	result, err := db.Query("SELECT* FROM movies WHERE released_year = $1;", params["year"])
	if err != nil {
		panic(err.Error())
	}
	var movie Movies
	for result.Next() {
		err := result.Scan(&movie.Title, &movie.Released_year, &movie.Rating, &movie.Movie_id, &movie.Genres)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(movie)
}

// Fetching movie details using rating

func getMovieByRating(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT* FROM movies WHERE rating = $1;", params["rating"])
	if err != nil {
		panic(err.Error())
	}
	var movie Movies
	for result.Next() {
		err := result.Scan(&movie.Title, &movie.Released_year, &movie.Rating, &movie.Movie_id, &movie.Genres)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(movie)
}

// Fetching movie details using movie id

func getMovieById(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT* FROM movies WHERE movie_id = $1;", params["movie_id"])
	if err != nil {
		panic(err.Error())
	}
	var movie Movies
	for result.Next() {
		err := result.Scan(&movie.Title, &movie.Released_year, &movie.Rating, &movie.Movie_id, &movie.Genres)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(movie)
}

// updating movie rating using movie id

func updateMovieRating(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE movies SET rating = $1 WHERE movie_id = $2")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	Rating := keyVal["rating"]

	_, err = stmt.Exec(Rating, params["movie_id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "movie_id = %s : rating is updated", params["movie_id"])
}

// updating movie genres using movie id

func updateMovieGenre(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE movies SET genres = $1 WHERE movie_id = $2")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	Genres := keyVal["genres"]

	_, err = stmt.Exec(Genres, params["movie_id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "movie_id = %s : Genres is updated", params["movie_id"])
}
