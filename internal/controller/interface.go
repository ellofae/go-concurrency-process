package controller

import (
	"github.com/gorilla/mux"
)

type IHandler interface {
	Register(*mux.Router)
}
