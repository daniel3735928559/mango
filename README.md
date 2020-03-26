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


### Node type definition

```
[config]
name <node_name>
command <cmd>
env <name> <val>
validate (yes|no)

[usage]
<usage_string>

[interface]
import <filename>
type <type_name> <typespec>
input <name> <typespec>
output <name> <typespec>
```

### EMP definition

```
[config]
name <name>

[nodes]
instance <type> <name> <args>...
dummy <name>
merge <name> <mergepoints>...
gen <name> <values>...

[routes]
<route>
```

### Route definition

```
route   
  : node '>' node
  | node '<' node
  | node '<' '>' node
  | node '>' transforms
node    
  : IDENT
  | IDENT '/' IDENT
transforms 
  : transform '>' node
  | transform '>' transforms
transform 
  : '?' transform_filter
  | '%' transform_edit
  | '=' transform_replace
  | '?' transform_filter '%' transform_edit
  | '?' transform_filter '=' transform_replace
transform_filter 
  : '{' expr '}'
  | IDENT '{' expr '}'
  | IDENT
transform_replace 
  : '{' mapexprs '}'
  | IDENT '{' mapexprs '}'
  | IDENT
transform_edit 
  : '{' script '}'
  | IDENT '{' script '}'
script 
  : stmt
  | stmt script
stmt 
  : dstexpr '=' expr ';'
  | dstexpr AE expr ';'
  | dstexpr OE expr ';'
  | dstexpr XE expr ';'
  | dstexpr PE expr ';'
  | dstexpr ME expr ';'
  | dstexpr TE expr ';'
  | dstexpr DE expr ';'
  | dstexpr RE expr ';'
  | VAR IDENT ';'
  | DEL IDENT ';'
expr 
  : NUMBER
  | TRUE
  | FALSE
  | '{' mapexprs '}'
  | '[' listexprs ']'
  | IDENT '(' listexprs ')'
  | IDENT '(' ')'
  | STRING
  | expr '~' expr
  | expr '?' expr ':' expr
  | '-' expr      %prec UNARY
  | '(' expr ')'
  | varexpr
  | expr EXP expr
  | expr '+' expr
  | expr '-' expr
  | expr '*' expr
  | expr '/' expr
  | expr '&' expr
  | expr '|' expr
  | expr '^' expr
  | expr '%' expr
  | expr EQ expr
  | expr NE expr
  | expr GE expr
  | expr LE expr
  | expr AND expr
  | expr OR expr
  | '!' expr %prec UNARY
  | expr '>' expr
  | expr '<' expr
mapexprs 
  : IDENT ':' expr
  | IDENT ':' expr ',' mapexprs
listexprs 
  : expr
  | expr ',' listexprs
varexpr 
  : expr '.' IDENT
  | expr '[' expr ']'
  | IDENT
dstexpr
  : IDENT
  | THIS
  | dstexpr '[' expr ']'
  | dstexpr '.' IDENT
```

### Value definition

```
value 
  : NUMBER
  | STRING
  | TRUE
  | FALSE
  | '{' mapvals '}'
  | '{' '}'
  | '[' listvals ']'
  | '[' ']'
mapvals 
  : IDENT ':' value
  | IDENT ':' value ',' mapvals
listvals 
  : value
  | value ',' listvals
```

### Value type definition

```
typedesc 
  : STR
  | NUM
  | IDENT
  | BOOL
  | ANY
  | '[' typedesc ']'
  | ONEOF '(' oneofentries ')'
  | '{' mapentries '}'
oneofentries 
  : typedesc
  | typedesc ',' oneofentries
mapentries 
  : mapentry
  | mapentry ',' mapentries
mapentry 
  : IDENT ':' typedesc
  | IDENT '*' ':' typedesc
  | IDENT ':' typedesc '=' value
value 
  : NUMBER
  | STRING
  | TRUE
  | FALSE
  | '{' mapvals '}'
  | '{' '}'
  | '[' listvals ']'
  | '[' ']'
mapvals 
  : IDENT ':' value
  | IDENT ':' value ',' mapvals
listvals 
  : value
  | value ',' listvals
```
