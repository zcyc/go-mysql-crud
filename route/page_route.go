package route

import (
	"go-example/service"

	"github.com/gorilla/mux"
)

func AddPageHandler(r *mux.Router) {
	r.HandleFunc("/", service.Index).Methods("GET", "POST")
	r.HandleFunc("/login", service.Login).Methods("GET", "POST")
	r.HandleFunc("/userinfo", service.Userinfo).Methods("GET", "POST")
	r.HandleFunc("/logout", service.Logout).Methods("GET", "POST")
	r.HandleFunc("/register", service.Register).Methods("GET", "POST")
}
