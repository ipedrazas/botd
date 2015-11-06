package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type WebhookJson struct {
	pushed    uint64 `json:pushed_at`
	pusher    string `json:pusher`
	name      string `json:name`
	namespace string `json:namespace`
}

type App struct {
	Id      uint32 `json:"id"`
	Appname string `json:"name"`
	Apptype string `json:"type"`
}

type AppParams struct {
	Appname string `json:"name"`
	Apptype string `json:"type"`
}

type Version struct {
	Apiversion string `json: version`
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
		Id:      appIdCounter,
		Appname: p.Appname,
		Apptype: p.Apptype,
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
	v := Version{
		Apiversion: "0.0.1",
	}

	version, err := json.Marshal(v)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(version)
}

func webhookHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var wh WebhookJson
	err = json.Unmarshal(body, &wh)
	if err != nil {
		panic(err)
	}
	log.Println(wh.name)
	log.Println(wh.namespace)
	log.Println(wh.pusher)
	log.Println(wh.pushed)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func dockerBuild() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	check(err)
	dockerfile := "Dockerfile"
	dir := "/home/ivan/go/src/github.com/ipedrazas/botd"
	var output bytes.Buffer
	opts := docker.BuildImageOptions{
		Name:         "botd",
		Dockerfile:   dockerfile,
		ContextDir:   dir,
		OutputStream: &output,
		Remote:       "github.com/ipedrazas/botd",
	}
	log.Println("Build options")
	if err := client.BuildImage(opts); err != nil {
		log.Fatal(err)
	}

}

func Handlers() *mux.Router {
	dockerBuild()

	r := mux.NewRouter()

	r.HandleFunc("/apps", createAppHandler).Methods("POST")
	r.HandleFunc("/botd", webhookHandler).Methods("POST")
	r.HandleFunc("/apps", listAppsHandler).Methods("GET")
	r.HandleFunc("/version", versionHandler).Methods("GET")

	return r
}
