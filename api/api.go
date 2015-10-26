package api

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

type App struct {
    Id          uint32 `json:"id"`
    Appname     string `json:"name"`
    Apptype     string `json:"type"`
}

type AppParams struct {
    Appname     string `json:"name"`
    Apptype     string `json:"type"`
}

var appIdCounter uint32 = 0

var appStore = []App{}

func createAppHandler(w http.ResponseWriter, r *http.Request) {
    p := AppParams{}

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = json.Unmarshal(body, &p)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = validateUniqueness(p.Appname)

    if err != nil {
        fmt.Printf("Error: %s\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    a := App{
        Id:           appIdCounter,
        Appname:      p.Appname,
        Apptype:      p.Apptype,
    }

    appStore = append(appStore, a)

    appIdCounter += 1

    w.WriteHeader(http.StatusCreated)
}

func validateUniqueness(appname string) error {
    for _, a := range appStore {
        if a.Appname == appname {
            return errors.New("Appname is already used")
        }
    }

    return nil
}

func listAppsHandler(w http.ResponseWriter, r *http.Request) {
    apps, err := json.Marshal(appStore)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(apps)
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Version 0.0.1"))
}

func Handlers() *mux.Router {
    r := mux.NewRouter()

    r.HandleFunc("/apps", createAppHandler).Methods("POST")
    r.HandleFunc("/apps", listAppsHandler).Methods("GET")
    r.HandleFunc("/version", versionHandler).Methods("GET")

    return r
}
