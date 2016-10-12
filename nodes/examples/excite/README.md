# Excite (mango)

For each language implementation, the corresponding "excite" node
serves as the "hello world" example for that language.

In this document, we cover the prerequisites and process for running
each of the excite nodes in each language.

## Prerequisites

### Javascript

`excite.js` will require only NodeJS to be installed

### C

`excite.c` will need to be compiled, as will the libmango.c library.

To compile the library, go to `/libmango/c` and run `./make.sh`.  To compile `excite.c`, in this directory run `./make.sh`.

### Python

The Python expample will require Python 3 along with the `zmq` and `yaml` modules to be installed, and for `python` to refer to the Python 3 executable (as opposed to Python 2.7, say).

## Running

To run and test any/all of the excite nodes, you will need two nodes running in advance: mc and mx.

### Starting mc

Open a terminal, go to `/nodes/mc` and run `./mc_start 61453 61454`.
Leave this terminal window open.  When anything goes wrong, this is
where the error message will likely be.

### Starting mx

Open another terminal tab, go to `/nodes/mx` and run `./mx_start 61453
mx`.  This will put you into a shell that is configured to communicate
with `mc`.  NB: This shell currently has a lot of debug output.  If
this overwrites your current prompt, just use `C-c` or `Enter` a few
times to get your prompt back.  

You can send `mc` commands with the `mc` command.  For example, `mc
nodes` and `mc routes` will give you the currently running nodes and
the current routes between them, respectively.  `mc types` will show
you all the types of nodes that can be started.  You should see
`excite_*` for various values of `*` in that list.

### Starting excite

To start excite_py, for example, in the `mx` shell, run:

`mc launch -node excite_py -id foo`

You should see `success: True` somewhere in the debug spew, and `foo`
should now appear in the list got by running `mc nodes`.

### Using excite

Now to connect your shell's node, `mx`, to the excite node we just
launched: `foo`.  We need to create a route between them.  The
simplest such would be a two-way route directly between them:

`mc route -map 'mx <> foo'`

This should show `success: True` in the debug spew again.  At this
point, anything you send from the mx node will be received by `foo`,
and anything `foo` sends will be received by you.

The excite node specficially requires `excite` message only, each with
one argument only: `str`.  (This is documented in
[excite.yaml](excite.yaml) in this directory.)  To send an `excite`
message, we use

`mx excite`

And to send it with a `str` argument of, say, `"Hello World"`, we use

`mx excite -str "Hello World"`

You should receive `"excited": "Hello World!"` in the debug spew.  You
did it!

### A more complex example

Suppose you are going to send many messages from your `mx` shell, and
only want the `excite` messages to go to your `foo` node.  Further,
you want to be able to send these with argument `bored` instead of
`str`.

Then add a new route that will allow only `excite` messages through,
and that will change an argument called `bored` to one called `str`:

`mc route -map 'mx > ? excite > {"str":"$bored"} > foo'

(The old route from `foo` to `mx` will still exist since we did not
delete or otherwise overwrite that.)

You should be able to see this route if you do

`mc routes -src_node mx -src_port stdio`

Then if you do:

`mx excite -bored blah`

you will get `"excited":"blah!"` back, whereas if you do

`mx somethingelse -bored blah`

then this will get nothing back.

If you do

`mx excite -str blah`

you will get nothing back and see an error in the `mc` console.

If, instead of filtering only excite commands, you want to ensure that
all commands will go to `foo`, but to have them turned into `excite`
commands, you can use

`mx route -map 'mx > @ "excite" > {"str":"$bored"} > foo > mx'

### An exercise;

Try to launch two excite nodes (possibly in different languages), say
with ids `foo` and `bar`, and have `excite` messages from `mx` go to
foo and then be converted into valid `excite` messages and passed on
to `bar`, and then have the now-twice-excited resutl passed back to
mx`.

ANSWER:

```








mc launch -node excite_c -id foo
mc launch -node excite_js -id bar
mc route -map 'mx > ? "excite" > foo > @ "excite" > {"str":"$excited"} > bar > mx'
```