# libmango

libmango is the core library needed to handle mango messages.  Each
language needs an instance of this library in order that mango
programs can be written in that language.

This should not be a terribly onerous request, as the library should
be relatively simple (the current Python implementation is 257 lines
of code) and we intend that this document should describe the
requirements fully so that anyone with the requisite knowledge should
be able to implement libmango in a currently unsupported language.

A libmango implementation is more or less a simple RPC server with
ZeroMQ as the transport and JSON for serialisation.  The protocol used
differs from e.g. JSON-RPC and other standard RPC protocols in a few
ways.  For example, there is no expectation of a request/reply flow,
and there is a built-in addressing system of nodes and ports, rather
than using the transport for addressing as is standard.

## Current language support

* Python: Completed
* Javascript: Completed
* C: Completed (alpha)
* Clojure: Planned
* C++: Tentatively planned

## Outline of the task

### Example

To describe what a Mango program should look like, we shall work with
an example: An `excite` program, which has one function called
`excite` that takes in a string argument `str` and returns a string
argument `excited`, which is the input string with an '!'  appended to
it.

This program consists of two files:
`[excite.yaml](../nodes/example/excite/excite.yaml)`, which describes
the functions of the program and their inputs and outputs:

```
excite:
  args:
    str:
      type: string
  rets:
    excited:
      type: string
```

