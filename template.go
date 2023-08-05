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