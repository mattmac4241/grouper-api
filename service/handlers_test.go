package service

import (
    "bytes"
    "errors"
    "testing"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "strconv"
    "time"


    "github.com/unrolled/render"
    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
)

var (
    formatter = render.New(render.Options{
        IndentJSON: true,
    })
)

type repoTest struct {
    groups          []Group
    posts           []Post
    comments        []Comment
    groupMembers    []GroupMember
    groupAdmins     []GroupAdmin
    redis           map[string]string
}

func (r *repoTest) addGroup(group Group) error {
    r.groups = append(r.groups, group)
    return nil
}

func (r *repoTest) getGroups() ([]Group, error) {
    return r.groups, nil
}

func (r *repoTest) getGroup(id string) (Group, error) {
    userID, _ := strconv.ParseUint(id, 10, 32)

    for _, group := range r.groups {
        if group.ID == uint(userID) {
            return group, nil
        }
    }

    return Group{}, errors.New("Group not found")
}

func (r *repoTest) addPost(post Post) error {
    r.posts = append(r.posts, post)
    return nil
}

func (r *repoTest) getPostsByGroup(groupsIDs []string) ([]Post, error) {
    var posts []Post
    for _, post := range r.posts {
        for _, group := range groupsIDs {
            groupID,_ := strconv.ParseUint(group, 10, 32)
            if uint(groupID) == post.GroupID {
                posts = append(posts, post)
            }
        }
    }
    return posts, nil
}

func (r *repoTest) getPost(id string) (Post, error) {
    postID,_ := strconv.ParseUint(id, 10, 32)

    for _, post := range r.posts {
        if post.ID == uint(postID) {
            return post, nil
        }
    }
    return Post{}, errors.New("Post not found")
}

func (r *repoTest) addComment(comment Comment) error {
    r.comments = append(r.comments, comment)
    return nil
}

func (r *repoTest) getCommentsByPost(postIDs []string) ([]Comment, error) {
    var comments []Comment
    for _, comment := range r.comments {
        for _, post := range postIDs {
            postID,_ := strconv.ParseUint(post, 10, 32)
            if uint(postID) == comment.PostID {
                comments = append(comments, comment)
            }
        }
    }
    return comments, nil
}

func (r *repoTest) getComment(id string) (Comment, error) {
    commentID,_ := strconv.ParseUint(id, 10, 32)

    for _, comment := range r.comments {
        if comment.ID == uint(commentID) {
            return comment, nil
        }
    }
    return Comment{}, errors.New("Comment not found")
}

func (r *repoTest) addGroupMember(groupID, userID uint) error {
    groupMember := GroupMember{UserID: userID, GroupID: groupID}
    r.groupMembers = append(r.groupMembers, groupMember)
    return nil
}

func (r *repoTest) addGroupAdmin(groupID, userID uint) error {
    groupAdmin := GroupAdmin{UserID: userID, GroupID: groupID}
    r.groupAdmins = append(r.groupAdmins, groupAdmin)
    return nil
}

func (r *repoTest) redisGetValue(key string) (string, error) {
    value, prs := r.redis[key]
    if prs == false {
        return "", errors.New("Key not found:")
    }
    return value, nil
}

func (r *repoTest) redisSetValue(key, value string, seconds time.Duration) error {
    r.redis[key] = value
    return nil
}

func TestGetGroupsHandler(t *testing.T) {
    group1 := Group{Name:"Group1", Private:false}
    group2 := Group{Name:"Group2", Private:false}

    repo := &repoTest{}
    repo.addGroup(group1)
    repo.addGroup(group2)

    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(getGroupsHandler(formatter, repo)))
    defer server.Close()
    req, _ := http.NewRequest("GET", server.URL, nil)

    resp, err := client.Do(req)

    if err != nil {
        t.Error("Errored when sending request to the server", err)
        return
    }

    defer resp.Body.Close()
    payload, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Error("Failed to read response from server", err)
    }

    var groups []Group
    err = json.Unmarshal(payload, &groups)
    if err != nil {
        t.Errorf("Could not unmarshal payload into []groups slice")
    }

    if len(groups) != 2 {
        t.Errorf("Expected an group len 2 , got %d", len(groups))
    }
}

