package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"lenslocked/controllers"
	"lenslocked/models"
	"net/http"
	"os"
)

var notTemplate *template.Template
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "lenslocked_dev"
)

func notf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := notTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func main() {
	// Create a DB connection string and then use it to
	// create our model services.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)
	notTemplate, _ = template.ParseFiles("views/404.gohtml")

	var nf http.Handler = http.HandlerFunc(notf)
	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.NotFoundHandler = nf
	port := getPort()
	http.ListenAndServe(port, r)
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":3000"
}
