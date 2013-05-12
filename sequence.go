package main

import (
    "fmt"
    "github.com/j14159/gocmc"
)

const (
    SEQ_DESTROY = "destroy"
    SEQ_START = "start"
    SEQ_STOP = "stop"
)

type Sequence struct {
    Id string
    Running bool
    StepLength int
    StartTick int
    NoteOffOnStop bool
    Steps []Step    
    Output gocmc.Output
}

type Step struct {
    Notes []gocmc.NoteEvent
}

type NoteChange struct {
    step int
    note int
}

func (seq Sequence) SetOutput(output gocmc.Output) {
    seq.Output = output
}

func (seq *Sequence) Stepper(clockChan chan int, controlChan chan string) {
    step := 0
    //this flips to false when the sequence is told to destroy itself:
    live := true

    /*
    running means the sequence should be paying attention to step values.
    started means that it should actually output notes.

    The combination is used so that sequences can specify specific pulses they
    want to be started on in order to synchronize with others.  When an external
    signal says to start the sequence, running = true and started = false UNTIL
    the controlTick value received matches the desired start value specified
    in the original sequence.  Once controlTick modulo seq.StartTick == 0, started = true
    and the sequence stepper will begin to output steps to its output until told
    to stop or destroy.
    */
    running := seq.Running
    started := false

    for live {
        select {
        case controlTick := <-clockChan:            
            if (controlTick % seq.StepLength) == 0 {
                //fmt.Printf("Step %d is note %d\n", step, seq.Steps[step])
                if(running && started) {
                    for _, note := range(seq.Steps[step].Notes) {
                        seq.Output.EventOut(note)
                    }
                } else if running == true && controlTick % seq.StartTick == 0 {
                    started = true
                    for _, note := range(seq.Steps[step].Notes) {
                        seq.Output.EventOut(note)
                    }
                }

                if step += 1; step >= len(seq.Steps) {
                    step = 0
                }
            }
        case controlMsg := <-controlChan:
            fmt.Printf("Sequence stepper got control message:", controlMsg)

            switch controlMsg {
            case SEQ_DESTROY:
                seq.cleanup()
                live = false
            case SEQ_STOP:
                seq.cleanup()
                running = false
            case SEQ_START:
                running = true
            }
        }
    }
}

//Sends a note-off with velocity 0 to all notes on the Sequence's output.
func (seq *Sequence) cleanup() {
    for i := 0; i < 128; i++ {
        seq.Output.EventOut(gocmc.NoteEvent{false, i, 0})
    }
}

