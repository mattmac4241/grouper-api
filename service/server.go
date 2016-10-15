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
    repo := &repoHandler{}
    initRoutes(mx, formatter, repo)
    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, repo repository) {
    mx.HandleFunc("/api/groups", getGroupsHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/api/groups", postGroupHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/api/groups/{id}", getGroupHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/api/posts", getPostsHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/api/posts", postPostHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/api/posts/{id}", getPostHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/api/comments", getCommentsHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/api/comments", postCommentHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/api/comments/{id}", getCommentHandler(formatter, repo)).Methods("GET")
}
