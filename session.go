package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/sessions"
)

var (
    // key must be 16, 24 or 32 bytes long
    // AES-128, AES-192 OR AES-256
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden, login first!", http.StatusForbidden)
        return
    }

    // print secret mess
    fmt.Fprintln(w, "This is the super secret message, you are fooled")
}

func login(w http.ResponseWriter, r *http.Request){
    session, _ := store.Get(r, "cookie-name")

    // authentication goes here
    // set user as authenticated
    session.Values["authenticated"] = true
    session.Save(r, w)

    fmt.Println("There you are, you've logined")
}

func logout(w http.ResponseWriter, r *http.Request){
    session, _ := store.Get(r, "cookie-name")

    // revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)

    fmt.Println("There you are, you've logouted")
}

func main() {
    http.HandleFunc("/secret", secret)
    http.HandleFunc("/login", login) 
    http.HandleFunc("/logout", logout) 

    http.ListenAndServe(":8080", nil)
}