This uses the interface descriptor language from
[Pijemont](https://github.com/daniel3735928559/pijemont), and
documentation for the format can be found there.

Then, the program that implements this function and allows it to be
called through Mango RPC looks simply like: 

```
from libmango import *

class excite(m_node):
    def __init__(self):
        super().__init__()
        self.interface.add_interface('excite.yaml', {'excite':self.excite})
        self.run()
    def excite(self,header,args):
        return {'excited':args['str']+'!'}
t = excite()
```

Let us walk through this a little bit:

First, we import libmango to give us access to the RPC functionality: 

```
from libmango import *
```

This program is a mango node, and the main class should extend the m_node class:

```
class excite(m_node):
```

The first thing we do when we initialise the class is to initialise
the node as well, by calling the superclass's constructor:

```
  def __init__(self):
    super().__init__()
```

Once we have initialised the node, we have an interface member
variable that we can add to, using the add_interface function to
register that we are implementing the functions from excite.yaml, and
specifically that the 'excite' function specified in the YAML file is
implemented by the Python function self.excite: 

```
    self.interface.add_interface('excite.yaml', {'excite':self.excite})
```

Finally, we start the main loop of the program: 


```
    self.run()
```


We need also to write the function that implements the 'excite'
function from the interface descriptor.  All this needs to do is
return a dictionary with 'excited' as its only key.  In this case, the
corresponding value is the 'excited' version of the input value 'str'.

```
  def excite(self,header,args):    
    return {'excited':args['str']+'!'}
```

Having now defined the class, we can simply instantiate it:

```
excite()
```

### Behind the scenes

So the task in writing a libmango implementation, broadly, is to write
the library that enables this code to work.  Here, "work" means the
following:

* The program creates node object (here, on line ) a ZeroMQ dealer socket and connects it to a
  ZeroMQ router socket whose address is specified in the `MC_ADDR`
  environment variable.  Concretely, if the program is run with
  `MC_ADDR=tcp://localhost:61453`, then the program needs to have code
  to connect to this.

* The program sends a serialised Mango RPC 'hello' message (specified
  below) on this socket and receives and handles a serialised Mango
  RPC 'reg' message (also specified below) in response.  Concretely,
  this will look like sending a message such as:

  ```
MANGO0.1 json
{"header":{"command":"hello","src_node":"excite","src_port","mc"},"args":{"id":"ex","if":{"excite":{"args":{"str":{"type":"string"}},"rets":{"excited":{"type":"string"}}}}}}
  ```

  and will receive a response like:

  ```
MANGO0.1 json
{"header":{"command":"reg","src_node":"mc","src_port","mc"},"args":{"id":"ex.0"}}
  ```

  at which point the node's ID should be set to "ex.0" with this being
  used as the "src_node" parameter in all future communications.

* The program can then receive serialised Mango RPC messages calling
  the `excite` function with various parameters and turn these into
  actual calls to the `excite` function defined above, and will send
  back Mango RPC 'reply' messages with the return value of this
  function as the body.  Concretely, it might receive messages like:

  ```
MANGO0.1 json
{"header":{"command":"excite","src_node":"some_node","src_port","stdio","port":"stdio"},"args":{"str":"Hello World"}}
  ```

  and will respond like: 

  ```
MANGO0.1 json
{"header":{"command":"reply","src_node":"ex.0","src_port","stdio"},"args":{"excited":"Hello World!"}}
  ```

## Mango RPC specification

### Architecture

Every RPC endpoint is called a "node", which may have one or more
"ports".  Every node has a unique ID, every port of a given node has a
string name.  Messages are always marked with their source node/port
and their destination port, as well as the function they are calling.

### Message format

Every message is formatted thus:

  ```
MANGO[version nunber] [serialisation method]
[dictionary with "header" and "args" keys, serialised with the specified serialisation method]
  ```

Most current implementations simply use JSON to serialise the message
content.  Since Mango is currently on version 0.1, these messages look like: 

  ```
MANGO0.1 json
{"header":..., "args":...}
  ```

The header dictionary, when sent out, contains the following keys:

* `command`: The name of the function to call
* `src_port`: The port from which we are sending the call

When a program receives a command, the header dictionary will have the
following keys:

* `command`: The name of the function to call

* `src_node`: The node from which the call is being made

* `src_port`: The port from which the call is being made

* `port`: The port on our node on which the function call is being
  made.

This is because when you send a command, your identity and the
destination of your command is controlled by the central router, so
rather than specify where you want your command to do, you simply send
a command and an indication of which of your ports is emitting the
command, and then the router will attach all the information about
which node sent it and which node and port it gets routed to.

### Base interface:

Every node must implement the basic node interface, whose descriptor
is found in [node_if.yaml](node_if.yaml).  In addition, every node may
implement any number of non-conflicting interfaces loaded from
separate YAML files.  As mentioned before, all interface descriptor
files use the interface descriptor language from
[Pijemont](https://github.com/daniel3735928559/pijemont), which is
documented in that repository.

### Dataflow

Given this, the usual dataflow of the libmango implementation happens
in two steps: Initialisation (which usually happens when the libmango
Node class (in whatever form) is instantiated), interface definition
(where we set up the functions we want to be accessible from the
outside world, and the main loop (usually started by calling the `run`
function of the Node object).

#### Initialisation

* Make a ZMQ dealer socket and connect to the ZMQ server

* Send the "hello" message, which has command "hello" and port "mc",
  and therefore has header dictionary:

  ```
  {"command":"hello","port":"mc")
  ```

  and further has arguments `id` and `if`, which are, respectively the
  node ID that we want for this node and the dictionary describing the
  interface for the node (which is the dictionary compiled from the
  various interface descriptor files for the interfaces we are
  implementing).

#### Interface prep

* Load all the desired YAML files (libmango will load the default
  node.yaml and implement all its functions, but the program may at
  this stage load any number of other interface descriptor files), and
  create a mapping of function names in those interface files to the
  corresponding actual functions.

#### Main loop

* Polls the ZMQ dealer socket (and possibly any other sockets the user
  wishes to poll in this loop, although there is another pattern
  recommended for integrating multiple "main loops").

* When we receive input on this socket, first parse and deserialise
  the input into its consitutuent header and arguments dictionaries.

* From the header dictionary, determine the command they wish to call,
  look up the correspoding function in the aforementioned mapping, and
  call it with the supplied arguments.

* If it returns something, send that as a "reply" command on the port
  on which the message was received.  If an error occurs in the
  function, (e.g. an exception, in languages that support them) return
  a description of the error in an "error" command on the port "mc".
  Otherwise, do nothing.

* If in the course of the program, a command gets sent, which will
  specify the command name, the arguments, and optionally the port.
  If the port is unspecified, it should default to "stdio".  Sending a
  message should create the header:

  ```
  {"command":[command name], "src_port":[port name]}
  ```

  It should then serialise the message with this as the header and
  with the supplied arguments as the arguments and otuput them on the
  node's main dealer ZMQ socket.


## Components

A libmango implementation generally consists of code that handles the
following tasks (whether as separate files, functions, or just pieces
of one large blob of code, but in principle it should be easy to swap
out methods of transporting data, serialising, or whatever; if someone
wants to do all of their transport by carrier pigeon and serialise
messages as XML RPC calls, then that should be made easy):

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
