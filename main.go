package main

import (
    "fmt"
    "net/http"
    "gocmc"
    "flag"
)

func main() {
    bpm := flag.Int("bpm", 60, "beats per minute")
    ppqn := flag.Int("ppqn", 24, "pulses per quarter note, default is 24 and that's probably what you want")
    cycle := flag.Int("cycle", 24, "number of quarter notes before the cycle count resets, can affect sequence start/stop times")
    host := flag.String("host", "localhost", "address to listen for HTTP")
    port := flag.String("port", "4000", "port for HTTP requests")
    flag.Parse()

    stepLen := 60000 / *bpm / *ppqn
    pulsesPerCycle := *ppqn * *cycle
    listenAddress := *host + ":" + *port

    fmt.Printf("bpm:  %d\nppqn:  %d\ncycle length:  %d\n", *bpm, *ppqn, *cycle)
    fmt.Printf("pulse length (ms):  %d\npulse count per cycle:  %d\n", stepLen, pulsesPerCycle)
    fmt.Printf("listening for HTTP at:  %s\n", listenAddress)

    
    destinations := gocmc.GetDestinations()
    fmt.Println("Destinations: ", destinations)

    client := gocmc.MakeClient("goseq")

    fmt.Println("Starting controller and handler")
    con := NewController(pulsesPerCycle, stepLen)
    StartSequenceHandler(con, client)

    StartDestinationsHandler()
    StartControlHandler(con)

    fmt.Println("Starting server")
    http.ListenAndServe(listenAddress, nil)
}