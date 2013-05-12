package main

import (
    "fmt"
    "net/http"
)

func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    buf := make([]byte, r.ContentLength)
    r.Body.Read(buf)
    //TODO:  logging
    fmt.Println("Got controller payload:  ", string(buf))
    if r.URL.Path == "/control" {
        c.controlChannel <- string(buf)
    } else {
        id := lastPathPart(r.URL)
        c.SeqCon(id, string(buf))
    }
}

func StartControlHandler(controller Controller) {
    http.Handle("/control", controller)
    http.Handle("/control/", controller)
}