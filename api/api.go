package api

import (
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
	Data string
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
	log.Println(wh.Data)
}

func dockerTest() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
}

func Handlers() *mux.Router {
	dockerTest()

	r := mux.NewRouter()

	r.HandleFunc("/apps", createAppHandler).Methods("POST")
	r.HandleFunc("/botd", webhookHandler).Methods("POST")
	r.HandleFunc("/apps", listAppsHandler).Methods("GET")
	r.HandleFunc("/version", versionHandler).Methods("GET")

	return r
}
