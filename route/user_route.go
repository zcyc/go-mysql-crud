package route

import (
	"github.com/gorilla/mux"
	"go-example/service"
)

func AddUserHandler(r *mux.Router) {
	r.HandleFunc("/user/list", service.GetUserList).Methods("GET")
	r.HandleFunc("/user/{id}", service.GetUser).Methods("GET")
	r.HandleFunc("/user", service.CreateUser).Methods("POST")
	r.HandleFunc("/user", service.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", service.DeleteUser).Methods("DELETE")
}
