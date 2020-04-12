package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func testFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	// Default configuration settings
	viper.SetDefault("StaticDir", "./static/")
	viper.SetDefault("Listen", "localhost:8100")
	viper.SetDefault("Webhook", nil)

	// Parse configuration files
	viper.SetConfigName("werdz")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	r := mux.NewRouter()

	// sample API
	r.HandleFunc("/api/test/", testFunc)

	// another sample API
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	// serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(viper.GetString("StaticDir"))))

	// Test with webhook
	// post(fmt.Sprintf("# Starting server\nTime: %s", time.Now().String()))

	fmt.Printf("Serving static files from: %s\n", viper.GetString("StaticDir"))
	fmt.Printf("Listening: %s\n", viper.GetString("Listen"))
	http.Handle("/", r)
	if err := http.ListenAndServe(viper.GetString("Listen"), nil); err != nil {
		panic(fmt.Errorf("unable to start server: %s", err))
	}
}
