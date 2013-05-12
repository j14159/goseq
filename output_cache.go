//Caches gocmc.CoreMidiEndpoints so that we're not constantly 
//creating new output ports, etc when all that changes on an output
//is often the midi channel.
package main

import (
	"fmt"
	"github.com/j14159/gocmc"
)

//Holds details for the asynchronous cache request, you should
//use OutputCache.GetOutput instead of dealing with this directly.
type OutputRequest struct {
	outputName   string
	midiChannel  int
	replyChannel chan gocmc.Output
}

//Holds the CoreMIDI client and cache's request channel.  Use NewOutputCache
//rather than constructing this directly _unless_ you don't want the 
//underlying go routine started to handle cache requests right away.
type OutputCache struct {
	client         gocmc.Client
	requestChannel chan OutputRequest
}

//Creates a new OutputCache and spins up the underlying go routine.
func NewOutputCache(client gocmc.Client) (cache OutputCache) {
	requestChan := make(chan OutputRequest)
	cache = OutputCache{client, requestChan}
	go cache.outputCacheHandler()
	return
}

//Convenience method so that users of an OutputCache don't need to deal
//directly with channels and OutputRequests.  Same idea behind
//gen_server:call in a way.
func (cache OutputCache) GetOutput(outputName string, midiChannel int) gocmc.Output {
	replyChan := make(chan gocmc.Output)
	cache.requestChannel <- OutputRequest{outputName, midiChannel, replyChan}
	return <-replyChan
}

//Run this as a go routine to handle asynchronous output
//cache requests.
func (cache OutputCache) outputCacheHandler() {
	_cache := make(map[string]map[int]gocmc.Output)

	for true {
		select {
		case req := <-cache.requestChannel:
			//have we cached _anything_ for this output yet?
			if outputs, outputCached := _cache[req.outputName]; outputCached {
				if output, exists := outputs[req.midiChannel]; exists {
					req.replyChannel <- output
				} else {
					fmt.Printf("New channel for %s\n", req.outputName)
					destination := gocmc.GetDestinations()[req.outputName].Endpoint
					output := cache.client.NewOutput(req.outputName+"-output", req.midiChannel, destination)
					outputs[req.midiChannel] = output
					_cache[req.outputName] = outputs
					req.replyChannel <- output
				}
			} else {
				fmt.Println("No existing output, creating")
				destination := gocmc.GetDestinations()[req.outputName].Endpoint
				output := cache.client.NewOutput(req.outputName+"-output", req.midiChannel, destination)
				fmt.Println("New output: ", output)
				outputs := make(map[int]gocmc.Output)
				outputs[req.midiChannel] = output
				_cache[req.outputName] = outputs
				req.replyChannel <- output
			}
		}
	}
}
