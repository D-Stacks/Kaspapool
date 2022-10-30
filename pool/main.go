package main

import (
	"KPool/pool_server"
)

func main() {
	pool, err := pool_server.NewPool()
	if err != nil {
		panic(err)
	}
	pool.RunForever()
}