func TestGetGroupHandlerNotValidGroup(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    repo  := &repoTest{}

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/groups/1", nil)
    server.ServeHTTP(recorder, request)
    if recorder.Code != http.StatusNotFound {
        t.Errorf("Expected %v; received %v", http.StatusNotFound, recorder.Code)
    }
}

func TestGetGroupHandlerValidGroup(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )

    group := Group{Name: "test", Private: false}
    group.ID = 1
    repo  := &repoTest{}
    repo.addGroup(group)

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/groups/1", nil)
    server.ServeHTTP(recorder, request)
    if recorder.Code != http.StatusOK {
        t.Errorf("Expected %v; received %v", http.StatusOK, recorder.Code)
    }

    var groupResponse Group
    err := json.Unmarshal(recorder.Body.Bytes(), &groupResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  groupResponse.Name != group.Name && groupResponse.Private != group.Private {
        t.Errorf("Expected group recieved to equal; received %d", groupResponse)
    }
}

func TestPostGroupHandlerInvalidJSON(t *testing.T) {
    repo := &repoTest{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postGroupHandler(formatter, repo)))
    defer server.Close()

    body := []byte("this is not valid json")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
    }
    res, err := client.Do(req)
    if err != nil {
        t.Errorf("Error in POST to createMatchHandler: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending invalid JSON should result in a bad request from server.")
    }
}

func TestPostGroupHandlerNotGroup(t *testing.T) {
    repo := &repoTest{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postGroupHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"test\":\"Not user.\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }
}

func TestPostGroupHandlerValidGroup(t *testing.T) {
    repo := &repoTest{}
    repo.redis =  make(map[string]string)
    client := &http.Client{}
    seconds := time.Second * time.Duration(time.Now().Unix() - time.Now().Unix())
    repo.redisSetValue("token", "1", seconds)

    server := httptest.NewServer(http.HandlerFunc(postGroupHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"name\":\"testname\",\n\"private\":false}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", "token")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode == http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }

    if len(repo.groups) != 1 {
        t.Error("Expected one group")
    }

    group := repo.groups[0]
    if group.Name != "name" && group.Private != false {
        t.Error("Group was not set correctly")
    }
    if len(repo.groupMembers) != 1 {
        t.Error("Expected one group member")
    }

    if len(repo.groupAdmins) != 1 {
        t.Error("Expected one group member")
    }
}

func TestPostPostHandlerInvalidJSON(t *testing.T) {
    repo := &repoTest{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postPostHandler(formatter, repo)))
    defer server.Close()

    body := []byte("this is not valid json")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
    }
    res, err := client.Do(req)
    if err != nil {
        t.Errorf("Error in POST to createMatchHandler: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending invalid JSON should result in a bad request from server.")
    }
}

func TestPostPostHandlerNotPost(t *testing.T) {
    repo := &repoTest{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postPostHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"test\":\"Not user.\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }
}

func TestPostPostHandlerSuccess(t *testing.T) {
    repo := &repoTest{}
    repo.redis =  make(map[string]string)
    client := &http.Client{}
    seconds := time.Second * time.Duration(time.Now().Unix() - time.Now().Unix())
    repo.redisSetValue("token", "1", seconds)

    server := httptest.NewServer(http.HandlerFunc(postPostHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"group_id\":1,\n\"title\":\"test\",\n\"content\":\"this is a test\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", "token")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode == http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }

    if len(repo.posts) != 1 {
        t.Error("Expected one post")
    }

    post := repo.posts[0]
    if post.Title != "test" || post.Content != "this is a test" || post.GroupID != 1 || post.UserID != 1 {
        t.Error("post not equal to what is expected")
    }
}

func TestGetPostHandlerNoPosts(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )

    repo  := &repoTest{}

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/posts/1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusNotFound {
        t.Errorf("Expected %v; received %v", http.StatusNotFound, recorder.Code)
    }
}

