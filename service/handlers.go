package service

import (
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

func getGroupsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        groups := []Group{}
        DB.Find(&groups)
        formatter.JSON(w, http.StatusOK, groups)
    }
}

func getGroupHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        id := vars["id"]
        group := Group{}
        err := DB.Find(&group, id).Error
        if err != nil {
            formatter.JSON(w, http.StatusFound, "Group not found")
            return
        }
        formatter.JSON(w, http.StatusOK, group)
    }
}

func postGroupHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var group Group
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &group)
        if err != nil {
            formatter.Text(w, http.StatusBadRequest, "Failed to parse add group command.")
            return
        }

        err = DB.Create(&group).Error
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create group.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "Group succesfully created.")
    }
}

func postPostHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var post Post
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &post)
        if err != nil {
            formatter.Text(w, http.StatusBadRequest, "Failed to parse post.")
            return
        }
        err = DB.Create(&post).Error
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create post.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "Post succesfully created.")
    }
}

func getPostHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        id := vars["id"]
        post := Post{}
        err := DB.Find(&post, id).Error
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Post not found")
            return
        }
        formatter.JSON(w, http.StatusFound, post)
    }
}

func getPostsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        posts := []Post{}
        err := DB.Find(&posts).Error
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find post")
        }
        formatter.JSON(w, http.StatusFound, posts)
    }
}

func getCommentsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        comments := []Comment{}
        DB.Find(&comments)
        formatter.JSON(w, http.StatusFound, comments)
    }
}

func getCommentHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        comment := Comment{}
        vars := mux.Vars(req)
        id := vars["id"]
        err := DB.Find(&comment, id).Error
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find comment")
        }
        formatter.JSON(w, http.StatusFound, comment)
    }
}

func postCommentHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var comment Comment
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &comment)
        if err != nil {
            formatter.Text(w, http.StatusBadRequest, "Failed to parse comment.")
            return
        }
        err = DB.Create(&comment).Error
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create comment.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "Comment succesfully created.")
    }
}
