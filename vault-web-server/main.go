package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pashpashpash/vault/serverutil"

	"github.com/pashpashpash/vault/vault-web-server/postapi"

	openai "github.com/sashabaranov/go-openai"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const (
	NegroniLogFmt = `{{.StartTime}} | {{.Status}} | {{.Duration}}
          {{.Method}} {{.Path}}`
	NegroniDateFmt = time.Stamp
)

var (
	debugSite = flag.Bool(
		"debug", false, "debug site")
	port = flag.String(
		"port", "80", "server port")
	siteConfig = map[string]string{
		"DEBUG_SITE": "false",
	}
)

func main() {
	// Parse command line flags + override defaults
	flag.Parse()
	siteConfig["DEBUG_SITE"] = strconv.FormatBool(*debugSite)
	rand.Seed(time.Now().UnixNano())

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if len(openaiApiKey) == 0 {
		log.Fatalln("MISSING OPENAI API KEY ENV VARIABLE")
	}
	openaiClient := openai.NewClient(openaiApiKey)

	pineconeApiKey := os.Getenv("PINECONE_API_KEY")
	if len(pineconeApiKey) == 0 {
		log.Fatalln("MISSING PINECONE API KEY ENV VARIABLE")
	}

	pineconeApiEndpoint := os.Getenv("PINECONE_API_ENDPOINT")
	if len(pineconeApiEndpoint) == 0 {
		log.Fatalln("MISSING PINECONE API ENDPOINT ENV VARIABLE")
	}

	// Initialize modules
	postapi.Run(openaiClient, pineconeApiKey, pineconeApiEndpoint)

	// Configure main web server
	server := negroni.New()
	server.Use(negroni.NewRecovery())
	l := negroni.NewLogger()
	l.SetFormat(NegroniLogFmt)
	l.SetDateFormat(NegroniDateFmt)
	server.Use(l)
	mx := mux.NewRouter()

	// Path Routing Rules: [POST]
	mx.HandleFunc("/api/questions", postapi.QuestionHandler).Methods("POST")
	mx.HandleFunc("/upload", postapi.UploadHandler).Methods("POST")

	// Path Routing Rules: Static Handlers
	mx.HandleFunc("/github", StaticRedirectHandler("https://github.com/pashpashpash/vault"))
	mx.PathPrefix("/").Handler(ReactFileServer(http.Dir(serverutil.WebAbs(""))))

	// Start up web server
	server.UseHandler(mx)
	server.Run(":" + *port)
}

// Forwards all traffic to React, except basic file serving
func ReactFileServer(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(serverutil.WebAbs(r.URL.Path)); os.IsNotExist(err) {
			// Do nothing, and let the request pass through
		} else {
			fsh.ServeHTTP(w, r)
		}
	})
}

func StaticRedirectHandler(to string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r,
to, http.StatusFound)
}
}
