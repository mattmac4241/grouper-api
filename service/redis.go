package service

import (
    "fmt"
    "gopkg.in/redis.v4"
)

//REDIS client to be shared throughout service
var REDIS *redis.Client


//InitRedisClient returns a redis client
func InitRedisClient(address, password string) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     address,
        Password: password, // no password set
        DB:       0,
    })
    pong, err := client.Ping().Result()
    fmt.Println(pong, err)
    if err != nil {
        fmt.Println(err)
        return client, err
    }
    return client, nil
}
