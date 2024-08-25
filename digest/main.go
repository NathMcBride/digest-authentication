package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/NathMcBride/web-authentication/digest/authentication/authenticator"
	"github.com/NathMcBride/web-authentication/digest/authentication/contexts"
	"github.com/NathMcBride/web-authentication/digest/authentication/digest"
	"github.com/NathMcBride/web-authentication/digest/authentication/middleware"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

/*
	func main() {
		addHelloHandler()
		addCookieHandler()
		addSessionHandler()
		addBasicAuthHandler()
		addFormBasedAuthHandler()
		addDigestAuthHandler()

		log.Fatal(http.ListenAndServe(":8080", nil))
	}
*/

func somethingProtected(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("Something protected"))
}

func somethingThatNeedsAuthenticatingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := contexts.GetSession(ctx)
	somethingThatNeedsAuthenticating(w, r, session)
}

func somethingThatNeedsAuthenticating(w http.ResponseWriter, r *http.Request, session *authenticator.Session) {
	log.Print("Executing somethingThatNeedsUser")
	pp.Println(session)
	w.Write([]byte("Something that needs a user"))
}

func somethingPublic(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("Something public"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	options := middleware.Options{
		Realm:              "A-Realm",
		Opaque:             digest.RandomKey(),
		ShouldHashUsername: true,
	}

	da := middleware.NewDigestAuth(options)

	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(somethingThatNeedsAuthenticatingHandler)

	mux.Handle("/protected", da.RequireAuth(finalHandler))
	mux.Handle("/public", http.HandlerFunc(somethingPublic))
	mux.Handle("/health", http.HandlerFunc(handleHealth))
	mux.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
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

func addBasicAuthHandler() {
	pmap := map[string]string{"jdoe": "password"}
	handler := func(w http.ResponseWriter, req *http.Request) {
		if u, p, ok := req.BasicAuth(); ok {
			if pmap[u] == p {
				str := fmt.Sprintf("User %s authenticated", u)
				io.WriteString(w, str)
				log.Default().Print(str)
			} else {
				str := fmt.Sprintf("User %s failed to authenticate", u)
				w.WriteHeader(http.StatusUnauthorized)
				log.Default().Print(str)
			}
		} else {
			w.Header().Add("WWW-Authenticate", "Basic Realm=\"Access Server\"")
			w.WriteHeader(http.StatusUnauthorized)
			log.Default().Print("Basic authentication needed.")
		}

	}
	http.HandleFunc("/basicauth", handler)
}

type dClient struct {
	nc       uint64
	lastSeen int64
}

/*
func addDigestAuthHandler() {
	// pmap := map[string]string{"jdoe": "password"}
	opaque := RandomKey()
	// var clients map[string]*dClient

	handler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("HERE")
		if auth, ok := DigestAuth(req); ok {
			fmt.Println("HERE2")

			pp.Println(auth)
			// if pmap[u] == p {
			// 	str := fmt.Sprintf("User %s authenticated", u)
			// 	io.WriteString(w, str)
			// 	log.Default().Print(str)
			// } else {
			// 	str := fmt.Sprintf("User %s failed to authenticate", u)
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	log.Default().Print(str)
			// }
		} else {

			nonce := RandomKey()
			// clients[nonce] = &dClient{nc: 0, lastSeen: time.Now().UnixNano()}

			header := fmt.Sprintf(`
				Digest Realm="Access Server",
				qop="auth",
				algorithm=SHA-256,
				nonce="%s",
				opaque="%s"`,
				nonce, opaque)

			w.Header().Add("WWW-Authenticate", header)
			w.WriteHeader(http.StatusUnauthorized)
			log.Default().Print("Digest authentication needed.")
		}

	}
	http.HandleFunc("/digestauth", handler)
}
*/
// func KeyFound(m map[string]interface{}, key string) bool {
// 	_, ok := m[key]
// 	return ok;
// }

// func KeyNotFound(m map[string]interface{}, key string) bool {
// 	_, ok := m[key]
// 	return !ok;
// }

func KeyNotFound[T any](m map[string]T) func(key string) bool {
	return func(key string) bool {
		_, ok := m[key]
		return !ok
	}
}

/*
	func DigestAuth(r *http.Request) (authorization.AuthorizationInfo, bool) {
		authHeader := r.Header.Get("Authorization")
		var auth authorization.AuthorizationInfo
		emptyAuthInfo := authorization.AuthorizationInfo{}
		if authHeader == "" {
			return auth, false
		}

		parsed, ok := ParseDigestAuth(authHeader)
		if !ok {
			return emptyAuthInfo, false
		}

		var fieldExists bool
		if auth.Response, fieldExists = parsed["response"]; !fieldExists {
			return emptyAuthInfo, false
		}

		if auth.UserID, fieldExists = parsed["username"]; !fieldExists {
			return emptyAuthInfo, false
		}

		if auth.Realm, fieldExists = parsed["realm"]; !fieldExists {
			return emptyAuthInfo, false
		}

		if auth.Algorithm, fieldExists = parsed["algorithm"]; !fieldExists {
			return emptyAuthInfo, false
		}

		if auth.Qop, fieldExists = parsed["qop"]; !fieldExists {
			return emptyAuthInfo, false
		}

		if auth.Cnonce, fieldExists = parsed["cnonce"]; !fieldExists {
			return emptyAuthInfo, false
		}

		if auth.Nc, fieldExists = parsed["nc"]; !fieldExists {
			return emptyAuthInfo, false
		}
		return auth, fieldExists
	}
*/
func addFormBasedAuthHandler() {
	smap := map[string]string{}
	pmap := map[string]string{"jdoe": "password1"}

	http.HandleFunc("/resource", func(w http.ResponseWriter, req *http.Request) {
		if cookie, err := req.Cookie("session"); err != nil {
			w.Header().Add("Location", "/login")
			w.WriteHeader(http.StatusFound)
		} else {
			uid := cookie.Value
			user := smap[uid]
			if user != "" {
				str := fmt.Sprintf("User %s authenticated.", user)
				log.Default().Printf("Session %s found. Allowing user %s to access", uid, user)
				io.WriteString(w, str)
			} else {
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusFound)
			}
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		form := `<form method="GET" enctype="application/X-WWW-form-urlencoded">
		<label for="user">Username:</label><br>
		<input type="text" id="user" name="user"><br>
		<label for="password">Password:</label><br>
		<input type="text" id="password" name="password"><br>
		<input type="submit" value="Submit">
		</form>`
		user := req.FormValue("user")
		pass := req.FormValue("password")

		if user == "" || pass == "" {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(form))
		} else {
			if pmap[user] == pass {
				str := fmt.Sprintf("User %s authenticated.", user)
				log.Default().Print(str)
				uid := uuid.NewString()
				log.Default().Printf("No session found. Creating a new session: %s", uid)
				http.SetCookie(w, &http.Cookie{
					Name:  "session",
					Value: uid,
				})
				smap[uid] = user
				w.Header().Add("Location", "/resource")
				w.WriteHeader(http.StatusFound)
			} else {
				str := fmt.Sprintf("User %s failed to authenticate.", user)
				log.Default().Print(str)
				w.Header().Add("Content-Type", "text/html")
				w.Write([]byte(form))
			}
		}
	})

}
