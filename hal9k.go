package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ralreegorganon/cylon"
)

func main() {
	match := getMatch()
	port := getListenPort()

	remoteRoot := "http://localhost:3008"
	localRoot := "http://localhost:" + port
	ch := make(chan bool)

	h := &hal9k{}

	server := cylon.NewServer(h, remoteRoot, localRoot, ch)
	r, err := server.CreateRouter()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", r)

	go func() {
		log.Printf("AI started at http://%s\n", "localhost:"+port)
		err := http.ListenAndServe("localhost:"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = server.Join(match)
	if err != nil {
		ch <- false
	}

	<-ch
}

func getMatch() string {
	var dat map[string]interface{}
	f, err := ioutil.ReadFile("match.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
	json.Unmarshal(f, &dat)
	match := dat["match"].(string)
	return match
}

func getListenPort() string {
	rand.Seed(time.Now().Unix())
	port := strconv.Itoa(rand.Intn(999) + 9000)
	return port
}

type hal9k struct {
}

func (r *hal9k) Think(s *cylon.RobotState) *cylon.RobotCommands {
	commands := &cylon.RobotCommands{
		Turn:       0,
		TurnGun:    0,
		TurnRadar:  1,
		Accelerate: 0,
		Fire:       0,
	}
	return commands
}
