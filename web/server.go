package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const METABASE_SITE_URL = "http://34.170.245.114:3000"

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

	//claims := CustomClaims{
	//	Resource: map[string]int{"dashboard": 1},
	//	Params:   map[string]interface{}{},
	//	RegisteredClaims: jwt.RegisteredClaims{
	//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
	//	},
	//}
	//
	//if METABASE_SECRET_KEY == "" {
	//	panic("METABASE_SECRET_KEY not set")
	//}
	//fmt.Println("METABASE_SECRET_KEY: " + METABASE_SECRET_KEY)
	//// Create a new token object
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//secretKey := []byte(METABASE_SECRET_KEY)
	//
	//tokenString, err := token.SignedString(secretKey)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	fmt.Printf("error: %v\n", err)
	//	return
	//}

	tmpl, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error: %v\n", err)
		return
	}

	//iframeUrl := METABASE_SITE_URL + "/embed/dashboard/" + tokenString + "#bordered=true&titled=true"
	//
	tmplData := PageData{
		//IFrameUrl: iframeUrl,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error: %v\n", err)
		return
	}
}

const tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Hivemaper Dashboard</title>
</head>
<body>
<iframe
	style="position:fixed; top:0; left:0; bottom:0; right:0; width:100%; height:100%; border:none; margin:0; padding:0; overflow:hidden; z-index:999999;"
    src="http://metabase.streamingfast.io/public/dashboard/3e029abe-66bf-4cad-895c-e39922f03927"
    frameborder="0"
    width="100%"
    height="100%"
	allow="fullscreen"
></iframe>
</body>
</html>
`

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
