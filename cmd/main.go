package main

import "github.com/thaison199py/multi-threaded-redis/internal/server"

func main() {
	server.RunIoMultiplexingServer()
}
