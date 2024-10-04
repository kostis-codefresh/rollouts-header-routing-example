package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type SmartCanary struct {
	AppVersion         string
	AppColor           string
	CanaryHeaderName   string
	CurrentHeaderValue string
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	smartCanary := SmartCanary{}
	smartCanary.CanaryHeaderName = "X-Canary"

	smartCanary.AppVersion = os.Getenv("APP_VERSION")
	if len(smartCanary.AppVersion) == 0 {
		smartCanary.AppVersion = "1.0"
	}

	smartCanary.AppColor = os.Getenv("APP_COLOR")
	if len(smartCanary.AppColor) == 0 {
		smartCanary.AppColor = "#44B3C2" //Blue 44B3C2 and Yellow F1A94E
	}

	http.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<div class='pod' style='background:%s'> ver: %s\n </div>", smartCanary.AppColor, smartCanary.AppVersion)
	})

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	http.HandleFunc("/", smartCanary.serveFiles)

	fmt.Println("Listening now at port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func (smartCanary *SmartCanary) serveFiles(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	p := "." + upath
	fmt.Printf("Path:Upath is %s:%s\n", p, upath)

	if p == "./" {
		smartCanary.home(w, r)
		return
	} else {
		p = filepath.Join("./static/", path.Clean(upath))
	}

	http.ServeFile(w, r, p)
}

func (smartCanary *SmartCanary) home(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	err = t.Execute(w, smartCanary)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}
