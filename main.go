package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var PORT = ":8080"
var counter = 0

func main() {
	r := mux.NewRouter()

	log.Printf("Веб-сервер запущен на http://127.0.0.1%s", PORT)

	targetURL1, errProxy1 := url.Parse("http://127.0.0.1:4000")

	if errProxy1 != nil {
		log.Fatal("Ошибка при парсинге первого URL:", errProxy1)
	}

	targetURL2, errProxy2 := url.Parse("http://127.0.0.1:4010")

	if errProxy2 != nil {
		log.Fatal("Ошибка при парсинге второго URL:", errProxy2)
	}

	// Создаем два прокси для каждого целевого сервера
	proxy1 := httputil.NewSingleHostReverseProxy(targetURL1)
	proxy2 := httputil.NewSingleHostReverseProxy(targetURL2)

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if counter%2 == 0 {
			log.Println("request server 1")
			proxy1.ServeHTTP(writer, request)
		} else {
			log.Println("request server 2")
			proxy2.ServeHTTP(writer, request)
		}
		counter++

	})
	r.HandleFunc("/create", func(writer http.ResponseWriter, request *http.Request) {
		if counter%2 == 0 {
			log.Println("request server 1")
			proxy1.ServeHTTP(writer, request)
		} else {
			log.Println("request server 2")
			proxy2.ServeHTTP(writer, request)
		}
		counter++

	})
	r.HandleFunc("/friends/{id}", func(writer http.ResponseWriter, request *http.Request) {
		if counter%2 == 0 {
			log.Println("request server 1")
			proxy1.ServeHTTP(writer, request)
		} else {
			log.Println("request server 2")
			proxy2.ServeHTTP(writer, request)
		}
		counter++

	})
	r.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		if counter%2 == 0 {
			log.Println("request server 1")
			proxy1.ServeHTTP(writer, request)
		} else {
			log.Println("request server 2")
			proxy2.ServeHTTP(writer, request)
		}
		counter++

	})
	r.HandleFunc("/make_friends", func(writer http.ResponseWriter, request *http.Request) {
		if counter%2 == 0 {
			log.Println("request server 1")
			proxy1.ServeHTTP(writer, request)
		} else {
			log.Println("request server 2")
			proxy2.ServeHTTP(writer, request)
		}
		counter++

	})
	r.HandleFunc("/{user_id}", func(writer http.ResponseWriter, request *http.Request) {
		if counter%2 == 0 {
			log.Println("request server 1")
			proxy1.ServeHTTP(writer, request)
		} else {
			log.Println("request server 2")
			proxy2.ServeHTTP(writer, request)
		}
		counter++

	})

	var err = http.ListenAndServe(PORT, r)

	log.Fatal(err)
}
