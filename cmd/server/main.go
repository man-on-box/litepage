package main

import (
	"fmt"
	"log"
	"net/http"
)

const devPort = ":3001"

func main() {
	staticDir := "dist"

	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileServer)

	fmt.Printf("Dev server listening on localhost%s\n", devPort)
	err := http.ListenAndServe("localhost"+devPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
