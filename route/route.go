package route

import (
	"github.com/gorilla/mux"
	"go-example/dao"
)

func AddUserHandler(r *mux.Router) {
	r.HandleFunc("/user/list", dao.GetUserList).Methods("GET")
	r.HandleFunc("/user/{id}", dao.GetUser).Methods("GET")
	r.HandleFunc("/user", dao.CreateUser).Methods("POST")
	r.HandleFunc("/user", dao.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", dao.DeleteUser).Methods("DELETE")
}
