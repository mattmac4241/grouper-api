package service

import (
    "bytes"
    "encoding/json"
    "net/http"
    "fmt"
    "os"
)

type serviceCatalogClient interface {
    publishService(name, url string)
}

//CatalogWebClient allows communication to registry
type CatalogWebClient struct {
    RootURL string
}

//PublishService publishes name and url to service
func (client CatalogWebClient) PublishService() {
    name := "api"
    url := os.Getenv("URL")
    httpclient := &http.Client{}
    service := Service{Name: name, URL: url}
    marshalled, _ := json.Marshal(service)
    req, _ := http.NewRequest("POST", client.RootURL, bytes.NewBuffer(marshalled))
    fmt.Println(req.URL)
    _, err := httpclient.Do(req)
    if err != nil {
        fmt.Println(err.Error())
    }
}
