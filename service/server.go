package service

import (
    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

// NewServer configures and returns a server.
func NewServer() *negroni.Negroni {
    formatter := render.New(render.Options{
        IndentJSON: true,
    })

    n := negroni.Classic()
    mx := mux.NewRouter()

    initRoutes(mx, formatter)
    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
    mx.HandleFunc("/api/groups", getGroupsHandler(formatter)).Methods("GET")
    mx.HandleFunc("/api/groups", postGroupHandler(formatter)).Methods("POST")
    mx.HandleFunc("/api/groups/{id}", getGroupHandler(formatter)).Methods("GET")
}
