// Copyright 2013 Tim Ray, All rights reserved.
// Use of this source code is governed by an Apache-style
// License that can be found in the LICENSE file.

package main

import (
	"flag"
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
		defaultCertFile = ""
		usage           = "full path to SSL certificate file: '--cert=/path/to/your.crt'"
	)
	flag.StringVar(&certFile, "cert", defaultCertFile, usage)
}

var keyFile string

func init() {
	const (
		defaultKeyFile = ""
		usage          = "full path to SSL key file: '--key=/path/to/your.key'"
	)
	flag.StringVar(&keyFile, "key", defaultKeyFile, usage)
}

var host string

func init() {
	const (
		defaultHost = "localhost"
		usage       = "host or ip to serve on: '--host=localhost'"
	)
	flag.StringVar(&host, "host", defaultHost, usage)
	flag.StringVar(&host, "h", defaultHost, usage+" (shorthand)")
}

var user string

func init() {
	const (
		defaultUser = ""
		usage       = "user for http basic auth: '--user=joe'"
	)
	flag.StringVar(&user, "user", defaultUser, usage)
}

var password string

func init() {
	const (
		defaultPassword = ""
		usage           = "password for http basic auth: '--password=secret'"
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

func parseFlags() {
	flag.Parse()
}
