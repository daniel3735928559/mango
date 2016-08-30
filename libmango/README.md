# libmango

libmango is the core library needed to handle mango messages.  Each
language needs an instance of this library in order that mango
programs can be written in that language.

This should not be a terribly onerous request, as the library should
be relatively simple (the current Python implementation is 257 lines
of code) and we intend that this document should describe the
requirements fully so that anyone with the requisite knowledge should
be able to implement libmango in a currently unsupported language.

## Current language support

* Python: Completed
* Javascript: Planned
* Clojure: Tentatively planned
* C++: Tentatively planned

## Outline of the task

At a high level, a mango program runs the following steps:

1. Store a mapping of function names to actual functions

2. Connect to a central hub (currently via ZMQ dealer socket)

3. Send the central hub a "hello" message, which looks like:

```
MANGO [version number]
{"header":{"command":"hello",...},"args":{"id":[id requested],"if":[dictionary describing functions and their arguments]}
```

4. The response will be a registration message, which will look like: 

```
MANGO [version number]
{"header":{"command":"reg",...},"args":{"id":[id given]}
```

This should be handled, and all future communications should use the
given ID as the "src" field in the header.

5. Then, after any further program-specific initialisation, the socket
connecting the node to the hub should be polled for any incoming data,
which should be deserialised into the header dictionary and the
arguments dictionary.  The header dictionary will contain a "command"
field, which will specify the function being called.  This should
exist in the mapping stored in step 1.  This function should be
called, with the header and arguments dictionaries passed in as the
two arguments.  The function should return a dictionary, which should
be sent out as the arguments dictionary to a "reply" command.  The
"mid" in the header of the reply command should be the same as the
"mid" in the header of the command being replied to.

6. To facilitate all this sending of messages, there should be a "send
message" function (often called `m_send`) which serialises and sends a
message given: a command, a source port (which should default to
"stdio"), and an arguments dictionary.  

## Components

A libmango implementation generally consists of code that handles the
following tasks (whether as separate files or as :

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

