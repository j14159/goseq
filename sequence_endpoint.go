package main

import (
	"encoding/json"
	"fmt"
	"github.com/j14159/gocmc"
	"net/http"
	"net/url"
	"strings"
)

type SequenceEndpoint struct {
	controller  Controller
	client      gocmc.Client
	outputCache OutputCache
}

type DestinationUpdates struct{}

type JsonSequence struct {
	Running       bool
	MidiChannel   int
	StepLength    int
	StartTick     int
	NoteOffOnStop bool
	Output        string
	Steps         []Step
}

type SimpleJsonSequence struct {
	MidiChannel int
	Output      string
	StepLength  int
	StartTick   int
	Steps       []int
}

//Converts a list of integers to a list of note on/off events suitable for a Sequence.  
//Used to change SimpleJsonSequence step arrays to real note on/off events.  -1 indicates 
//a rest (sends a note off for the last note seen) and a -2 indicates a tie (no immediate note off).
func intsToSteps(notes []int) []Step {
	steps := make([]Step, len(notes))

    

	if len(notes) > 0 {
        lastNote := notes[len(notes) -1]
		for i, n := range notes {
            var stepNotes []gocmc.NoteEvent
            if n > -1 {
                //if we have a note, put it in slot 1
                stepNotes = make([]gocmc.NoteEvent, 2)
                lastNote = n
                stepNotes[0] = gocmc.NoteEvent{ false, lastNote, 0 }
                stepNotes[1] = gocmc.NoteEvent{ true, n, 127 }
            } else if n == -1 {
                stepNotes = make([]gocmc.NoteEvent, 1)
                stepNotes[0] = gocmc.NoteEvent{ false, lastNote, 0 }
            } else {
                stepNotes = make([]gocmc.NoteEvent, 0) //empty step for a tie
            }

            steps[i] = Step{ stepNotes }

            /*
			stepNotes := make([]gocmc.NoteEvent, 2)
            if n > -1 {
			 stepNotes[1] = gocmc.NoteEvent{true, n, 127}
             lastNote = n
            }
			if i == 0 && notes[len(notes) - 1] > -1 {
				stepNotes[0] = gocmc.NoteEvent{ false, notes[len(notes)-1], 0 }
			} else {
                if n != -2 {
				    stepNotes[0] = gocmc.NoteEvent{false, lastNote, 0}
                }
			}

            if n > -1 {
			 steps[i] = Step{stepNotes}
            }*/
		}
	}

	return steps
}

//Simple helper method to grab the last part of a URL path,
//used to grab sequence IDs.
func lastPathPart(url *url.URL) string {
	pathParts := strings.Split(url.Path, "/")
	return pathParts[len(pathParts)-1]
}

func (s SequenceEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := lastPathPart(r.URL)
	fmt.Println("Sequence ID: ", id)
	seq := JsonSequence{}
	simpleSeq := SimpleJsonSequence{}

	buf := make([]byte, r.ContentLength)
	r.Body.Read(buf)

	fmt.Println("Got sequence payload:  ", string(buf))
	if e := json.Unmarshal(buf, &seq); e == nil {
		fmt.Fprintf(w, "ok")
		output := s.outputCache.GetOutput(seq.Output, seq.MidiChannel)

		sequence := Sequence{id, seq.Running, seq.StepLength, seq.StartTick, seq.NoteOffOnStop, seq.Steps, output}
		s.controller.Add(sequence)
	} else if e := json.Unmarshal(buf, &simpleSeq); e == nil {
		fmt.Fprintf(w, "ok")
		output := s.outputCache.GetOutput(simpleSeq.Output, simpleSeq.MidiChannel)
		sequence := Sequence{id, true, simpleSeq.StepLength, simpleSeq.StartTick, true, intsToSteps(simpleSeq.Steps), output}
		s.controller.Add(sequence)
	} else {
		fmt.Fprintf(w, "Problem deserializing your payload:  %s", e)
	}
}

func StartSequenceHandler(controller Controller, client gocmc.Client) {
	http.Handle("/seq/", SequenceEndpoint{controller, client, NewOutputCache(client)})
	return
}
