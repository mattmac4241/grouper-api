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
    api := mux.NewRouter()
    mux := mux.NewRouter()
    repo := &repoHandler{}
    initRoutes(api, formatter, repo)
    mux.Handle("/api", negroni.New(
                NewMiddleware(),
                negroni.Wrap(api),
        ))
    initRoutesWithoutAuth(mux, formatter)
    n.UseHandler(mux)

    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, repo repository) {
    mx.HandleFunc("/groups", getGroupsHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/groups", postGroupHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/groups/{id}", getGroupHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/posts", getPostsHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/posts", postPostHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/posts/{id}", getPostHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/comments", getCommentsHandler(formatter, repo)).Methods("GET")
    mx.HandleFunc("/comments", postCommentHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/comments/{id}", getCommentHandler(formatter, repo)).Methods("GET")
}

func initRoutesWithoutAuth(mx *mux.Router, formatter *render.Render) {
    mx.HandleFunc("/ping", getPingHandler(formatter)).Methods("GET")
}
