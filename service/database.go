package service

import (
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

//DB shared connection through service
var DB *gorm.DB

//InitDatabase setsup the database
func InitDatabase(host, user, dbname, password string) *gorm.DB{
    dbFormat := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
    db, _ := gorm.Open("postgres", dbFormat)
    return db
}

//CloseDatabase closes the current database
func CloseDatabase() {
    DB.Close()
}

//CreateModels inits the database with the models
func CreateModels() {
    DB.CreateTable(&Group{}, &Post{}, &Comment{}, &GroupMember{}, &GroupAdmin{})
}

//MigrateModels updates the models in the database
func MigrateModels() {
    DB.AutoMigrate(&Group{}, &Post{}, &Comment{}, &GroupMember{}, &GroupAdmin{})
}

//DropModels deletes the models from the database
func DropModels() {
    DB.DropTable(&Group{}, &Post{}, &Comment{}, &GroupMember{}, &GroupAdmin{})
}
