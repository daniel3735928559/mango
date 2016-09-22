#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "args.h"
#include "error.h"

char *LIBMANGO_VERSION = "0.1";

typedef struct m_node {
  char *version;
  char *node_id;
  uint32 mid;
  char **ports;
  char debug;
  char *server_addr;
  m_interface_t *interface;
  m_serialiser_t *serialiser;
  m_transport_t *local_gateway;
  m_dataflow_t *dataflow;
} m_node_t;

void m_node_new(char debug){
  m_node_t *n = malloc(sizeof(m_node_t));
  n->version = LIBMANGO_VERSION;
  n->debug = debug;
  n->node_id = NULL;//...
  n->mid = 0;
  /*
    this.serialiser = new Serialiser(this.version);
    this.iface = new Interface();
    this.node_id = process.env['MANGO_ID'];
    this.mid = 0;
    this.ports = [];
    var server = process.env['MC_ADDR'];
this.iface.add_interface('/home/zoom/suit/mango/libmango/node_if.yaml',{
        'reg':self.reg, 'reply':self.reply, 'heartbeat':self.heartbeat
    });
    this.local_gateway = new Transport(server);
    var s = this.local_gateway.socket;
    this.dataflow = new Dataflow(self.iface, self.local_gateway, self.serialiser, self.dispatch, self.handle_error);
    s.on('message',function(data){ self.dataflow.recv(data); });

  */
  return n;
}

void m_node_ready(m_node_t *node, m_header_t *header, m_args_t *args){
	/* var ifce = {}; */
	/* for(var i in this.iface.iface){ */
	/*     ifce[i] = JSON.parse(JSON.stringify(this.iface.iface[i])); */
	/*     delete ifce[i]['handler']; */
	/* } */
	/* console.log("IF",ifce) */
	/* self.m_send('hello',{'id':self.node_id,'if':ifce,'ports':self.ports},"reg",null,"mc") */
}

void m_node_dispatch(m_node_t *node, m_header_t *header, m_args_t *args){
	/* console.log("DISPATCH",header,args,self.iface.iface,self.iface.iface[header['command']]); */
	/* try{ */
        /*     result = self.iface.iface[header['command']]['handler'](header,args); */
        /*     if(result && 'callback' in header){ */
	/* 	self.m_send(header['callback'],result,null,header['mid'],header['port']) */
	/*     } */
	/* } catch(e) { */
	/*     console.log(e); */
        /*     self.handle_error(header['src_node'],e+"") */
	/* } */
}

void m_node_dispatch(m_node_t *node, char *src, char *err){
  /* console.log('OOPS',src,err); */
  /* self.m_send('error',{'source':src,'message':err},null,null,"mc"); */
}

void m_node_reg(m_node_t *node, m_header_t *header, m_args_t *args){
  //self.node_id = args["id"];
}

void m_node_reply(m_node_t *node, m_header_t *header, m_args_t *args){
  //console.log("REPLY",header,args);
}

void m_node_heartbeat(m_node_t *node, m_header_t *header, m_args_t *args){
  //self.m_send("alive",{},null,null,"mc");
}

m_header_t *m_node_heartbeat(m_node_t *node, char *command, char *callback, int mid, char *src_port){
  if(!callback) callback = LIBMANGO_REPLY;
  if(!mid) mid = m_node_get_mid(node);
  if(!src_port) src_port = LIBMANGO_STDIO;
  m_header_t *header = m_header_new();
  header->src_node = node->node_id;
  header->src_port = src_port;
  header->mid = mid;
  header->command = command;
  header->callback = callback;
  return header;
}
    
int m_node_get_mid(m_node_t *node){
  return node->mid++;
}

int m_node_send(m_node_t *node, char *command, m_args_t *msg, char *callback, int mid, char *port){
  /* console.log('sending',command,msg,mid,port) */
  /*   header = self.make_header(command,callback,mid,port) */
  /*   self.dataflow.send(header,msg) */
  /*   return header['mid'] */
}
    

function MError(message){
    this.message = message;
    this.name = "Mango Error";
}

function Interface(){
    var self = this;
    this.iface = {};
    
    this.add_interface = function(if_file,handlers){
	try {
	    var spec = jsyaml.safeLoad(fs.readFileSync(if_file, 'utf8'));
	    console.log(JSON.stringify(spec));
	    var missing = [], extra = [], existing = [];
	    for(var i in spec)
		if(!(i in handlers)) missing.push(i);
	    for(var i in handlers)
		if(!(i in spec)) extra.push(i);
	    for(var i in spec)
		if(i in self.iface) existing.push(i);
	    if(missing.length > 0) throw new MError("Functions not implemented: "+missing.join(", "));
	    if(extra.length > 0) throw new MError("Functions not in interface: "+extra.join(", "));
	    if(existing.length > 0) throw new MError("Functions already implemented: "+existing.join(", "));

	    for(var i in handlers){
		self.iface[i] = spec[i] ? spec[i] : {};
		self.iface[i]['handler'] = handlers[i];
	    }
	} catch (e) {
	    console.log(e);
	    throw new MError("Failed to load interface");
	}
    }

    this.validate = function(fn){
	return fn in self.iface;
    }
}

function Transport(target){
    var self = this;
    this.target = target
    this.socket = zmq.socket("dealer");
    this.socket.connect(target);
    
    this.tx = function(data){
	self.socket.send(data);
    }

}

function Dataflow(iface,transport,serialiser,dispatch_cb,error_cb){
    var self = this;
    this.iface = iface;
    this.transport = transport;
    this.serialiser = serialiser;
    this.dispatch_cb = dispatch_cb;
    this.error_cb = error_cb;
    console.log(this.dispatch_cb);
    this.send = function(header,args){
	self.transport.tx(self.serialiser.serialise(header,args));
    }
    
    this.recv = function(data){
	try{
	    var m = self.serialiser.deserialise(data);
	    if(!self.iface.validate(m[0]["command"])) throw new MError("Unknown function");
	    console.log(self.dispatch_cb);
	    self.dispatch_cb(m[0],m[1]);
	} catch(e) {
	    console.log(e);
	    self.error_cb(e);
	}
    }

}


function Serialiser(version){
    var self = this;
    this.version = version;
    this.method = "json";
    
    this.make_preamble = function(){
	return "MANGO"+self.version+" json\n";
    }

    this.parse_preamble = function(data){
	console.log("GOT",data);
	var nl1 = data.indexOf('\n');
	var m = data.substring(0,nl1).match(/^MANGO([0-9.]*) ([^ ]*)$/);
	if(!m || m.length < 1) throw new MError("Preamble failed to parse");
	return [m[1],m[2],data.substring(nl1)];
    }

    this.serialise = function(header,args){
	return self.make_preamble()+JSON.stringify({"header":header,"args":args});
    }
    
    this.deserialise = function(data){
	data = data.toString('utf8');
	p = self.parse_preamble(data);
	var ver = p[0]; var method = p[1]; var message = p[2];
	if(ver != self.version) throw new MError("Version mismatch");
	if(method != self.method) throw new MError("Unsupported method");
	try{
	    var d = JSON.parse(message);
	    return [d['header'],d['args']]
	} catch (e) {
	    console.log(e);
	    throw new MError("Failed to parse message");
	}
    }
}

module.exports = MNode;
