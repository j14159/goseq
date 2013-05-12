package main

import (
    "gocmc"
    "encoding/json"
    "net/http"
)

func destinationsToJson(d map[string] gocmc.Destination) ([]byte) {
    builder := make([]string, len(d))
    i := 0
    for _, dest := range d {
        builder[i] = dest.Name
        i++
    }
    ret, _ := json.Marshal(builder)
    return ret
}

func (du DestinationUpdates) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    destinations := gocmc.GetDestinations()

    w.Write(destinationsToJson(destinations))
}

func StartDestinationsHandler() {
    http.Handle("/destinations", DestinationUpdates{})
}