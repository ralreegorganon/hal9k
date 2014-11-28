package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	rand.Seed(time.Now().Unix())
	port := strconv.Itoa(rand.Intn(999) + 9000)

	router, err := createRouter()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Handle("/", router)

	go func() {
		log.Printf("AI started at http://%s\n", "localhost:"+port)
		err = http.ListenAndServe("localhost:"+port, nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	j := &joinMatchRequestMessage{
		Endpoint: "http://localhost:" + port,
		Match:    "68267d65-6f63-4c18-8afa-dce5c91e0e73",
	}

	log.Printf("joining match %+v\n", j)
	js, _ := json.Marshal(j)
	res, err := http.Post("http://localhost:3008/join", "application/json", bytes.NewBuffer(js))
	if err != nil {
		log.Println(err)
	}
	log.Println(res.StatusCode)

	/*
		s := &startMatchRequestMessage{
			Match: "68267d65-6f63-4c18-8afa-dce5c91e0e73",
		}

		log.Printf("starting match %+v\n", s)
		js, _ = json.Marshal(s)
		log.Println("%v", string(js))
				res, err = http.Post("http://localhost:3008/start", "application/json", bytes.NewBuffer(js))
				if err != nil {
					log.Println(err)
				}
				log.Println(res.StatusCode)
	*/

	ch := make(chan bool)
	<-ch
}

type joinMatchRequestMessage struct {
	Endpoint string `json:"endpoint"`
	Match    string `json:"match"`
}

type startMatchRequestMessage struct {
	Match string `json:"match"`
}

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	m := map[string]map[string]HttpApiFunc{
		"GET": {
			"/status": Status,
		},
		"POST": {
			"/status": Status,
			"/think":  Think,
		},
	}

	for method, routes := range m {
		for route, handler := range routes {
			localRoute := route
			localHandler := handler
			localMethod := method
			f := makeHttpHandler(localMethod, localRoute, localHandler)

			r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
		}
	}

	return r, nil
}

func makeHttpHandler(localMethod string, localRoute string, handlerFunc HttpApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeCorsHeaders(w, r)
		if err := handlerFunc(w, r, mux.Vars(r)); err != nil {
			httpError(w, err)
		}
	}
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) error

func writeJSON(w http.ResponseWriter, code int, thing interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	val, err := json.Marshal(thing)
	w.Write(val)
	return err
}

func httpError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	if err != nil {
		http.Error(w, err.Error(), statusCode)
	}
}

func Status(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func Think(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	decoder := json.NewDecoder(r.Body)
	state := &RobotState{}

	err := decoder.Decode(state)
	if err != nil {
		log.Println(err)
		return err
	}

	// THINK HERE
	log.Printf("%+v\n", state)

	commands := &RobotCommands{
		Turn:       20,
		TurnGun:    10,
		TurnRadar:  5,
		Accelerate: 1,
		Fire:       1,
	}

	log.Printf("sending %+v\n", commands)

	err = writeJSON(w, http.StatusOK, commands)
	if err != nil {
		log.Println(err)
	}

	return nil
}

type RobotState struct {
	Position     Point   `json:"position"`
	Heading      float64 `json:"heading"`
	GunHeading   float64 `json:"gunHeading"`
	RadarHeading float64 `json:"radarHeading"`
	Velocity     float64 `json:"velocity"`
	Heat         float64 `json:"heat"`
	Health       float64 `json:"health"`
	Alive        bool    `json:"alive"`
}

type RobotCommands struct {
	Turn       float64 `json:"turn"`
	TurnGun    float64 `json:"turnGun"`
	TurnRadar  float64 `json:"turnRadar"`
	Accelerate float64 `json:"accelerate"`
	Fire       float64 `json:"fire"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
