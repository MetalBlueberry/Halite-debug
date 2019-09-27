package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	svg "github.com/ajstarks/svgo/float"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.Path("/visor/{game}").HandlerFunc(serveGamePage)
	prefix := router.PathPrefix("/img/{game}/{id:[0-9]+}").Subrouter()
	prefix.Methods("POST").HandlerFunc(postHandler)
	prefix.Methods("GET").HandlerFunc(getHandler)
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

type Turn struct {
	bytes.Buffer
}

func NewTurn() *Turn {
	return &Turn{
		Buffer: bytes.Buffer{},
	}
}

func (t Turn) Svg() []byte {
	buf := &bytes.Buffer{}
	canvas := svg.New(buf)
	//svginitfmt := `<svg width="%.*f%s" height="%.*f%s">`
	//buf.WriteString(fmt.Sprintf(svginitfmt, 2, 500.0, "", 2, 500.0, ""))
	//buf.WriteString("<svg >")
	canvas.Startraw("class=\"fillscreen\"")
	canvas.Gid("scene")
	buf.Write(t.Bytes())
	canvas.Gend()
	canvas.End()
	return buf.Bytes()
}

var turns = make(map[string]Turn)

func serveGamePage(w http.ResponseWriter, r *http.Request) {
	log.Println("handle GamePage")
	v := mux.Vars(r)
	game := v["game"]
	tpl, err := template.New("Visor").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv='content-type' content='text/html; charset=utf-8' />
    <meta name='viewport' content='width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge' >
    <META NAME='Description' content='Pan and zoom svg elements demo '>
    <meta name='keywords' content='svg, pan, zoom' />
    <meta name='author' content='Andrei Kashcha'>
    <meta name='title' content='SVG panzoom demo' />
	<link rel="stylesheet" href="/styles.css">
    <title id="title" ></title>
  </head>
  <body onkeydown="keyboardNavigation(event)">
	<div class="fillscreen" id="viewport"></div>
    <script src='https://unpkg.com/panzoom@8.4.0/dist/panzoom.min.js'></script>
    <script>
	var turn = 1;
	var game = {{ . }}

	console.log("start")

	function initController(){
		console.log("Controller initialization")
		httpGetAsync("/img/" + game + "/" + turn, function(text){
			console.log("image requested")
			document.getElementById('viewport').innerHTML = text
			var area = document.getElementById('scene')
			window.pz = panzoom(area, {
				autocenter: true,
				bounds: true,
				  filterKey: function(/* e, dx, dy, dz */) {
					// don't let panzoom handle this event:
					return true;
				  }
			})
		})
		document.getElementById('title').innerHTML = game + " turn " + turn
	}

	function httpGetAsync(theUrl, callback)
	{
		var xmlHttp = new XMLHttpRequest();
		xmlHttp.onreadystatechange = function() { 
			if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
				callback(xmlHttp.responseText);
		}
		xmlHttp.open("GET", theUrl, true); // true for asynchronous 
		xmlHttp.send(null);
	}
	initController()

	function keyboardNavigation(e){
		switch(e.keyCode) {
		  case 37:
			if (turn > 1) {
				turn = turn-1
				initController()
			}
			break;
		  case 39:
			if (turn < 300) {
				turn = turn+1
				initController()
			}
			break;
		} 
	}

    </script>
  </body>
</html>
	`)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "text/html")
	tpl.Execute(w, game)

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("handle get")
	v := mux.Vars(r)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	t := turns[v["id"]]
	w.Write(t.Svg())
}

type Action map[string]interface{}

//Method string
//Params []json.Token

type HaliteCanvas struct {
	*svg.SVG
}

func NewHaliteCanvas(w io.Writer) *HaliteCanvas {
	return &HaliteCanvas{
		SVG: svg.New(w),
	}
}

var emptyclose = "/>\n"

func (canvas HaliteCanvas) Planet(x float64, y float64, r float64, ownerID string, s ...string) {
	d := canvas.Decimals
	fmt.Fprintf(canvas.Writer, `<circle class="planet %s" cx="%.*f" cy="%.*f" r="%.*f" %s`, ownerID, d, x, d, y, d, r, endstyle(s, emptyclose))
}
func (canvas HaliteCanvas) Entity(x float64, y float64, r float64, class []string, s ...string) {
	d := canvas.Decimals
	fmt.Fprintf(canvas.Writer, `<circle class="planet %s" cx="%.*f" cy="%.*f" r="%.*f" %s`, strings.Join(class, " "), d, x, d, y, d, r, endstyle(s, emptyclose))
}

// endstyle modifies an SVG object, with either a series of name="value" pairs,
// or a single string containing a style
func endstyle(s []string, endtag string) string {
	if len(s) > 0 {
		nv := ""
		for i := 0; i < len(s); i++ {
			if strings.Index(s[i], "=") > 0 {
				nv += (s[i]) + " "
			} else {
				nv += style(s[i]) + " "
			}
		}
		return nv + endtag
	}
	return endtag

}
func style(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`style="%s"`, s)
	}
	return s
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("handle post")
	v := mux.Vars(r)
	t := turns[v["id"]]

	actions := []Action{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&actions)
	log.Printf("%#v", actions)

	canvas := NewHaliteCanvas(&t)

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

	turns[v["id"]] = t
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
