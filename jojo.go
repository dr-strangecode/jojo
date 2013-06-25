package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
  "github.com/gorilla/mux"
  "strings"
  "encoding/base64"
)

import (
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

var configFile string
func init() {
	const (
		defaultConfigFile = "/etc/jojo.yaml"
    usage             = "full path to config file: '--config=/path/to/your.yaml'"
	)
	flag.StringVar(&configFile, "config", defaultConfigFile, usage)
	flag.StringVar(&configFile, "c", defaultConfigFile, usage+" (shorthand)")
}

var certFile string
func init() {
	const (
		defaultCertFile   = ""
    usage             = "full path to SSL certificate file: '--cert=/path/to/your.crt'"
	)
	flag.StringVar(&certFile, "cert", defaultCertFile, usage)
}

var keyFile string
func init() {
	const (
		defaultKeyFile    = ""
    usage             = "full path to SSL key file: '--key=/path/to/your.key'"
	)
	flag.StringVar(&keyFile, "key", defaultKeyFile, usage)
}

var host string
func init() {
	const (
		defaultHost       = "localhost"
    usage             = "host or ip to serve on: '--host=localhost'"
	)
	flag.StringVar(&host, "host", defaultHost, usage)
	flag.StringVar(&host, "h", defaultHost, usage+" (shorthand)")
}

var user string
func init() {
  const (
		defaultUser       = ""
    usage             = "user for http basic auth: '--user=joe'"
	)
	flag.StringVar(&user, "user", defaultUser, usage)
}

var password string
func init() {
  const (
		defaultPassword   = ""
    usage             = "password for http basic auth: '--password=secret'"
	)
	flag.StringVar(&password, "password", defaultPassword, usage)
}

var port uint64
func init() {
	const (
		defaultPort = 3000
    usage       = "port to serve on: '--port=3000'"
	)
	flag.Uint64Var(&port, "port", defaultPort, usage)
	flag.Uint64Var(&port, "p", defaultPort, usage+" (shorthand)")
}

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
    log.Printf("[Info] %s %s %s %s %s", r.Proto, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
    log.Printf("Args: %s", r.URL.Query())
    urlCmdArgs := strings.Join(strings.Split(strings.Join(strings.Split(r.URL.RawQuery, "&"), " "), "=")," ")
    log.Printf("[DEBUG] <Header>%s</Header>", r.Header)
    fmt.Fprintf(w, "<p>Hey, I'm going to call: %s %s</p>", script, urlCmdArgs)
    log.Printf("[Info] Running %s %s", script, urlCmdArgs)
    cmd := exec.Command(script, urlCmdArgs)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    cmdErr := cmd.Run()
    if cmdErr != nil {
      fmt.Fprintf(w, "<p>Got an error: %s</p>", cmdErr)
      log.Printf("[ERROR] %s", cmdErr)
    }
    fmt.Fprintf(w, "<p>Results: %s</p>", out.String())
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
    parameters, _ := config.Get(fmt.Sprintf("routes[%d].parameters", i))
    log.Printf("Params: %s", parameters)
    if method == "" {
      method = "GET"
    }
		r.HandleFunc(url, scriptHandlerGenerator(script)).Methods(method)
	}

  http.Handle("/", r)
}

func main() {

	flag.Parse()

  var useSSL = false
  if keyFile != "" && certFile != "" {
    useSSL = true
  }

	loadConfig(useSSL)

  if useSSL {
    log.Printf("[INFO] Starting server on https://%s:%s",host,strconv.FormatUint(port, 10))
    log.Fatalf("[FATAL] %s", http.ListenAndServeTLS(host+":"+strconv.FormatUint(port, 10), certFile, keyFile, nil))
  } else {
    log.Printf("[INFO] Starting server on http://%s:%s",host,strconv.FormatUint(port, 10))
    log.Fatalf("[FATAL] %s", http.ListenAndServe(host+":"+strconv.FormatUint(port, 10), nil))
  }
}
