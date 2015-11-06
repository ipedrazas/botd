package api

import (
	"encoding/json"
	// "os"
	// "errors"
	// "fmt
	// "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"log"
	"net/http"
)

type Version struct {
	Apiversion string `json: version`
}

type Repository struct {
	Name      string `json: name`
	Namespace string `json: namespace`
}

type Data struct {
	Pushed float64 `json:"pushed_at,string"`
	Pusher string  `json: pusher`
}

type WebhookJson struct {
	PushData Data       `json:"push_data"`
	Repo     Repository `json:"repository"`
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

func hooksHandler(w http.ResponseWriter, req *http.Request) {
	hooks := getHooks(redisClient)

	shooks := "{["
	for i := 0; i < len(hooks); i++ {
		if i > 0 {
			shooks = shooks + ", "
		}
		shooks = shooks + hooks[i]
	}
	shooks = shooks + "]}"

	jhooks, err := json.Marshal(hooks)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jhooks)

}

func webhookHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))

	key, err := setHook(redisClient, string(body))

	if err != nil {
		panic(err)
	}

	log.Println(key)
	log.Println("...h")

	// var wh WebhookJson
	// err = json.Unmarshal([]byte(body), &wh)
	// if err != nil {
	// 	panic(err)
	// }

	// repo := wh.Repo
	// data := wh.PushData
	// log.Println("")
	// log.Println(repo.Name)
	// log.Println(repo.Namespace)
	// log.Println(data.Pusher)
	// log.Println(data.Pushed)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// func dockerBuild() {
// 	endpoint := "unix:///var/run/docker.sock"
// 	client, err := docker.NewClient(endpoint)
// 	check(err)
// 	dockerfile := "Dockerfile"
// 	dir := "/home/ivan/go/src/github.com/ipedrazas/botd"
// 	var output bytes.Buffer
// 	opts := docker.BuildImageOptions{
// 		Name:         "botd",
// 		Dockerfile:   dockerfile,
// 		ContextDir:   dir,
// 		OutputStream: &output,
// 		Remote:       "github.com/ipedrazas/botd",
// 	}
// 	log.Println("Build options")
// 	if err := client.BuildImage(opts); err != nil {
// 		log.Fatal(err)
// 	}

// }

var redisClient *Redis

type Redis struct {
	*redis.Client
}

func Handlers() *mux.Router {
	// dockerBuild()

	redisClient = &Redis{
		redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}

	r := mux.NewRouter()

	r.HandleFunc("/botd", webhookHandler).Methods("POST")
	r.HandleFunc("/hooks", hooksHandler).Methods("GET")
	r.HandleFunc("/hooks", webhookHandler).Methods("POST")
	r.HandleFunc("/version", versionHandler).Methods("GET")

	return r
}
