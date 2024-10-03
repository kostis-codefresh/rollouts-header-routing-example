package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

func main() {

	version := "1.0"

	color := "#44B3C2" //Blue 44B3C2 and Yellow F1A94E

	http.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<div class='pod' style='background:%s'> ver: %s\n </div>", color, version)
	})

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	http.HandleFunc("/", serveFiles)

	fmt.Println("Listening now at port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	p := "." + upath
	fmt.Printf("Path:Upath is %s:%s\n", p, upath)

	p = filepath.Join("./static/", path.Clean(upath))

	http.ServeFile(w, r, p)
}
