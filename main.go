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
	"github.com/metalblueberry/Halite-debug/internal/actions"
)

func main() {
	router := mux.NewRouter()

	router.Path("/visor/{game}").HandlerFunc(serveGamePage)

	prefix := router.PathPrefix("/img/{game}/{turn:[0-9]+}").Subrouter()
	prefix.Methods("POST").HandlerFunc(postImgHandler)
	prefix.Methods("GET").HandlerFunc(getImgHandler)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Handler: logging(router),
		Addr:    "127.0.0.1:8888",
		// Gfmtood practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Start")
	log.Fatal(srv.ListenAndServe())
}

func serveGamePage(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	gameName := v["game"]

	game := GetGame(gameName)

	w.Header().Add("Content-Type", "text/html")
	tpl.Execute(w, game)
}

func getImgHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	gameName := v["game"]
	turnNumber, _ := strconv.Atoi(v["turn"])

	turn := GetGame(gameName).GetTurn(turnNumber)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	w.Write(turn.Svg())
}

func postImgHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	gameName := v["game"]
	turnNumber, _ := strconv.Atoi(v["turn"])

	actions := []actions.Action{}
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
			x, y, r := action.Circle()
			class := action.Class()
			canvas.Circle(x, y, r, class)
		case "Line":
			x1, y1, x2, y2 := action.Line()
			class := action.Class()
			canvas.Line(x1, y1, x2, y2, class)
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

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s:%s in %s", r.Method, r.URL.Path, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}
