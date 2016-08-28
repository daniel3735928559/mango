# libmango

libmango is the core library needed to handle mango messages.  Each
language needs an instance of this library in order that mango
programs can be written in that language.

This should not be a terribly onerous request, as the library should
be relatively simple (the current Python implementation is under 300
lines of code) and we intend that this document should describe the
requirements fully so that anyone with the requisite knowledge should
be able to implement libmango in a currently unsupported language.

## Current language support

* Python: Completed
* Javascript: Planned
* Clojure: Tentatively planned
* C++: Tentatively planned

## Specification

A libmango implementation generally consists of code that handles the
following tasks:

* Transport: Sending and receiving actual data (currently, via
  ZeroMQ).

* Serializer: Decoding the data received into a map or dictionary
  (currently, by JSON parsing).

* Interface: Storing a dictionary of which functions should be called
  for which commands and dispatching commands accordingly when they
  come in.

* Dataflow: Polling for new data on the transport and, when it comes
  in, passing it off through the serializer for decoding an interface
  for dispatching

* Initialization: The code that defines the functions to call that add
  functions to the interface, registers the node, and starts the main
  polling loop.

* Error: Something to codify and translate the various mango error
  codes.

