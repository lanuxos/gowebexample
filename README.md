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