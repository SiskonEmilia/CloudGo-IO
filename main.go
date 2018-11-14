package main

import (
	"os"

	"github.com/siskonemilia/CloudGo-IO/routes"
	flag "github.com/spf13/pflag"
)

const (
	// Default host configuration
	defaultHost string = "localhost"
	// Default port configuration
	defaultPort string = "8000"
)

func main() {
	// Get the router from routes model
	r := routes.Router()

	// Trying to fetch the host configuration in environment variables
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = defaultHost
	}

	// Trying to fetch the port configuration in environment variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	// Trying to fetch the host & port configuration in the flags
	hostFlag := flag.StringP("hostname", "h", defaultHost, "The host to deploy at.")
	portFlag := flag.StringP("port", "p", defaultPort, "The port for listening.")
	flag.Parse()
	if len(*hostFlag) != 0 {
		host = *hostFlag
	}
	if len(*portFlag) != 0 {
		port = *portFlag
	}

	// Run the server @host:port
	r.Run(host + ":" + port)
}
