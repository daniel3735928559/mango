# Mango

Mango is an IPC routing library designed so that mango-based programs written in any language can all work together cleanly, and in a way that is easy for the user to control.

For example, suppose we write a drawing program.  A mango-based
drawing program will come as two separate programs--the backend, and
some graphical interface.  Mango allows you to connect the interface
to the backend and draw as usual, but now, if your friend is also
running the graphical interface piece on a separate computer, you can
connect that also (over the network) to the backend you're running.
Suddenly, our naive drawing program is functioning as a shared
whiteboard!  

The goal of mango is to make it easy to route input and output
messages between mango programs, meaning that you can with equal ease
(and with no further work from the programmer) use mango programs from
a GUI, shell, other program, or over the network, and that you can
readily chain together functions from disparate programs when needed.
