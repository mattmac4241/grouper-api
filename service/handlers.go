package service

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

func getGroupsHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        groups, err := repo.getGroups()
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to get groups")
            return
        }
        formatter.JSON(w, http.StatusOK, groups)
    }
}

func getGroupHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        id := vars["id"]
        group, err := repo.getGroup(id)
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Group not found")
            return
        }
        formatter.JSON(w, http.StatusOK, group)
    }
}

func postGroupHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var group Group

        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &group)
        if err != nil || (group == Group{}) {
            formatter.JSON(w, http.StatusBadRequest, "Failed to parse add group command.")
            return
        }

        err = repo.addGroup(group)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create group.")
            return
        }

        key := req.Header.Get("Authorization")
        user, err := repo.redisGetValue(key)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to get user from token.")
            return
        }

        userID, _ := strconv.ParseUint(user, 10, 32)

        err = repo.addGroupMember(group.ID, uint(userID))
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to join group.")
            return
        }

        err = repo.addGroupAdmin(group.ID, uint(userID))
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to join admin of group.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "Group succesfully created.")
    }
}

func postPostHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var post Post

        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &post)
        if err != nil || (post == Post{}) {
            formatter.JSON(w, http.StatusBadRequest, "Failed to parse post.")
            return
        }

        key := req.Header.Get("Authorization")
        user, err := repo.redisGetValue(key)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to get user from token.")
            return
        }
        userID, _ := strconv.ParseUint(user, 10, 32)

        post.UserID = uint(userID)
        err = repo.addPost(post)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create post.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "Post succesfully created.")
    }
}

func getPostHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        id := vars["id"]
        post, err := repo.getPost(id)
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Post not found")
            return
        }
        formatter.JSON(w, http.StatusOK, post)
    }
}

func getPostsHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        groups := req.URL.Query()["group"]
        posts, err := repo.getPostsByGroup(groups)
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find posts")
        }
        formatter.JSON(w, http.StatusOK, posts)
    }
}

func getCommentsHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        posts := req.URL.Query()["post"]
        comments, err := repo.getCommentsByPost(posts)
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find comments")
        }
        formatter.JSON(w, http.StatusOK, comments)
    }
}

func getCommentHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        id := vars["id"]
        comment, err := repo.getComment(id)
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find comment")
        }
        formatter.JSON(w, http.StatusOK, comment)
    }
}

func postCommentHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var comment Comment

        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &comment)
        if err != nil ||( comment == Comment{}) {
            formatter.Text(w, http.StatusBadRequest, "Failed to parse comment.")
            return
        }

        key := req.Header.Get("Authorization")
        user, err := repo.redisGetValue(key)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to get user from token.")
            return
        }
        userID, _ := strconv.ParseUint(user, 10, 32)
        comment.UserID = uint(userID)
        err = repo.addComment(comment)
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
        REDIS.Set(token.Key, string(token.UserID), seconds)
    }
    next(w, req)
}
