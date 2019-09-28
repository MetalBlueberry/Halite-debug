package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Path("/visor/{game}").HandlerFunc(serveGamePage)

	prefix := router.PathPrefix("/img/{game}/{turn:[0-9]+}").Subrouter()
	prefix.Methods("POST").HandlerFunc(postImgHandler)
	prefix.Methods("GET").HandlerFunc(getImgHandler)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8888",
		// Gfmtood practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Start")
	log.Fatal(srv.ListenAndServe())
}

func serveGamePage(w http.ResponseWriter, r *http.Request) {
	log.Println("handle GamePage")
	v := mux.Vars(r)
	gameName := v["game"]

	game := GetGame(gameName)

	w.Header().Add("Content-Type", "text/html")
	tpl.Execute(w, game)
}

func getImgHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("handle get")
	v := mux.Vars(r)
	gameName := v["game"]
	turnNumber, _ := strconv.Atoi(v["turn"])

	turn := GetGame(gameName).GetTurn(turnNumber)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	w.Write(turn.Svg())
}

func postImgHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("handle post")
	v := mux.Vars(r)
	gameName := v["game"]
	turnNumber, _ := strconv.Atoi(v["turn"])

	actions := []Action{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&actions)
	//log.Printf("%#v", actions)

	game := GetGame(gameName)

	turn := game.GetTurn(turnNumber)
	canvas := NewHaliteCanvas(turn)

	for _, action := range actions {
		method, ok := action["Method"]
		if !ok {
			log.Print("Method field not present")
			continue
		}
		switch method {
		case "Circle":
			x, ok := action["X"].(float64)
			if !ok {
				log.Printf("Expected %T, got %T", x, action["X"])
				continue
			}
			y, ok := action["Y"].(float64)
			if !ok {
				log.Printf("Expected %T, got %T", y, action["Y"])
				continue
			}
			r, ok := action["R"].(float64)
			if !ok {
				log.Printf("Expected %T, got %T", r, action["R"])
				continue
			}
			classInterface, ok := action["Class"].([]interface{})
			if !ok {
				log.Printf("Expected %T, got %T", classInterface, action["Class"])
				continue
			}
			class := make([]string, len(classInterface))
			for i, item := range classInterface {
				class[i], ok = item.(string)
				if !ok {
					log.Printf("Expected %T, got %T", class[i], item)
				}
			}
			canvas.Entity(x, y, r, class)
		}
	}
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
