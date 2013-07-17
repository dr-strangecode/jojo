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
	"time"
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

func scriptHandlerGenerator(script string, stringTimeout string) http.HandlerFunc {
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
		exitRe := regexp.MustCompile("^exit status ")
		paramsRe := regexp.MustCompile("=|&")
		jsonRe := regexp.MustCompile("&")
		jsonParamRe := regexp.MustCompile("=")
		urlCmdArgs := strings.Split(paramsRe.ReplaceAllString(r.URL.RawQuery, " "), " ")
		jsonUrlCmdArgs := strings.Split(jsonRe.ReplaceAllString(r.URL.RawQuery, " "), " ")

		// Misc Logging junk
		log.Printf("[Info] %s %s %s %s %s", r.Proto, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
		//log.Printf("Args: %s", r.URL.Query())
		//log.Printf("[DEBUG] <Header>%s</Header>", r.Header)
		log.Printf("[Info] Running %s %s", script, urlCmdArgs)

		// JSON - start json and script name
		fmt.Fprintf(w, "{\"script\": %s, ", strconv.Quote(script))

		// JSON - Script arguments - spit out either an empty list of an array of k/v pairs
		fmt.Fprintf(w, "\"arguments\": [")
		for i, paramPair := range jsonUrlCmdArgs {
			entries := strings.Split(jsonParamRe.ReplaceAllString(paramPair, " "), " ")
			if i + 1 < len(jsonUrlCmdArgs) {
				for j, entry := range entries {
					if len(entry) == 0 {
					} else if j == 0 {
						fmt.Fprintf(w, "{%s: ", strconv.Quote(entry))
					} else {
						fmt.Fprintf(w, "%s}, ", strconv.Quote(entry))
					}
				}
			} else {
				for j, entry := range entries {
					if len(entry) == 0 {
					} else if j == 0 {
						fmt.Fprintf(w, "{%s: ", strconv.Quote(entry))
					} else {
						fmt.Fprintf(w, "%s}", strconv.Quote(entry))
					}
				}
			}
		}
		fmt.Fprintf(w, "], ")

		// Run the script passing in the arguments
		cmd := exec.Command(script, urlCmdArgs...)
		var stderr, stdout bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Timer and running command
		t0 := time.Now()

		timeout, err := time.ParseDuration(stringTimeout)
		if err != nil {
			cmdErr := cmd.Run()
			if cmdErr != nil {
				fmt.Fprintf(w, "\"error\": \"%s\", ", cmdErr)
				log.Printf("[ERROR] %s", cmdErr)
			}
		} else {
			cmd.Start()
			done := make(chan error)
			go func() {
				done <- cmd.Wait()
			}()
			select {
			case <-time.After(timeout):
				if err := cmd.Process.Kill(); err != nil {
					fmt.Fprintf(w, "\"error\": \"failed to kill process: %s\", ", err)
					log.Printf("[ERROR] %s", err)
				}
				//<-done // allow goroutine to exit
				fmt.Fprintf(w, "\"error\": \"process timed out and was killed\", ")
				log.Printf("[INFO] process killed after timeout")
			case err := <-done:
				if err != nil {
					log.Printf("[INFO] process finished with error = %v", err)
				}
			}
		}
		t1 := time.Now()
		fmt.Fprintf(w, "\"duration\": \"%v\", ", t1.Sub(t0))


		// JSON - exit status
		exitStatus, _ := strconv.ParseInt(exitRe.ReplaceAllString(cmd.ProcessState.String(), ""), 10, 16)
		fmt.Fprintf(w, "\"exit-status\": %d, ", exitStatus)

		// JSON - stdout - array of lines
		stdOutString := strings.Split(strings.TrimSpace(stdout.String()), "\n")
		fmt.Fprintf(w, "\"stdout\": [")
		for i, line := range stdOutString {
			if i + 1 < len(stdOutString) {
				fmt.Fprintf(w, "%s, ", strconv.Quote(line))
			} else if len(line) == 0 {
			} else {
				fmt.Fprintf(w, "%s", strconv.Quote(line))
			}
		}
		fmt.Fprintf(w, "], ")

		// JSON - stderr - array of lines
		stdErrString := strings.Split(strings.TrimSpace(stderr.String()), "\n")
		fmt.Fprintf(w, "\"stderr\": [")
		for i, line := range stdErrString {
			if i + 1 < len(stdErrString) {
				fmt.Fprintf(w, "%s, ", strconv.Quote(line))
			} else if len(line) == 0 {
			} else {
				fmt.Fprintf(w, "%s", strconv.Quote(line))
			}
		}
		fmt.Fprintf(w, "]")
		fmt.Fprintf(w, "}")

		// Moar logging
		log.Printf("[Info] State: %s", cmd.ProcessState)
		log.Printf("[Info] Duration: %v", t1.Sub(t0))
		log.Printf("[Info] Stdout: %s", stdout.String())
		log.Printf("[Info] Stderr: %s", stderr.String())
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
		stringTimeout, _ := config.Get(fmt.Sprintf("routes[%d].timeout", i))
		if method == "" {
			method = "GET"
		}
		r.HandleFunc(url, scriptHandlerGenerator(script, stringTimeout)).Methods(method)
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
