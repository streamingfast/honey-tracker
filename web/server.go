package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const METABASE_SITE_URL = "https://metabase.streamingfast.io"

const HIVEMPPER_SITE_URL = "https://hivemapper.streamingfast.io"

//const HIVEMPPER_SITE_URL = "http://localhost:8080"

var METABASE_SECRET_KEY = os.Getenv("SECRET_KEY")

type CustomClaims struct {
	Resource map[string]int         `json:"resource"`
	Params   map[string]interface{} `json:"params"`
	jwt.RegisteredClaims
}

type PageData struct {
	IFrameUrl string
}

func handleIFrame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling request:", r.URL.Path)
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

	iframeURL := HIVEMPPER_SITE_URL + "/embed/dashboard/" + tokenString + "#bordered=false&titled=false"
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

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		handleIFrame(w, r)
		return
	}
	r.Host = s.metabaseUrl.Host
	s.proxy.ServeHTTP(w, r)
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
	<div style="display: flex;  flex-direction: row;  flex-wrap: wrap; align-items: center; padding-left: 40px; padding-top: 5px; padding-bottom: 5px;">
		<div style="display: flex; flex-grow:2; align-items: center;">
			<a href="https://hivemapper.com/">
				<img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRxugF0W5SMUWI-Pp0abVDJ8zWDP8UJ5ccWJQ8zzYgsgUev7maTTasAWaD123-LTPuAl30&usqp=CAU" height=50/>
			</a>
			<span style="flex-grow:2; padding-left:5px;font-family:sans-serif; font-weight: bold; font-size: 26px; color:#4F5671;">Hivemapper</span>
		</div>
		<div style="display: flex; align-items: center;">
			<span style="font-family:sans-serif; font-weight: bold; font-size: 18px; color:#4F5671; white-space: nowrap; padding-right:15px">Powered by:</span>
			<a href="https://thegraph.com/">
				<img src="https://cdn.prod.website-files.com/649aef4c9068091b737a9baf/66ad5cd70dec61e829ac393d_The%20Graph%20-%20Logo%20-%20Dark.svg"  height=40/>
			</a>
			<a href="https://www.streamingfast.io/#substreams">
				<img style="padding-top: 5px" src="https://cdn.prod.website-files.com/649aef4c9068091b737a9baf/65119d95beb11cb3b1c4583c_Substreams-Tagline-color-p-800.png" height=75/>
			</a>
		</div>
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
		sandbox="allow-scripts allow-same-origin allow-downloads"
    	frameborder="0"
    	width="100%"
    	height="100%"
		allow="fullscreen"

	></iframe>
</div>

</body>
</html>
`

// style="position:fixed; top:0px; left:0; bottom:100; right:0; width:100%; height:100%; border:none; margin-top:100px; padding:0; overflow:hidden; z-index:999999;"
// const tmpl = `
// <!DOCTYPE html>
// <html lang="en">
// <head>
//
//	<meta charset="UTF-8">
//	<title>Hivemaper Dashboard</title>
//
// </head>
// <body>
// <iframe
//
//	src="{{.IFrameUrl}}"
//	frameborder="0"
//	width="90%"
//	height="90%"
//	allowtransparency
//
// ></iframe></body>
// </html>
// `
type Server struct {
	proxy       *httputil.ReverseProxy
	metabaseUrl *url.URL
}

func (s *Server) ServeHttp() {
	url, err := url.Parse(METABASE_SITE_URL)
	if err != nil {
		log.Fatal(err)
	}
	s.metabaseUrl = url

	proxy := httputil.NewSingleHostReverseProxy(url)

	s.proxy = proxy

	http.HandleFunc("/", s.handleRoot)

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