func TestGetPostHandlerSuccess(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )

    post := Post{GroupID: 1, Title: "Test", Content: "This is a test"}
    post.ID = 1
    repo  := &repoTest{}
    repo.addPost(post)
    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/posts/1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("Expected %v; received %v", http.StatusOK, recorder.Code)
    }

    var postResponse Post
    err := json.Unmarshal(recorder.Body.Bytes(), &postResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  postResponse.GroupID != post.GroupID && postResponse.Title != post.Title && postResponse.Content != post.Content {
        t.Errorf("Expected post recieved to equal; received %d", postResponse)
    }
}

func TestGetPostsHandlerNoGroup(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )

    repo  := &repoTest{}

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/posts?group=1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("received %v",recorder.Code)
    }

    var postResponse []Post
    err := json.Unmarshal(recorder.Body.Bytes(), &postResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  len(postResponse) != 0 {
        t.Error("Expected no posts")
    }
}

func TestGetPostsHandlerGroupSuccess(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    post := Post{GroupID: 1, Title: "Test", Content: "This is a test"}
    post2 := Post{GroupID: 2, Title: "Test", Content: "This is a test"}

    repo  := &repoTest{}
    repo.addPost(post)
    repo.addPost(post2)

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/posts?group=1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("received %v",recorder.Code)
    }

    var postResponse []Post
    err := json.Unmarshal(recorder.Body.Bytes(), &postResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  len(postResponse) != 1 {
        t.Error("Expected one post")
    }
}

func TestGetCommentsHandlerNoComments(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )

    repo  := &repoTest{}

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/posts?group=1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("received %v",recorder.Code)
    }

    var commentResponse []Comment
    err := json.Unmarshal(recorder.Body.Bytes(), &commentResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  len(commentResponse) != 0 {
        t.Error("Expected no posts")
    }
}

func TestGetCommentsHandlerWithComments(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )

    comment := Comment{PostID: 1,  Content: "This is a test"}
    comment2 := Comment{PostID: 2, Content: "This is a test"}

    repo  := &repoTest{}
    repo.addComment(comment)
    repo.addComment(comment2)

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/comments?post=1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("received %v",recorder.Code)
    }

    var commentResponse []Comment
    err := json.Unmarshal(recorder.Body.Bytes(), &commentResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  len(commentResponse) != 1 {
        t.Error("Expected one post")
    }
}

func TestGetCommentHandlerNoComment(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    repo  := &repoTest{}

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/comments/1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusNotFound {
        t.Errorf("received %v",recorder.Code)
    }
}

func TestGetCommentHandlerSuccess(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    comment := Comment{PostID: 1,  Content: "This is a test"}
    comment.ID = 1
    repo  := &repoTest{}
    repo.addComment(comment)

    server := MakeTestServer(repo)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/comments/1", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("received %v",recorder.Code)
    }
    var commentResponse Comment
    err := json.Unmarshal(recorder.Body.Bytes(), &commentResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if  commentResponse.Content != comment.Content && commentResponse.PostID != comment.PostID {
        t.Error("Not correct post")
    }
}

func TestPostCommentHandlerInvalidJSON(t *testing.T) {
    repo := &repoTest{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postCommentHandler(formatter, repo)))
    defer server.Close()

    body := []byte("this is not valid json")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
    }
    res, err := client.Do(req)
    if err != nil {
        t.Errorf("Error in POST to createMatchHandler: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending invalid JSON should result in a bad request from server.")
    }
}

func TestPostCommentHandlerNotComment(t *testing.T) {
    repo := &repoTest{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postCommentHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"test\":\"Not comment.\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }
}

func TestPostCommentHandlerSuccess(t *testing.T) {
    repo := &repoTest{}
    repo.redis =  make(map[string]string)
    client := &http.Client{}
    seconds := time.Second * time.Duration(time.Now().Unix() - time.Now().Unix())
    repo.redisSetValue("token", "1", seconds)

    server := httptest.NewServer(http.HandlerFunc(postCommentHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"post_id\":1,\n\"content\":\"this is a test\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", "token")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode == http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }

    if len(repo.comments) != 1 {
        t.Error("Expected one comment")
    }

    comment := repo.comments[0]
    if comment.Content != "this is a test" || comment.PostID != 1 || comment.UserID != 1 {
        t.Error("Comment not equal to what is expected")
    }
}

func MakeTestServer(repository *repoTest) *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter, repository)
	server.UseHandler(mx)
	return server
}
