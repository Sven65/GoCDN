package main

import (
	"fmt"
	"net/http"
	"flag"
	"io"
	"os"
	"encoding/json"
	"github.com/BurntSushi/toml"
)

type Upload struct{
	Url string
}

type Config struct {
	Baseurl string
}

var conf Config

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		w.WriteHeader(404)
		w.Write([]byte("404 Not found"))
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./f/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		resp := &Upload{Url: conf.Baseurl+"f/"+handler.Filename}
		res, err := json.Marshal(resp)
		fmt.Print(res)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(200)
  		w.Write(res)
		//fmt.Fprintf(w, "%s", res)
	}
}

func main(){

	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil{
		fmt.Println(err);
	}

	fmt.Println(conf)

	port := flag.Int("port", 80, "The port to listen to") // Self explanatory

	flag.Parse() // Parse the flags

	http.HandleFunc("/", upload)
	http.HandleFunc("/ul", upload)
	fmt.Printf("Listening on port %d", *port) // Print serving
	http.ListenAndServe(fmt.Sprintf(":%d",*port), nil) // Listen for requests
}
