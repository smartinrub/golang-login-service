package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

// User store email and BCrypt password
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="email">Email</label>
    <input type="text" id="email" name="email">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

const successPage = `
<h1>Welcome</h1>
<a href="/">Go to Home</a>
`

const failPage = `
<h1>Denied</h1>
<a href="/">Go to Home</a>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	pass := request.FormValue("password")
	redirectTarget := "/fail"
	if checkPassword(email, pass) {
		redirectTarget = "/success"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func successPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, successPage)
}

func failPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, failPage)
}

func main() {
	port := ":8090"
	router := mux.NewRouter()
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/success", successPageHandler)
	router.HandleFunc("/fail", failPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	http.Handle("/", router)
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(port, nil)
}

func checkPassword(email string, password string) bool {
	mysqlURL := os.Getenv("MYSQL_URL")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDb := os.Getenv("MYSQL_DB")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlConnection := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlURL + ":" + mysqlPort + ")/" + mysqlDb
	db, err := sql.Open("mysql", mysqlConnection)
	if err != nil {
		panic(err.Error())
	}

	// defer the close until the main function has finished
	// executing
	defer db.Close()

	// Execute the query
	result, err := db.Query("SELECT email, password FROM user WHERE email=?", email)
	if err != nil {
		panic(err.Error())
	}

	if result.Next() == false {
		log.Printf("User not found")
		return false
	}

	var user User
	err = result.Scan(&user.Email, &user.Password)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("\nEmail: %s \nHashed Password: %s", user.Email, user.Password)
	if checkHash(user.Password, password) {
		log.Printf("Valid Password")
		return true
	}

	log.Printf("Wrong Password")
	return false
}

func checkHash(hashed string, supplied string) bool {
	error := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(supplied))

	return error == nil
}
