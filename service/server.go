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
    n.Use(negroni.HandlerFunc(checkTokenHandler))
    mx := mux.NewRouter()

    initRoutes(mx, formatter)
    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
    mx.HandleFunc("/api/groups", getGroupsHandler(formatter)).Methods("GET")
    mx.HandleFunc("/api/groups", postGroupHandler(formatter)).Methods("POST")
    mx.HandleFunc("/api/groups/{id}", getGroupHandler(formatter)).Methods("GET")
    mx.HandleFunc("/api/posts", getPostsHandler(formatter)).Methods("GET")
    mx.HandleFunc("/api/posts", postPostHandler(formatter)).Methods("POST")
    mx.HandleFunc("/api/posts/{id}", getPostHandler(formatter)).Methods("GET")
    mx.HandleFunc("/api/comments", getCommentsHandler(formatter)).Methods("GET")
    mx.HandleFunc("/api/comments", postCommentHandler(formatter)).Methods("POST")
    mx.HandleFunc("/api/comments/{id}", getCommentHandler(formatter)).Methods("GET")
}
