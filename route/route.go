package route

import (
	"github.com/gorilla/mux"
)

func AddHandler(r *mux.Router) {
	AddUserHandler(r)
	AddPageHandler(r)
}
