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