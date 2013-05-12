gosēq
=====

Google Go + OSX CoreMIDI + REST = gosēq, a simple MIDI sequencer with a RESTful JSON interface

#What

A bare-bones (so far) MIDI sequencer using the [OSX CoreMIDI](http://developer.apple.com/library/ios/#documentation/MusicAudio/Reference/CACoreMIDIRef/_index.html) library.  Sequences are created, destroyed, started and stopped by POSTing JSON documents to the running sequencer.

This depends on my minimal [gocmc](https://github.com/j14159/gocmc) library.  I don't provide packages for this yet but likely will once a rudimentary UI is built.

#Shut Up And Tell Me How To Use It
See [the wiki](https://github.com/j14159/goseq/wiki)

#Why

I used to use [Seq24](http://en.wikipedia.org/wiki/Seq24) a lot on Linux to control a pile of synths ages ago.  While I don't have all the kit I used to, I do have a small [Eurorack](http://en.wikipedia.org/wiki/Doepfer_A-100) modular synth and a couple other outboard pieces alongside a developing fascination with [Pure Data](http://puredata.info/).

While I write software all day to pay the bills, I'm terrible at UI so I thought providing a simple interface for those better in that department could hack together something sensible if they so choose (nudge, nudge, wink, wink).

#What's Missing

The short list of stuff I intend to get to in no particular order:

* a basic HTML/JS UI
* MIDI clock input/output
* proper use of the DELETE method to destroy/remove existing sequences
* error handling (there's almost none right now)
* sequencing MIDI CC messages

Consider this to be pretty brittle for the time being.

#How

See [the wiki](https://github.com/j14159/goseq/wiki).  I'll try to keep that as up to date as possible but will probably fail.
