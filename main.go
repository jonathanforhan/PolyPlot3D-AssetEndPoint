package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var whitelist = []string{
	"https://poly-plot-3d.netlify.app",
}

/* Get the request's IP */
func readIp(r *http.Request) string {
	ipaddr := r.Header.Get("X-Real-Ip")
	if ipaddr == "" {
		ipaddr = r.Header.Get("X-Forwarded-For")
	}
	if ipaddr == "" {
		ipaddr = r.RemoteAddr
	}
	if ipaddr == "" {
		ipaddr = "UNKNOWN"
	}

	return ipaddr
}

/* if api_key matches .env api key we allow request */
func cors(w *http.ResponseWriter, r *http.Request) {
	origin := (*r).Header.Get("Origin")
	if strings.HasPrefix(origin, "http://localhost:") {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	} else if slices.Contains(whitelist, origin) {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	}
}

/* GET "/" */
func getRoot(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	log.Printf("GET %s %s", r.URL.String(), readIp(r))
	w.Write([]byte("PolyPlot3D Asset Import Endpoint"))
}

/* GET "/import" */
func getImport(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	asset := r.URL.Query().Get("asset")
	ft := r.URL.Query().Get("ft")
	status := "GET-INVALID"

	if len(asset) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing asset name in request"))
	} else if len(ft) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing filetype in request"))
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./assets/%s/%s.%s", asset, asset, ft))

		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf("asset %s.%s not found", asset, ft)))
		} else {
			/* ONLY VALID BRANCH */
			status = "GET"
			w.Write(data)
		}
	}
	log.Printf("%s %s %s", status, r.URL.String(), readIp(r))
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC ", err)
			main()
		}
	}()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	logfile, err := os.Create(fmt.Sprintf("log/server-log-%s.log", time.Now().UTC().Format(time.DateTime)))
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/import", getImport)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	fmt.Printf("Server Listening on Port " + port + "\n")

	// err = http.ListenAndServeTLS(":"+port, "server.crt", "server.key", mux)
	err = http.ListenAndServe(":"+port, mux)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("SERVER-CLOSED ")
	} else if err != nil {
		log.Printf("ERROR %s\n", err)
	}
}
