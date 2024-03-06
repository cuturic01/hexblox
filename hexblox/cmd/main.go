package main

import (
	"flag"
	"hexblox/internal/p2p"
)

func main() {
	httpPort := flag.String("http", "", "HTTP server port")
	hostPort := flag.String("host", "", "WebSocket server port")
	flag.Parse()

	p2p.Run(*httpPort, *hostPort)

	select {}
}
