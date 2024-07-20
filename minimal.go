package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func main() {
	addHelloHandler()
	addCookieHandler()
	addSessionHandler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addSessionHandler() {
	cmap := map[string]int{}
	handler := func(w http.ResponseWriter, req *http.Request) {
		uid := ""
		if cookie, err := req.Cookie("session"); err != nil {
			uid = uuid.NewString()
			log.Default().Printf("No session found. Creating a new session: %s", uid)
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: uid,
			})
			cmap[uid] = 0
		} else {
			uid = cookie.Value
		}
		cmap[uid] += 1
		str := fmt.Sprintf("You have visited: %d times.", cmap[uid])
		log.Default().Printf(str)
		io.WriteString(w, str)
	}
	http.HandleFunc("/session", handler)
}

func addHelloHandler() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, World!\n")
	}
	http.HandleFunc("/hello", handler)
}

func addCookieHandler() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		count := 0
		if c, err := req.Cookie("count"); err == nil {
			if count, err = strconv.Atoi(c.Value); err != nil {
				log.Default().Print(err)
				count = 0
			}
		}
		count += 1
		http.SetCookie(w, &http.Cookie{
			Name:  "count",
			Value: strconv.Itoa(count),
		})

		str := fmt.Sprintf("You have visited %d times.", count)
		log.Default().Print(str)
		io.WriteString(w, str)
	}

	http.HandleFunc("/count", handler)
}
