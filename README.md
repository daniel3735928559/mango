# Mango

Mango is an IPC routing library designed so that mango-based programs
written in any language can all work together cleanly under the user's
complete control.

For example, a mango-based drawing program will come as two separate
programs--the backend (which stores the current image buffer and
implements any functions that can modify the image) and the frontend
(which can display an image and allows the user to modify the image).
Mango allows you to connect the interface to the backend and draw as
usual, but now, if a friend is also running the frontend program on a
separate computer, you can connect that program also (over the
network) to the backend that you are running.  Suddenly, our naive
drawing program is functioning as a shared whiteboard!

The goal of mango is to make it easy to route input and output
messages between mango programs, meaning that you can with equal ease
(and with no further work from the programmer) use mango programs from
a GUI, shell, other program, or over the network, and that you can
readily chain together functions from disparate programs when needed,
in the style of a UNIX pipeline.
