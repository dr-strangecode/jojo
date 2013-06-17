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
		usage             = "full path to config file"
	)
	flag.StringVar(&configFile, "config", defaultConfigFile, usage)
	flag.StringVar(&configFile, "c", defaultConfigFile, usage+" (shorthand)")
}

var port uint64
func init() {
	const (
		defaultPort = 3000
		usage       = "port to serve on"
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
	log.Println("[INFO] Starting server on port "+strconv.FormatUint(port, 10))
	log.Fatalf("[FATAL] %s", http.ListenAndServe(":"+strconv.FormatUint(port, 10), nil))
}
