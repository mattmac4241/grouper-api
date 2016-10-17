package service

import (
    "net/http"
    "os"
    "time"
)

type Middleware struct {
    auth bool
}

// New`Middleware is a struct that has a ServeHTTP method
func NewMiddleware() *Middleware {
    return &Middleware{true}
}

// The middleware handler
func (l *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
    serviceClient := authtWebClient{
        rootURL: os.Getenv("AUTH_URL"),
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
