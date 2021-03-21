package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"lenslocked/controllers"
	"net/http"
	"os"
)

var notTemplate *template.Template


func notf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := notTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()
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
