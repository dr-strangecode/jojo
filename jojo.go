package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
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

var port uint64
func init() {
	const (
		defaultPort = 3000
    usage       = "port to serve on: '--port=3000'"
	)
	flag.Uint64Var(&port, "port", defaultPort, usage)
	flag.Uint64Var(&port, "p", defaultPort, usage+" (shorthand)")
}

func scriptHandler(script string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    log.Printf("[Info] %s %s %s %s", r.Proto, r.Method, r.Host, r.URL.Path)
    log.Printf("[DEBUG] <Header>%s</Header>", r.Header)
    fmt.Fprintf(w, "<p>Hey, I'm going to call: %s</p>", script)
    log.Printf("[Info] Running %s", script)
    cmd := exec.Command(script, "")
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
      fmt.Fprintf(w, "<p>Got an error: %s</p>", err)
      log.Printf("[ERROR] %s", err)
    }
    fmt.Fprintf(w, "<p>Results: %s</p>", out.String())
    log.Printf("[Info] Results: %s", out.String())
	}
}

func loadConfig() {
	config := yaml.ConfigFile(configFile)
	numRoutes, err := config.Count("routes")
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	for i := 0; i < numRoutes; i++ {
		url, _ := config.Get(fmt.Sprintf("routes[%d].url", i))
		script, _ := config.Get(fmt.Sprintf("routes[%d].script", i))
		http.HandleFunc(url, scriptHandler(script))
	}
}

func main() {
	flag.Parse()
	loadConfig()
  log.Printf("[INFO] Starting server on %s:%s",host,strconv.FormatUint(port, 10))
  if keyFile != "" && certFile != "" {
    log.Fatalf("[FATAL] %s", http.ListenAndServeTLS(host+":"+strconv.FormatUint(port, 10), certFile, keyFile, nil))
  } else {
    log.Fatalf("[FATAL] %s", http.ListenAndServe(host+":"+strconv.FormatUint(port, 10), nil))
  }
}
