package service

import (
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

func getGroupsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        groups := Group{}
        DB.Find(&groups)
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
        formatter.JSON(w, http.StatusCreated, "Group succesfully created.")
    }
}
