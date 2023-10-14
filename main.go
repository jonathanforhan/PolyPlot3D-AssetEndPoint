package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func readIp(r *http.Request) (string, error) {
	ipaddr := r.Header.Get("X-Real-Ip")
	if ipaddr == "" {
		ipaddr = r.Header.Get("X-Forwarded-For")
	}
	if ipaddr == "" {
		ipaddr = r.RemoteAddr
	}
    if ipaddr == "" {
        return ipaddr, errors.New("could not parse ip")
    }

    return ipaddr, nil
}

func cors(w *http.ResponseWriter, r *http.Request) {
    ipaddr, err := readIp(r)
    if err != nil {
        return
    }
    if strings.HasPrefix(ipaddr, "localhost:") || strings.HasPrefix(ipaddr, "[::1]:"){
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
	} else if strings.HasPrefix(ipaddr, "poly-plot-3d") {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)
    ipaddr, err := readIp(r);
    if err != nil { ipaddr = "UNKNOWN" }

	log.Printf("GET %s %s", r.URL.String(), ipaddr)
	io.WriteString(w, "PolyPlot3D Asset Import Endpoint")
}

func getImport(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)
    ipaddr, err := readIp(r);
    if err != nil { ipaddr = "UNKNOWN" }

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
	log.Printf("%s %s %s", status, r.URL.String(), ipaddr)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC ", err)
			main()
		}
	}()

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

	err = http.ListenAndServe(":8000", mux)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("SERVER-CLOSED ")
	} else if err != nil {
		log.Printf("ERROR %s\n", err)
	}
}
