package main

import "html/template"

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.New("Visor").Parse(`
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
	<svg class="fillscreen" >
		<g class="fillscreen" id="scene"></g>
	</svg>
    <script src='https://unpkg.com/panzoom@8.4.0/dist/panzoom.min.js'></script>
    <script>
	var turn = 1;
	var game = {{ .Name }}

	console.log("start")

	function loadTurn(callback){
		httpGetAsync("/img/" + game + "/" + turn, function(text){
			console.log("image requested")
			document.getElementById('scene').innerHTML = text
			if(callback && typeof callback === "function") {
				callback()
			}
		})
		document.getElementById('title').innerHTML = game + " turn " + turn
	}

	function initController(){
		console.log("Controller initialization")
		var area = document.getElementById('scene')
		window.pz = panzoom(area, {
			autocenter: true,
			bounds: true,
			  filterKey: function(/* e, dx, dy, dz */) {
				// don't let panzoom handle this event:
				return true;
			  }
		})
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

	loadTurn(function(){
		initController()
	})

	function keyboardNavigation(e){
		switch(e.keyCode) {
		  case 37:
			if (turn > 1) {
				turn = turn-1
				loadTurn()
			}
			break;
		  case 39:
			if (turn < 300) {
				turn = turn+1
				loadTurn()
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
}
