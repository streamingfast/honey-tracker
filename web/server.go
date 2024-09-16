package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const METABASE_SITE_URL = "http://metabase.streamingfast.io"

var METABASE_SECRET_KEY = os.Getenv("SECRET_KEY")

type CustomClaims struct {
	Resource map[string]int         `json:"resource"`
	Params   map[string]interface{} `json:"params"`
	jwt.RegisteredClaims
}

type PageData struct {
	IFrameUrl string
}

func handler(w http.ResponseWriter, r *http.Request) {
	dashboard := 1
	dashboardString := r.URL.Query().Get("dashboard")
	if dashboardString != "" {
		d, err := strconv.Atoi(dashboardString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Printf("error: %v\n", err)
			return
		}
		dashboard = d

	}

	claims := CustomClaims{
		Resource: map[string]int{"dashboard": dashboard},
		Params:   map[string]interface{}{},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}
	//
	if METABASE_SECRET_KEY == "" {
		panic("METABASE_SECRET_KEY not set")
	}
	fmt.Println("METABASE_SECRET_KEY: " + METABASE_SECRET_KEY)
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(METABASE_SECRET_KEY)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error: %v\n", err)
		return
	}

	tmpl, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error: %v\n", err)
		return
	}

	iframeURL := METABASE_SITE_URL + "/embed/dashboard/" + tokenString + "#bordered=false&titled=false"
	fmt.Println("iframeURL: " + iframeURL)

	tmplData := PageData{
		IFrameUrl: iframeURL,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("Done rendering")
}

const tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Hivemaper Dashboard</title>
	<style>
		body, html {width: 100%; height: 100%; margin: 0; padding: 0}
		.row-container {display: flex; width: 100%; height: 100%; flex-direction: column; background-color: white; overflow: hidden;}
		.first-row {}
		.second-row { flex-grow: 1; border: none; margin: 0; padding: 0; }

        .round-button {
            padding: 10px 20px;
            margin: 10px;
            border-radius: 25px;
            border: none;
            background-color: #679BDF; /* Green background */
            color: white;
            cursor: pointer;
            font-size: 16px;
            text-decoration: none;
            display: inline-block;
        }
        
        .round-button:hover {
            background-color: #45a049; /* Darker green on hover */
        }

	</style>
</head>
<body>

<div class="row-container">
  <div class="first-row">
	<div style="display: flex; align-items: center; padding-left: 40px; padding-top: 20px; padding-bottom: 20px;">
		<!-- img src="https://cdn.prod.website-files.com/649aef4c9068091b737a9baf/64c3e6ca809ba2f4f2d76359_streamingfast-p-500.png" width=150 height=63></img -->
		<span style="font-family:sans-serif; font-weight: bold; font-size: 26px; color:#4F5671;">Hivemapper - Powered by The Graph Substreams</span>
	</div>
	<div style="padding-left: 25px;">
    <button class="round-button" onclick="window.location.href='./'">Overview</button>
    <button class="round-button" onclick="window.location.href='./?dashboard=33'">Fleet Dashboard</button>
    <button class="round-button" onclick="window.location.href='./?dashboard=34'">Driver Dashboard</button>
	</div>

  </div>
	<iframe	
		class="second-row"
    	src="{{ .IFrameUrl }}"
		sandbox="allow-scripts allow-same-origin"
    	frameborder="0"
    	width="100%"
    	height="100%"
		allow="fullscreen"
	></iframe>
</div>

</body>
</html>
`

//style="position:fixed; top:0px; left:0; bottom:100; right:0; width:100%; height:100%; border:none; margin-top:100px; padding:0; overflow:hidden; z-index:999999;"
//const tmpl = `
//<!DOCTYPE html>
//<html lang="en">
//<head>
//    <meta charset="UTF-8">
//    <title>Hivemaper Dashboard</title>
//</head>
//<body>
//<iframe
//    src="{{.IFrameUrl}}"
//    frameborder="0"
//    width="90%"
//    height="90%"
//    allowtransparency
//></iframe></body>
//</html>
//`

func ServeHttp() {
	http.HandleFunc("/", handler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
