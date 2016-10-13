package service

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

func getGroupsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("CALLED")
        groups := []Group{}
        DB.Find(&groups)
        fmt.Println(req.Body)
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
        var member GroupMember
        var admin GroupAdmin

        key := req.Header.Get("Authorization")
        user, _ := REDIS.Get(key).Result()
        userID, _ := strconv.ParseUint(user, 10, 32)

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

        member.GroupID = group.ID
        admin.GroupID = group.ID
        member.UserID = uint(userID)
        admin.UserID = uint(userID)

        DB.Create(&member)
        DB.Create(&admin)

        formatter.JSON(w, http.StatusCreated, "Group succesfully created.")
    }
}

func postPostHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var post Post
        key := req.Header.Get("Authorization")
        user, _ := REDIS.Get(key).Result()
        userID, _ := strconv.ParseUint(user, 10, 32)

        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &post)
        if err != nil {
            formatter.Text(w, http.StatusBadRequest, "Failed to parse post.")
            return
        }
        post.UserID = uint(userID)
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
        key := req.Header.Get("Authorization")
        user, _ := REDIS.Get(key).Result()
        userID, _ := strconv.ParseUint(user, 10, 32)
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &comment)
        if err != nil {
            formatter.Text(w, http.StatusBadRequest, "Failed to parse comment.")
            return
        }
        comment.UserID = uint(userID)
        err = DB.Create(&comment).Error
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create comment.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "Comment succesfully created.")
    }
}

func checkTokenHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
    serviceClient := authtWebClient{
		rootURL: "http://localhost:3001/auth/token",
	}
    key := req.Header.Get("Authorization")
    w.Header().Set("Content-Type", "application/json")
    if key == "" {
        http.Error(w, "Failed to find token", http.StatusInternalServerError)
        return
    }

    _, err := REDIS.Get(key).Result()
    if err != nil {
        // if the token is not in redis get it and then set it
        token, err := serviceClient.getUserIDFromToken(key)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        now := time.Now().Unix()
        seconds := time.Second * time.Duration(token.ExpiresAt - now)
        REDIS.Set(token.Key, token.UserID, seconds)
    }
    next(w, req)
}
