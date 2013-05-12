package main

import (
    "time"
)

//Controller tracks a number of channels used for manipulating the
//sequencer state.
type Controller struct {
    //Used to add sequences to the sequencer
    sequenceChannel chan Sequence
    //Used to start and stop the controller's master clock
    controlChannel chan string
    //Provides control (start, stop, destroy) over individual sequences
    sequenceControlChannel chan SequenceControl
}

//Encapsulates necessary information for individual sequence control messages.
//Message should be one of "start", "stop" or "destroy"
type SequenceControl struct {
    SequenceId string
    Message string
}

//Used only within Controllers for tracking an individual sequences channels.
type sequenceChannels struct {
    clockChan chan int
    controlChan chan string
}

//Uses the controller's sequence channel to add a new sequence.
func (con Controller) Add(seq Sequence) {
    con.sequenceChannel <- seq
}

//Uses the controller's sequence control channel to pass control
//messages to individual Sequences.
func (con Controller) SeqCon(seqId string, msg string) {
    con.sequenceControlChannel <- SequenceControl{ seqId, msg }
}

//Simplifies construction of a new Controller, will start the main
//control function as a go routine.
func NewController(pulsesPerCycle int, stepLen int) Controller {
    seqChan := make(chan Sequence)
    conChan := make(chan string)
    seqConChan := make(chan SequenceControl)

    ret := Controller{ seqChan, conChan, seqConChan }
    go ret.controller(pulsesPerCycle, stepLen, make(map[string] sequenceChannels))
    return ret
}

//Runs the master control loop.
func (con *Controller) controller(stepLimit int, stepLen int, controls map[string] sequenceChannels) {
    running := false
    tickChan := time.Tick((time.Duration)(stepLen) * time.Millisecond)
    step := 0

    for true {
        select {
            case <-tickChan:
                if running {
                    for _, con := range controls {
                        con.clockChan <- step
                    }
                    if step += 1; step >= stepLimit {
                        step = 0
                    }
                }
            case newSeq := <-con.sequenceChannel:
                clock, control := make(chan int), make(chan string)
                //need to stop and remove the old version of an extant sequence:
                if oldSeq, exists := controls[newSeq.Id]; exists {
                    oldSeq.controlChan <- SEQ_DESTROY
                    delete(controls, newSeq.Id)
                }

                go newSeq.Stepper(clock, control)
                controls[newSeq.Id] = sequenceChannels{clock, control}
            case command := <-con.controlChannel:
                switch command {
                case "start":
                    running = true
                case "stop":
                    running = false
                    step = 0
                case "reset":
                    step = 0
                }
            case seqCommand := <-con.sequenceControlChannel:
                controls[seqCommand.SequenceId].controlChan <- seqCommand.Message
                if seqCommand.Message == SEQ_DESTROY {
                    delete(controls, seqCommand.SequenceId)
                }
        }
    }
}