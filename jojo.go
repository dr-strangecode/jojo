package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	auth, ok := r.Header["Authorization"]
	if !ok {
		w.Header().Add("WWW-Authenticate", "basic realm=\"jojo\"")
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Unauthorized access to %s", r.URL)
		return false
	}
	encoded := strings.Split(auth[0], " ")
	if len(encoded) != 2 || encoded[0] != "Basic" {
		log.Printf("Strange Authorizatoion %q", auth)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded[1])
	if err != nil {
		log.Printf("Cannot decode %q: %s", auth, err)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		log.Printf("Unknown format for credentials %q", decoded)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if parts[0] == user && parts[1] == password {
		return true
	}
	return false
}

func scriptHandlerGenerator(script string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user != "" && password != "" {
			authorized := checkAuth(w, r)
			if authorized != true {
				log.Printf("Unauthorized")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// Replace "=" and "&" with a " "
		// Creates []string urlCmdArgs
		re := regexp.MustCompile("=|&")
		urlCmdArgs := strings.Split(re.ReplaceAllString(r.URL.RawQuery, " "), " ")

		// Misc Logging junk
		log.Printf("[Info] %s %s %s %s %s", r.Proto, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
		log.Printf("Args: %s", r.URL.Query())
		//log.Printf("[DEBUG] <Header>%s</Header>", r.Header)
		log.Printf("[Info] Running %s %s", script, urlCmdArgs)
		fmt.Fprintf(w, "{\"script\": \"%s\",\"arguments\": \"%s\",", script, urlCmdArgs)

		// Run the script passing in the arguments
		cmd := exec.Command(script, urlCmdArgs...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		cmdErr := cmd.Run()
		if cmdErr != nil {
			fmt.Fprintf(w, "\"error\": \"%s\",", cmdErr)
			log.Printf("[ERROR] %s", cmdErr)
		}
		fmt.Fprintf(w, "\"results\": \"%s\"}", strings.Split(out.String(), "\n"))
		log.Printf("[Info] Results: %s", out.String())
	}
}

func loadConfig(useSSL bool) {
	r := mux.NewRouter()

	if useSSL {
		log.Println("using ssl")
		r.Schemes("https")
	} else {
		r.Schemes("http")
	}

	config := yaml.ConfigFile(configFile)
	numRoutes, err := config.Count("routes")
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	for i := 0; i < numRoutes; i++ {
		url, _ := config.Get(fmt.Sprintf("routes[%d].url", i))
		script, _ := config.Get(fmt.Sprintf("routes[%d].script", i))
		method, _ := config.Get(fmt.Sprintf("routes[%d].method", i))
		if method == "" {
			method = "GET"
		}
		r.HandleFunc(url, scriptHandlerGenerator(script)).Methods(method)
	}

	http.Handle("/", r)
}

func main() {

	parseFlags()

	var useSSL = false
	if keyFile != "" && certFile != "" {
		useSSL = true
	}

	loadConfig(useSSL)

	if useSSL {
		log.Printf("[INFO] Starting server on https://%s:%s", host, strconv.FormatUint(port, 10))
		log.Fatalf("[FATAL] %s", http.ListenAndServeTLS(host+":"+strconv.FormatUint(port, 10), certFile, keyFile, nil))
	} else {
		log.Printf("[INFO] Starting server on http://%s:%s", host, strconv.FormatUint(port, 10))
		log.Fatalf("[FATAL] %s", http.ListenAndServe(host+":"+strconv.FormatUint(port, 10), nil))
	}
}
