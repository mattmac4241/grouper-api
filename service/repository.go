package service

import (
    "time"
)

type repository interface {
    addGroup(group Group) error
	getGroups() ([]Group, error)
	getGroup(id string) (Group, error)
    addPost(post Post) error
    getPostsByGroup(groupIDs []string) ([]Post, error)
    getPost(id string) (Post, error)
    addComment(comment Comment) error
    getCommentsByPost(postIDs []string) ([]Comment, error)
    getComment(id string) (Comment, error)
    addGroupMember(groupID ,userID uint) error
    addGroupAdmin(groupID, userID uint) error
    redisGetValue(key string) (string, error)
    redisSetValue(key, value string, seconds time.Duration) error
}

type repoHandler struct{}

func (r *repoHandler) addGroup(group Group) error {
    return DB.Create(&group).Error
}

func (r *repoHandler) getGroups() ([]Group, error) {
    var groups []Group
    err := DB.Find(&groups).Error
    return groups, err
}

func (r *repoHandler) getGroup(id string) (Group, error) {
    var group Group
    err := DB.Find(&group, id).Error
    return group, err
}

func (r *repoHandler) addPost(post Post) error {
    return DB.Create(&post).Error
}

func (r *repoHandler) getPostsByGroup(groupIDs []string) ([]Post, error) {
    var posts []Post
    err := DB.Where("group_id in (?)", groupIDs).Find(&posts).Error
    return posts, err
}

func (r *repoHandler) getPost(id string) (Post, error) {
    var post Post
    err := DB.Find(&post).Error
    return post, err
}

func (r *repoHandler) addComment(comment Comment) error {
    return DB.Create(&comment).Error
}

func (r *repoHandler) getCommentsByPost(postIDs []string) ([]Comment, error) {
    var comments []Comment
    err := DB.Where("post_id in (?)", postIDs).Find(&comments).Error
    return comments, err
}

func (r *repoHandler) getComment(id string) (Comment, error) {
    var comment Comment
    err := DB.Find(&comment, id).Error
    return comment, err
}

func (r *repoHandler) addGroupMember(groupID, userID uint) error {
    groupMember := GroupMember{UserID: userID, GroupID: groupID}
    return DB.Create(&groupMember).Error
}

func (r *repoHandler) addGroupAdmin(groupID, userID uint) error {
    adminMember := GroupAdmin{UserID: userID, GroupID: groupID}
    return DB.Create(&adminMember).Error
}

func (r *repoHandler) redisGetValue(key string) (string, error) {
    return REDIS.Get(key).Result()
}

func (r *repoHandler) redisSetValue(key, value string, seconds time.Duration) error {
    return REDIS.Set(key, value, seconds).Err()
}
