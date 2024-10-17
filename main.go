package main

import "net/http"

func main() {
	app := http.NewServeMux()
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	
	http.ListenAndServe(":8080", app)
}