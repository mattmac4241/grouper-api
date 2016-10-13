package service

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type authClient interface {
    getUserIDFromToken(key string) (token Token, err error)
}

type authtWebClient struct {
    rootURL string
}

func (client authtWebClient) getUserIDFromToken(key string) (token Token, err error) {
    httpclient := &http.Client{}

    tokenURL := fmt.Sprintf("%s/%s", client.rootURL, key)
    req, _ := http.NewRequest("GET", tokenURL, nil)

    resp, err := httpclient.Do(req)

    if err != nil {
        fmt.Printf("Errored when sending request to the server: %s\n", err.Error())
        return
    }

    defer resp.Body.Close()
    payload, _ := ioutil.ReadAll(resp.Body)
    err = json.Unmarshal(payload, &token)
    if err != nil {
        fmt.Println("Failed to unmarshal server response")
    }

    return token, err
}
