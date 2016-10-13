package service

import (
    "github.com/jinzhu/gorm"
)

//Group model used for all groups
type Group struct {
    gorm.Model
    Name        string      `json:"name";gorm:not null`
    Private     bool        `json:"private"`
}

//GroupMember many2many for groups
type GroupMember struct {
    UserID      int     `json:"user_id"`
    GroupID     int     `json:"group_id"`
}

//GroupAdmin denotes who is an admin on a group
type GroupAdmin struct {
    UserID      int     `json:"user_id"`
    GroupID     int     `json:"group_id"`
}

//Post used for group posts
type Post struct {
    gorm.Model
    GroupID     int     `json:"group_id"`
    UserID      int     `json:"user_id"`
    Content     string  `json:"content";gorm:"type:varchar(500)`
    Title       string  `json:"title"`
}

//Comment connects to posts
type Comment struct {
    gorm.Model
    PostID      int     `json:"post_id"`
    Content     string  `json:"content";gorm:"type:varchar(500)`
}

//Token struct handles authentication
type Token struct {
    Key         string   `json:"token"`
    UserID      uint     `json:"user_id"`
    ExpiresAt   int64    `json:"expires_at"`
}
