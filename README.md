# gowebexample
[Official Golang Web Development Example](https://gowebexamples.com)

# Hello world
- Registering a Request Handler
`func (w http.ResponseWriter, r *http.Request)`
- Listen for HTTP Connections
`http.ListenAndServe(":80", nil)`
```
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello go, you've requested: %s\n", r.URL.Path)
    })

    http.ListenAndServe(":80", nil)
}
```
# http server
- process dynamic requests
    - http.Request
        - r.URL.Path
        - r.URL.Query().Get("token") // get parameter
        - r.FormValue("email") // post parameter
- serve static assets [JavaScript, CSS, images]
    - http.FileServer
```
fs := http.FileServer(http.Dir("static/"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
```
- accept connections
`http.ListenAndServe(":80", nil)`
```
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "Welcome to Go HTTP Server")
    })

    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(":80", nil)
}
```
# Routing (using gorilla/mux)
- Installing the gorilla/mux package
`go get -u github.com/gorilla/mux`
- Creating a new Router
`r := mux.NewRouter()`
- Registering a Request Handler
`r.HandleFunc()`
- URL Parameters
```
// the url as /books/go-programming-blueprint/page/10
// divide into two dynamic segments as
// book title slug, which is /go-programming-blueprint and
// page (10)
r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    vars["title"] // get the book's title from url
    vars["page"] // get the page from url
})
```
- Setting the HTTP server's router
`http.ListenAndServe(":80", r)`
```
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/books/{title}/page/{page}", func (w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        title := vars["title"]
        page := vars["page"]

        fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
    })

    http.ListenAndServe(":80", r)
}
// visit: http://localhost/books/golangReceipt/page/999
```
# MySQL
- Installing the go-sql-driver/mysql package
`go get -u github.com/go-sql-driver/mysql`
- Connecting to a MySQL database
```
import "database/sql"
import _ "go-sql-driver/mysql"

db, err := sql.Open("mysql", "USERNAME:PASSWORD@(127.0.0.1:3306)/dbname?parseTime=true")

err := db.Ping()
```
- Creating our first database table
```
query := `
    CREATE TABLE users (
        id INT AUTO_INCREMENT,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME,
        PRIMARY KEY(id)
    );`

_, err := db.Exec(query)
```
- Inserting our first record
```
import "time"

username := "johndoe"
password := "secret"
createAt := time.Now()

result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createAt)

userId, err := result.LastInsertId()
```
- Querying our records table [db.Query, db.QueryRow]
```
import "time"

var (
    id int
    username string
    password string
    createdAt time.Time
)

query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt)
```
- Querying all records
```
import "time"

var (
    id int
    username string
    password string
    createdAt time.Time
)

rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
defer rows.Close()

var users []user
for rows.Next() {
    var u user
    err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
    users = append(users, u)
}
err := rows.Err()

/*
users {
    user {
        id:         1,
        username:   "jonhdoe",
        password:   "secret",
        createdAt:  time.Time{wall: 0x0, ext: 63701044325, loc: (*time.Location)(nil)}
    }, 
    user {
        id:         2,
        username:   "janedoe",
        password:   "secretly",
        createdAt:  time.Time{wall: 0x0, ext: 63701044622, loc: (*time.Location)(nil)}
    }, 
}
*/
```
- Deleting a record from our table
`_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1)`
```
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
)

func main(){
    db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/root?parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    { // create a new table
        query := `
        CREATE TABLE users (
            id INT AUTO_INCREMENT,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
            created_at DATETIME,
            PRIMARY KEY(id)
        );`

        if _, err := db.Exec(query); err != nil {
            log.Fatal(err)
        }
    }

    { // insert a new user
        username := "johndoe"
        password := "secret"
        createdAt := time.Now()

        result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
        if err != nil {
            log.Fatal(err)
        }

        id, err := result.LastInsertId()
        fmt.Println(id)
    }

    { // query a single user
        var (
            id int
            username string
            password string
            createdAt time.Time
        )

        query := "SELECT id, username, password, created_at FROM users WHERE id = ?"
        if err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt); err != nil {
            log.Fatal(err)
        }

        fmt.Println(id, username, password, createdAt)
    }

    { // query all users
        type user struct {
            id int
            username string
            password string
            createdAt time.Time
        }

        rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        var users []user
        for rows.Next(){
            var u user

            err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
            if err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%#v", users)
    }

    {
        _, err := db.Exec(`DELETE FROM users WHERE id = ?`)
        if err != nil {
            log.Fatal(err)
        }
    }
}
```
# Templates
- First template
```
data := TodoPageData{
    PageTitle: "My Todo List",
    Todos: []Todo{
        {Title: "Task 1", Done: true},
        {Title: "Task 2", Done: false},
        {Title: "Task 3", Done: true},
    },
}

<h1>{{.PageTitle}}</h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <li class="done">{{.Title}}</li>
        {{else}}
            <li>{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
```
- Control Structures
```
{{/* a comment */}}             // defines a comment
{{.}}                           // renders the root element
{{.Title}}                      // renders the "Title" field in a nested element
{{if .Done}} {{else}} {{end}}   // defines an if-statement
{{range .Todos}} {{.}} {{end}}  // loops over all "Todos" and renders each using {{.}}
{{block "content" .}}{{end}}    // defines a block with the name "content"
```
- Parsing Templates from Files
```
// the layout.html is in the same directory as the go program

tmpl, err := template.ParseFiles("layout.html")
// or
tmpl := template.Must(template.ParseFiles("layout.html"))
```
- Execute a Template in a Request Handler
```
func(w http.ResponseWriter, r *http.Request) {
    tmpl.Execute(w, "tata goes here")
}
```
```
package main

import (
    "html/template"
    "net/http"
)

type Todo struct {
    Title string
    Done bool
}

type TodoPageData struct {
    PageTitle string
    Todos []Todo
}

func main() {
    tmpl := template.Must(template.ParseFiles("template.html"))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := TodoPageData{
            PageTitle: "My Todo List",
            Todos: []Todo{
                {Title: "Task 1", Done: true},
                {Title: "Task 2", Done: false},
                {Title: "Task 3", Done: true},
            },
        }
        tmpl.Execute(w, data)
    })
    http.ListenAndServe(":80", nil)
}
```
```html
    <h1>{{.PageTitle}}</h1>
    <ul>
        {{range .Todos}}
            {{if .Done}}
                <li class="done">{{.Title}}[{{.Done}}]</li>
            {{else}}
                <li>{{.Title}}[{{.Done}}]</li>
            {{end}}
        {{end}}
    </ul>
```
# Assets and Files [static files like CSS, JS, Images]
```
package main

import "net/http"

func main() {
    fs := http.FileServer(http.Dir("assets/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(":8080", nil)
}
```
# Forms
```
package main

import (
    "html/template"
    "net/http"
    "fmt"
)

type ContactDetails struct {
    Email string
    Subject string
    Message string
}

func main(){
    tmpl := template.Must(template.ParseFiles("form.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }
        details := ContactDetails{
            Email: r.FormValue("email"),
            Subject: r.FormValue("subject"),
            Message: r.FormValue("message"),
        }

        _ = details // do something with details var, else compile error
        fmt.Printf("Email: %s\nSubject: %s\nMessage: %s\n",details.Email, details.Subject, details.Message)

        tmpl.Execute(w, struct{Success bool}{true})
    })

    http.ListenAndServe(":8080", nil)
}
```
```html
{{if .Success}}
<h1>Thanks for your feedback messages.<a href="http://localhost:8080">Refresh</a></h1>

{{else}}
<h1>Leave us feedback</h1>
<form action="" method="post">
    <label for="email">Email:</label><br/>
    <input type="email" name="email" id=""><br/>
    <label for="subject">Subject:</label><br/>
    <input type="text" name="subject" id=""><br/>
    <label for="message">Message:</label><br/>
    <textarea name="message" id=""></textarea><br/>
    <input type="submit" value="Submit">
</form>
{{end}}
```
# Middleware (Basic)
```
// middleware [basic]
package main

import (
    "fmt"
    "log"
    "net/http"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL.Path)
        f(w, r)
    }
}

func foo(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "foo")
}

func bar(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "bar")
}

func main(){
    http.HandleFunc("/foo", logging(foo))
    http.HandleFunc("/bar", logging(bar))

    http.ListenAndServe(":8080", nil)
}
```
# Middleware (Advanced)
```
func createNewMiddleware() Middleware{
    // create a new middleware
    middleware := func(next http.HandlerFunc) http.HandlerFunc {
        // define the http.HandleFunc which is called by the server eventually
        handler := func(w http.ResponseWriter, r *http.Request) {
            // do middleware thing

            // call the next middleware/handler in chain
            next(w, r)
        }

        // return newly created handler
        return handler
    }

    // return newly created middleware
    return middleware
}
```
```
// middleware [advanced]
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// logging logs all requests with its path and 
// the time it took to process
func Logging() Middleware {
    
    // create a new Middleware
    return func(f http.HandlerFunc) http.HandlerFunc {

        // define the http.HandlerFunc
        return func(w http.ResponseWriter, r *http.Request) {

            // do middlewarethings
            start := time.Now()
            defer func() {log.Println(r.URL.Path, time.Since(start))}()

            // call the next middleware/handler in chain
            f(w, r)
        }
    }
}

// method ensures that url can only be 
// requested with a specific method, 
// else returns a 400 bad request
func Method(m string) Middleware {

    // create a anew middleware
    return func(f http.HandlerFunc) http.HandlerFunc {

        // define the http.HandlerFunc
        return func(w http.ResponseWriter, r *http.Request) {

            // do middleware things
            if r.Method != m {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                return
            }

            // call the next middleware/handler in chain
            f(w, r)
        }
    }
}

// chain applies middlewares to a http.lHandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    for _, m := range middlewares {
        f = m(f)
    }
    return f
}

func Hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world")
}

func main() {
    http.HandleFunc("/", Chain(Hello, Method("GET"), Logging()))
    http.ListenAndServe(":8080", nil)
}
```
# Sessions
```
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
```
# JSON
```
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type User struct {
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
    Age       int    `json:"age"`
}

func main() {
    http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
        var user User
        json.NewDecoder(r.Body).Decode(&user)

        fmt.Fprintf(w, "%s %s is %d years old!", user.Firstname, user.Lastname, user.Age)
    })

    http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
        peter := User{
            Firstname: "John",
            Lastname:  "Doe",
            Age:       25,
        }

        json.NewEncoder(w).Encode(peter)
    })

    http.ListenAndServe(":8080", nil)
}
```
# Websockets
```
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func main(){
    http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request){
        conn, _ := upgrader.Upgrade(w, r, nil) 

        for {
            // read message from browser
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }

            // print the message to the console
            fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

            // write message back to browser
            if err = conn.WriteMessage(msgType, msg); err != nil {
                return
            }
        }
    })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        http.ServeFile(w, r, "websocket.html")
    })

    http.ListenAndServe(":8080", nil)
}
```
```
<input id="input" type="text" />
<button onclick="send()">Send</button>
<pre id="output"></pre>
<script>
    var input = document.getElementById("input");
    var output = document.getElementById("output");
    var socket = new WebSocket("ws://localhost:8080/echo");

    socket.onopen = function () {
        output.innerHTML += "Status: Connected\n";
    };

    socket.onmessage = function (e) {
        output.innerHTML += "Server: " + e.data + "\n";
    };

    function send() {
        socket.send(input.value);
        input.value = "";
    }
</script>
```
# Password Hashing
```
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func main() {
    password := "secret"
    hash, _ := HashPassword(password)

    fmt.Println("Password:", password)
    fmt.Println("Hash:", hash)

    match := CheckPasswordHash(password, hash)
    fmt.Println("Match: ", match)
}
```